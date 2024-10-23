package flow

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"os/exec"
	"os/user"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/Techzy-Programmer/d2m/config/db"
	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/Techzy-Programmer/d2m/config/types"
)

var currentDeployID uint
var deployLogStorage []db.DeploymentLog
var deploymentRunning bool = false
var deployQueue []*types.DeploymentRequest

func preInit(req *types.DeploymentRequest) (*db.Deployment, *deployLogHandler, string, error) {
	if runtime.GOOS == "windows" {
		req.LocalParentPath = strings.ReplaceAll(req.LocalParentPath, "/", "\\")
	} else {
		req.LocalParentPath = strings.ReplaceAll(req.LocalParentPath, "\\", "/")
	}

	deploymentRunning = true
	currentDeployID = generateSecure4DigitNumber()

	deployment := &db.Deployment{
		ID:      currentDeployID,
		StartAt: time.Now(),
	}

	depLogger := newDeployLogHandler("Warm Up", "Preparing for deployment...").
		logStep(fmt.Sprintf("Starting deployment of with id \"%d\"", currentDeployID)).
		logStep("Running pre-deployment checks...")

	homeDir, hdirErr := getUserHomeDirectory(req.LocalUser)
	if hdirErr != nil {
		depLogger.logErr(fmt.Sprintf("Error getting home directory for local user \"%s\"", req.LocalUser)).save(errLvl)
		paint.Error("Error getting user home directory: ", hdirErr)
		return nil, nil, "", hdirErr
	}

	// Ensure the parent path exists
	parentPath := path.Join(homeDir, req.LocalParentPath)
	if info, ppErr := os.Stat(parentPath); os.IsNotExist(ppErr) || !info.IsDir() {
		depLogger.logErr(fmt.Sprintf("Parent path \"%s\" does not exist or is not a directory", parentPath)).save(errLvl)
		paint.Error("Parent path is not a valid directory: ", parentPath)
		return nil, nil, "", ppErr
	}

	depLogger.logOk("[✓] Checks passing").save(okLvl)

	return deployment, depLogger, parentPath, nil
}

func generateSecure4DigitNumber() uint {
	min := 1000
	max := 9999
	rangeVal := big.NewInt(int64(max - min + 1))

	n, err := rand.Int(rand.Reader, rangeVal)
	if err != nil {
		return 0
	}

	return uint(n.Int64() + int64(min))
}

const (
	infoLvl uint = iota
	warnLvl
	errLvl
	okLvl
)

func getUserHomeDirectory(username string) (string, error) {
	usr, err := user.Lookup(username)
	if err != nil {
		if runtime.GOOS == "windows" {
			host, hostErr := os.Hostname()
			if hostErr != nil {
				return "", errors.New("error getting machine hostname")
			}

			usr, err = user.Lookup(host + "\\" + username)
			if err != nil {
				return "", errors.New("user not found: " + username)
			}

			return usr.HomeDir, nil
		}

		return "", errors.New("user not found: " + username)
	}

	// Return the home directory
	return usr.HomeDir, nil
}

func saveDeployment(deploy *db.Deployment, success *bool) {
	if *success {
		deploy.Status = "success"
	} else {
		deploy.Status = "failed"
	}

	newDeployLogHandler("Finish Up", "Deployment completed with status: "+deploy.Status).save(infoLvl)
	deploy.Logs = deployLogStorage
	deploy.EndAt = time.Now()

	db.SaveDeployment(deploy)
	deploymentRunning = false
	deployLogStorage = []db.DeploymentLog{}
}

func execCmds(cmds []string, wdPath string, stopOnErr bool, logger *deployLogHandler) error {
	for _, cmd := range cmds {
		cmdParts := strings.Split(cmd, " ")
		ex := exec.Command(cmdParts[0], cmdParts[1:]...)
		ex.Dir = wdPath

		stdoutPipe, _ := ex.StdoutPipe()
		stderrPipe, _ := ex.StderrPipe()

		logger.logStep(fmt.Sprintf("$ %s", cmd))
		if startErr := ex.Start(); startErr != nil {
			logger.logErr(fmt.Sprintf("[X] %v", startErr))

			if stopOnErr {
				return startErr
			}

			logger.logWarn("[!] Recovery, continuing deployment...")
			paint.WarnF("Error starting command (%s): %v\n%s", cmd, startErr, "Continuing...")

			continue
		}

		var stdoutBuilder, stderrBuilder strings.Builder
		stdoutBytes, _ := io.ReadAll(stdoutPipe)
		stdoutBuilder.Write(stdoutBytes)

		stderrBytes, _ := io.ReadAll(stderrPipe)
		stderrBuilder.Write(stderrBytes)
		runErr := ex.Wait()

		if runErr != nil {
			logger.logErr(fmt.Sprintf("[X] %s (%v)", stderrBuilder.String(), runErr))

			if stopOnErr {
				return runErr
			}

			logger.logWarn("[!] Recovery, continuing deployment...")
			paint.WarnF("Error running command (%s): %v\n%s", cmd, runErr, "Continuing...")
		} else {
			logger.logOk(fmt.Sprintf("[✓] %s", stdoutBuilder.String()))
		}
	}

	return nil
}

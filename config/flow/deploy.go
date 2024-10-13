package flow

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/Techzy-Programmer/d2m/config/db"
	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/Techzy-Programmer/d2m/config/types"
)

func StartDeployment(req *types.DeploymentRequest) {
	if deploymentRunning {
		deployQueue = append(deployQueue, req)
		return
	}

	if runtime.GOOS == "windows" {
		req.LocalParentPath = strings.ReplaceAll(req.LocalParentPath, "/", "\\")
	} else {
		req.LocalParentPath = strings.ReplaceAll(req.LocalParentPath, "\\", "/")
	}

	success := false
	deploymentRunning = true
	currentDeployID = generateSecure4DigitNumber()

	deployment := &db.Deployment{
		ID:      currentDeployID,
		StartAt: time.Now(),
		Branch:  req.Branch,
		Repo:    req.RepoPath,
	}

	defer saveDeployment(deployment, &success)
	depLogger := newDeployLogHandler("Warm Up", "Preparing for deployment...").
		logStep(fmt.Sprintf("Starting deployment of \"%s\" from branch \"%s\" with id \"%d\"", req.RepoPath, req.Branch, currentDeployID)).
		logStep("Running pre-deployment checks...")

	homeDir, hdirErr := getUserHomeDirectory(req.LocalUser)
	if hdirErr != nil {
		depLogger.logErr(fmt.Sprintf("Error getting home directory for local user \"%s\"", req.LocalUser)).save(errLvl)
		paint.Error("Error getting user home directory: ", hdirErr)
		return
	}

	// Ensure the parent path exists
	parentPath := path.Join(homeDir, req.LocalParentPath)
	if info, ppErr := os.Stat(parentPath); os.IsNotExist(ppErr) || !info.IsDir() {
		depLogger.logErr(fmt.Sprintf("Parent path \"%s\" does not exist or is not a directory", parentPath)).save(errLvl)
		paint.Error("Parent path is not a valid directory: ", parentPath)
		return
	}

	depLogger.logOk("[✓] Checks passing").save(okLvl)

	if req.PreDeployCmds != nil && len(req.PreDeployCmds) > 0 {
		depLogger.reset("Pre-Deploy Commands", "Starting execution...")

		preExErr := execCmds(req.PreDeployCmds, parentPath, req.FailOnCmdError, depLogger)
		if preExErr != nil {
			depLogger.logErr("[X] Fatal error running pre-deployment commands").save(errLvl)
			paint.Error("Error running pre-deployment commands: ", preExErr)
			return
		}

		depLogger.save(okLvl)
	}

	depLogger.reset("Repository Fetch", "")

	// Let's fetch the repo from GitHub
	hash, msg, ghErr := ensureGHRepo(req.RepoPath, parentPath, depLogger)
	if ghErr != nil {
		depLogger.logStep("[X] Fatal error fetching remote repository").save(errLvl)
		paint.Error("Error fetching GitHub repository: ", ghErr)
		return
	}

	deployment.CommitHash = hash
	deployment.CommitMsg = msg
	depLogger.save(okLvl)

	// ToDo: Implement AutoSetupDeps with smart inference

	if req.PostDeployCmds != nil && len(req.PostDeployCmds) > 0 {
		depLogger.reset("Post-Deploy Commands", "Starting execution...")

		postExErr := execCmds(req.PostDeployCmds, parentPath, req.FailOnCmdError, depLogger)
		if postExErr != nil {
			depLogger.logErr("[X] Fatal error running post-deployment commands").save(errLvl)
			paint.Error("Error running post-deployment commands: ", postExErr)
			return
		}

		depLogger.save(okLvl)
	}

	success = true
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
			return startErr
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

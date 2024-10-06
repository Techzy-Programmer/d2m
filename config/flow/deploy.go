package flow

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/Techzy-Programmer/d2m/config/db"
	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/Techzy-Programmer/d2m/config/types"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func StartDeployment(req *types.DeploymentRequest) {
	if deploymentRunning {
		deployQueue = append(deployQueue, req)
		return
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

	storeDeployLog(infoLvl, "Starting deployment...")
	defer saveDeployment(deployment, &success)

	storeDeployLog(infoLvl, "Running pre-deployment checks...")

	homeDir, hdirErr := getUserHomeDirectory(req.LocalUser)
	if hdirErr != nil {
		storeDeployLog(errLvl, fmt.Sprintf("Error getting local user home directory: %s", homeDir))
		paint.Error("Error getting user home directory: ", hdirErr)
		return
	}

	// Ensure the parent path exists
	parentPath := path.Join(homeDir, req.LocalParentPath)
	if info, ppErr := os.Stat(parentPath); os.IsNotExist(ppErr) || !info.IsDir() {
		storeDeployLog(errLvl, fmt.Sprintf("Parent path is not a valid directory: %s", parentPath))
		paint.Error("Parent path is not a valid directory: ", parentPath)
		return
	}

	storeDeployLog(infoLvl, "Starting execution of pre deployment commands...")

	preExErr := execCmds(req.PreDeployCmds, parentPath, req.FailOnCmdError)
	if preExErr != nil {
		storeDeployLog(errLvl, "[X]: Fatal error running pre-deployment commands")
		paint.Error("Error running pre-deployment commands: ", preExErr)
		return
	}

	// Let's fetch the repo from GitHub
	ghErr := ensureGHRepo(req.RepoPath, parentPath, 3, deployment)
	if ghErr != nil {
		storeDeployLog(errLvl, "[X]: Fatal error fetching GitHub repository")
		paint.Error("Error fetching GitHub repository: ", ghErr)
		return
	}

	// ToDo: Implement AutoSetupDeps with smart inference

	storeDeployLog(infoLvl, "Starting execution of post deployment commands...")

	// Run the deployment commands
	depExErr := execCmds(req.PostDeployCmds, parentPath, req.FailOnCmdError)
	if depExErr != nil {
		storeDeployLog(errLvl, "[X]: Fatal error running post-deployment commands")
		paint.Error("Error running post-deployment commands: ", depExErr)
		return
	}

	success = true
}

func ensureGHRepo(repoPth string, parentPath string, retry int, deployment *db.Deployment) error {
	storeDeployLog(infoLvl, "Fetching latest changes from remote repository...")
	repoParts := strings.Split(repoPth, "/")
	appName := repoParts[1]
	appPth := path.Join(parentPath, appName)
	authTok := db.GetConfig[string]("user.GHPAT")
	var authOpt transport.AuthMethod = nil

	if authTok != "" {
		authOpt = &http.BasicAuth{
			Username: repoParts[0],
			Password: authTok,
		}
	}

	if _, err := os.Stat(appPth); !os.IsNotExist(err) {
		paint.Info("Repository already exists: ", appName)
		repo, poErr := git.PlainOpen(appPth)
		if poErr != nil {
			storeDeployLog(errLvl, fmt.Sprintf("[X]: Error opening local repository at path `%s`", appPth))
			return poErr
		}

		// Pull the latest changes
		paint.Info("Pulling latest changes for: ", appName)
		if wt, err := repo.Worktree(); err != nil {
			storeDeployLog(errLvl, "[X]: Error getting worktree for repository")
			return err
		} else {
			pullErr := wt.Pull(&git.PullOptions{
				Auth: authOpt,
			})

			if pullErr != nil {
				storeDeployLog(errLvl, "[X]: Error pulling changes from remote repository")
				return pullErr
			}
		}

		hash, msg, commitErr := getCommitData(repo)
		if commitErr != nil {
			storeDeployLog(errLvl, "[X]: Error getting commit data for ref HEAD")
			return commitErr
		}

		storeDeployLog(okLvl, fmt.Sprintf("[✓]: Pull successful commit: %s", hash))
		deployment.CommitHash = hash
		deployment.CommitMsg = msg
		return nil
	}

	// Clone the repository
	storeDeployLog(infoLvl, "Repository does not exist locally, cloning...")
	paint.Info("Cloning repository: ", appName)
	repo, cloneErr := git.PlainClone(appPth, false, &git.CloneOptions{
		URL:  "https://github.com/" + repoPth,
		Auth: authOpt,
	})

	if cloneErr != nil {
		if retry > 0 {
			paint.Error("Error cloning repository: ", cloneErr)
			paint.Info("Retrying...")
			_ = os.RemoveAll(appPth)
			storeDeployLog(warnLvl, "[!] Clone failed, retrying...")

			return ensureGHRepo(repoPth, parentPath, retry-1, deployment)
		}

		storeDeployLog(errLvl, "[X]: Error cloning repository")
		return cloneErr
	}

	hash, msg, commitErr := getCommitData(repo)
	if commitErr != nil {
		storeDeployLog(errLvl, "[X]: Error getting commit data for ref HEAD")
		return commitErr
	}

	storeDeployLog(okLvl, fmt.Sprintf("[✓]: Clone successful commit: %s", hash))
	deployment.CommitHash = hash
	deployment.CommitMsg = msg
	return nil
}

func execCmds(cmds []string, wdPath string, stopOnErr bool) error {
	for _, cmd := range cmds {
		cmdParts := strings.Split(cmd, " ")
		ex := exec.Command(cmdParts[0], cmdParts[1:]...)
		ex.Dir = wdPath

		stdoutPipe, _ := ex.StdoutPipe()
		stderrPipe, _ := ex.StderrPipe()

		storeDeployLog(infoLvl, fmt.Sprintf("Starting command execution: `%s`", cmd))
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
			storeDeployLog(errLvl, fmt.Sprintf("[X]: %s", stderrBuilder.String()))

			if stopOnErr {
				return runErr
			}

			storeDeployLog(warnLvl, "[!] Recovering...")

			paint.ErrorF("Error running command (%s): %v\n%s", cmd, runErr, "Continuing...")
		} else {
			storeDeployLog(okLvl, fmt.Sprintf("[✓]: %s", stdoutBuilder.String()))
		}
	}

	return nil
}

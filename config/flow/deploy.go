package flow

import (
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/Techzy-Programmer/d2m/config/db"
	"github.com/Techzy-Programmer/d2m/config/helpers"
	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/Techzy-Programmer/d2m/config/types"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// ToDo: Implememt comprehensive logging for deployment with timestamps and unique identifier

func StartDeployment(req *types.DeploymentRequest) {
	success := false

	deployment := &db.Deployment{
		ID:      helpers.GenerateSecure4DigitNumber(),
		StartAt: time.Now(),
		Branch:  req.Branch,
		Repo:    req.RepoPath,
	}

	defer saveDeployment(deployment, &success)

	homeDir, hdirErr := getUserHomeDirectory(req.LocalUser)
	if hdirErr != nil {
		paint.Error("Error getting user home directory: ", hdirErr)
		return
	}

	// Ensure the parent path exists
	parentPath := path.Join(homeDir, req.LocalParentPath)
	if info, err := os.Stat(parentPath); os.IsNotExist(err) || !info.IsDir() {
		paint.Error("Parent path is not a valid directory: ", parentPath)
		return
	}

	preExErr := execCmds(req.PreDeployCmds, parentPath, req.FailOnCmdError)
	if preExErr != nil {
		paint.Error("Error running pre-deployment commands: ", preExErr)
		return
	}

	// Let's fetch the repo from GitHub
	ghErr := ensureGHRepo(req.RepoPath, parentPath, 3, deployment)
	if ghErr != nil {
		paint.Error("Error fetching GitHub repository: ", ghErr)
		return
	}

	// ToDo: Implement AutoSetupDeps with smart inference

	// Run the deployment commands
	depExErr := execCmds(req.PostDeployCmds, parentPath, req.FailOnCmdError)
	if depExErr != nil {
		paint.Error("Error running post-deployment commands: ", depExErr)
		return
	}

	success = true
}

func ensureGHRepo(repoPth string, parentPath string, retry int, deployment *db.Deployment) error {
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
			return poErr
		}

		// Pull the latest changes
		paint.Info("Pulling latest changes for: ", appName)
		if wt, err := repo.Worktree(); err != nil {
			return err
		} else {
			pullErr := wt.Pull(&git.PullOptions{
				Auth: authOpt,
			})

			if pullErr != nil {
				return pullErr
			}
		}

		hash, msg, commitErr := getCommitData(repo)
		if commitErr != nil {
			return commitErr
		}

		deployment.CommitHash = hash
		deployment.CommitMsg = msg
		return nil
	}

	// Clone the repository
	paint.Info("Cloning repository: ", appName)
	repo, cloneErr := git.PlainClone(appPth, false, &git.CloneOptions{
		URL:  "https://github.com/" + repoPth,
		Auth: authOpt,
	})

	if cloneErr != nil {
		if retry > 0 {
			paint.Error("Error cloning repository: ", cloneErr)
			paint.Info("Retrying...")
			return ensureGHRepo(repoPth, parentPath, retry-1, deployment)
		}

		return cloneErr
	}

	hash, msg, commitErr := getCommitData(repo)
	if commitErr != nil {
		return commitErr
	}

	deployment.CommitHash = hash
	deployment.CommitMsg = msg
	return nil
}

func execCmds(cmds []string, wdPath string, stopOnErr bool) error {
	for _, cmd := range cmds {
		cmdParts := strings.Split(cmd, " ")
		ex := exec.Command(cmdParts[0], cmdParts[1:]...)
		ex.Dir = wdPath

		if runErr := ex.Run(); runErr != nil {
			if stopOnErr {
				return runErr
			}

			paint.ErrorF("Error running command (%s): %v\n%s", cmd, runErr, "Continuing...")
		}
	}

	return nil
}

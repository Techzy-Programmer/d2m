package flow

import (
	"errors"
	"os/user"
	"time"

	"github.com/Techzy-Programmer/d2m/config/db"
	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/go-git/go-git/v5"
)

func getCommitData(repo *git.Repository) (string, string, error) {
	ref, refErr := repo.Head()
	if refErr != nil {
		return "", "", refErr
	}

	commit, commitErr := repo.CommitObject(ref.Hash())
	if commitErr != nil {
		return "", "", commitErr
	}

	paint.Info("Latest commit: ", commit.Hash.String(), " by ", commit.Author.Name)
	return commit.Hash.String(), commit.Message, nil
}

func getUserHomeDirectory(username string) (string, error) {
	usr, err := user.Lookup(username)
	if err != nil {
		return "", errors.New("User not found: " + username)
	}

	// Return the home directory
	return usr.HomeDir, nil
}

func saveDeployment(deploy *db.Deployment, success *bool) {
	deploy.Logs = db.DeployLogStorage[deploy.ID]
	deploy.EndAt = time.Now()

	if *success {
		deploy.Status = "success"
	} else {
		deploy.Status = "failed"
	}

	db.SaveDeployment(deploy)
	db.FlushDeployLogs(deploy.ID)
}

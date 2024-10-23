package flow

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/Techzy-Programmer/d2m/config/db"
	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func ensureRepo(repoPth string, parentPath string, logger *deployLogHandler) (string, string, error) {
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
			logger.logErr("Error opening local repository at path: " + appPth)
			return "", "", poErr
		}

		// Pull the latest changes
		logger.logStep("Local repository available, pulling latest changes...")
		paint.Info("Pulling latest changes for: ", appName)
		hash, msg, pullErr := pullRepo(repo, authOpt)
		if pullErr != nil {
			logger.logErr("[X] Pull failed with error: " + pullErr.Error())
			return "", "", pullErr
		}

		logger.logOk(fmt.Sprintf("[✓] Pull successful, commit hash: %s", hash))
		return hash, msg, nil
	}

	// Clone the repository
	paint.Info("Cloning repository: ", appName)
	logger.logStep("Repository does not exist locally, cloning to path: " + appPth)
	hash, msg, cloneErr := cloneRepo("https://github.com/"+repoPth+".git", appPth, authOpt)
	if cloneErr != nil {
		logger.logErr("[X] Clone failed with error: " + cloneErr.Error())
		return "", "", cloneErr
	}

	logger.logOk(fmt.Sprintf("[✓] Clone successful, commit hash: %s", hash))
	return hash, msg, nil
}

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

func cloneRepo(url string, path string, auth transport.AuthMethod) (string, string, error) {
	repo, cloneErr := git.PlainClone(path, false, &git.CloneOptions{
		URL:  url,
		Auth: auth,
	})

	if cloneErr != nil {
		return "", "", cloneErr
	}

	hash, msg, commitErr := getCommitData(repo)
	if commitErr != nil {
		return "", "", commitErr
	}

	return hash, msg, nil
}

func pullRepo(repo *git.Repository, auth transport.AuthMethod) (string, string, error) {
	wt, err := repo.Worktree()
	if err != nil {
		return "", "", err
	}

	hash, msg, commitErr := getCommitData(repo)
	if commitErr != nil {
		return "", "", commitErr
	}

	pullErr := wt.Pull(&git.PullOptions{
		Auth: auth,
	})

	if pullErr != nil {
		if pullErr == git.NoErrAlreadyUpToDate {
			return hash, msg, nil
		}

		return "", "", pullErr
	}

	return hash, msg, nil
}

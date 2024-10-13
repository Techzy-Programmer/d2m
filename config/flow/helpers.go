package flow

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"os/user"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/Techzy-Programmer/d2m/config/db"
	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/Techzy-Programmer/d2m/config/types"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

var currentDeployID uint
var deployLogStorage []db.DeploymentLog
var deploymentRunning bool = false
var deployQueue []*types.DeploymentRequest

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
	deployLogStorage = []db.DeploymentLog{}
}

func ensureGHRepo(repoPth string, parentPath string, logger *deployLogHandler) (string, string, error) {
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

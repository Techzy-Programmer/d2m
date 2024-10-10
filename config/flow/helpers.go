package flow

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os/user"
	"time"

	"github.com/Techzy-Programmer/d2m/config/db"
	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/Techzy-Programmer/d2m/config/types"
	"github.com/go-git/go-git/v5"
)

var currentDeployID uint
var deployLogStorage []db.Log
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

func storeDeployLog(lvl uint, title string, msg string) {
	deployLogStorage = append(deployLogStorage, db.Log{
		Timestamp: time.Now().Unix(),
		DeployID:  &currentDeployID,
		Title:     title,
		Message:   msg,
		Level:     lvl,
	})
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

func getUserHomeDirectory(username string) (string, error) {
	usr, err := user.Lookup(username)
	if err != nil {
		return "", errors.New("User not found: " + username)
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

	storeDeployLog(infoLvl, "Done", fmt.Sprintf("Deployment completed with status: %s", deploy.Status))

	deploy.Logs = deployLogStorage
	deploy.EndAt = time.Now()

	db.SaveDeployment(deploy)
	deployLogStorage = []db.Log{}
}

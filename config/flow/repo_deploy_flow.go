package flow

import (
	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/Techzy-Programmer/d2m/config/types"
)

func StartRepoDeployment(req *types.DeploymentRequest, ghSt *types.RepoDeploymentStrategy) {
	if deploymentRunning {
		deployQueue = append(deployQueue, req)
		return
	}

	dep, logger, parentPath, derr := preInit(req)
	if derr != nil {
		return
	}

	success := false
	dep.Branch = ghSt.Branch
	dep.Repo = ghSt.Path
	defer saveDeployment(dep, &success)

	if len(req.PreDeployCmds) > 0 {
		logger.reset("Pre-Deploy Commands", "Starting execution...")

		preExErr := execCmds(req.PreDeployCmds, parentPath, req.FailOnError, logger)
		if preExErr != nil {
			logger.logErr("[X] Fatal error running pre-deployment commands").save(errLvl)
			paint.Error("Error running pre-deployment commands: ", preExErr)
			return
		}

		logger.save(okLvl)
	}

	logger.reset("Repository Fetch", "")

	// Let's fetch the repo from GitHub
	hash, msg, ghErr := ensureRepo(ghSt.Path, parentPath, logger)
	if ghErr != nil {
		logger.logErr("[X] Fatal error fetching remote repository").save(errLvl)
		paint.Error("Error fetching GitHub repository: ", ghErr)
		return
	}

	dep.CommitMsg = msg
	dep.CommitHash = hash
	logger.save(okLvl)

	// ToDo: Implement AutoSetupDeps with smart inference

	if len(req.PostDeployCmds) > 0 {
		logger.reset("Post-Deploy Commands", "Starting execution...")

		postExErr := execCmds(req.PostDeployCmds, parentPath, req.FailOnError, logger)
		if postExErr != nil {
			logger.logErr("[X] Fatal error running post-deployment commands").save(errLvl)
			paint.Error("Error running post-deployment commands: ", postExErr)
			return
		}

		logger.save(okLvl)
	}

	success = true
}

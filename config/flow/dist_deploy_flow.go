package flow

import (
	"github.com/Techzy-Programmer/d2m/config/helpers"
	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/Techzy-Programmer/d2m/config/types"
)

func StartDistDeployment(req *types.DeploymentRequest, dstSt *types.DistDeploymentStrategy) {
	if deploymentRunning {
		deployQueue = append(deployQueue, req)
		return
	}

	dep, logger, parentPath, derr := preInit(req)
	if derr != nil {
		return
	}

	success := false
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

	helpers.UnzipFolder(dstSt.FileName, parentPath, dstSt.RootFolderName)

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

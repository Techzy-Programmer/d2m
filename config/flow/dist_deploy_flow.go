package flow

import (
	"github.com/Techzy-Programmer/d2m/config/types"
)

func StartDistDeployment(req *types.DeploymentRequest, ghSt *types.DistDeploymentStrategy) {
	if deploymentRunning {
		deployQueue = append(deployQueue, req)
		return
	}
}

package main

import (
	// "os"

	// "github.com/Techzy-Programmer/d2m/config/flow"
	// "github.com/Techzy-Programmer/d2m/config/types"
	"github.com/Techzy-Programmer/d2m/config/paint"
)

func startDebug() {
	paint.Warn("Running in debug mode")

	// flow.StartDeployment(&types.DeploymentRequest{
	// 	Branch: 				 "main",
	// 	AutoSetupDeps:   true,
	// 	LocalUser:       "risha",
	// 	LocalParentPath: "Documents\\test\\d2m-deploy",
	// 	RepoPath:        "Node-Bug-Hunter/Hunter-UI",
	// 	PreDeployCmds:   []string{"mkdir hello-world", "touch hello-world/index.html"},
	// 	PostDeployCmds:  []string{"touch hello-world/deploy-success.txt"},
	// 	FailOnCmdError:  false,
	// })

	// os.Exit(0)
}

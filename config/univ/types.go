package univ

type DeploymentRequest struct {
	PostDeployCmds  []string `json:"post_deploy_cmds"`  // Default: []
	PreDeployCmds   []string `json:"pre_deploy_cmds"`   // Default: []
	// If true, the deployment will fail if any of the commands fail
	FailOnCmdError  bool     `json:"fail_on_cmd_error"` // Default: true
	// Based on project type (Supported: Node.js, Python & golang)
	// dependencies will be installed before running the commands
	AutoSetupDeps   bool     `json:"auto_setup_deps"`   // Default: true
	// On behalf of which user the deployment is to be done
	LocalUser       string   `json:"local_user"`
	// This is the path where the repository will be cloned and commands will be executed relative to this path as cwd
	LocalPath       string   `json:"local_path"`
	// Format: GH_USER_NAME/REPO_NAME
	RepoPath        string   `json:"repo_path"`         // Format: GH_USER_NAME/REPO_NAME
	Branch          string   `json:"branch"`            // Default: main
}

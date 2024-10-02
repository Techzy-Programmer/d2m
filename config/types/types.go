package types

// [-] Note: Put up type definition here which are used across multiple packages

type DeploymentRequest struct {
	PostDeployCmds  []string `json:"post_deploy_cmds"`  // Commands to run after deployment (default: [])
	PreDeployCmds   []string `json:"pre_deploy_cmds"`   // Commands to run before deployment (default: [])
	FailOnCmdError  bool     `json:"fail_on_cmd_error"` // Stop deployment if any of the command fails (default: true)
	AutoSetupDeps   bool     `json:"auto_setup_deps"`   // Automatically install dependencies based on project type (default: true)
	LocalUser       string   `json:"local_user"`        // User to execute deployment [LocalParentPath is relative to user's home-path] (default: root)
	LocalParentPath string   `json:"local_parent_path"` // Path where the repository will be cloned and post-deployment commands executed (relative to user's home directory)
	RepoPath        string   `json:"repo_path"`         // Repository path in the (format: {GH_USER_NAME}/{REPO_NAME})
	Branch          string   `json:"branch"`            // Branch to deploy (default: main)
}

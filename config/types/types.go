package types

import (
	"encoding/json"
)

// [-] Note: Put up type definition here which are used across multiple packages

type Strategy interface{}

type RepoDeploymentStrategy struct {
	Path          string `json:"path"`            // Repository path in the (format: {GH_USER_NAME}/{REPO_NAME})
	Branch        string `json:"branch"`          // Branch to deploy (default: main)
	AutoSetupDeps bool   `json:"auto_setup_deps"` // Automatically install dependencies based on project type (default: true)
	// ToDo: Activate following fields
	// Origin          string   `json:"origin"`            // Origin URL (default: https://github.com)
	// IsPrivate       bool     `json:"is_private"`        // Is the repository private (default: false)
	// Username        string   `json:"username"`          // Username for the repository (default: "")
	// AccessToken     string   `json:"access_token"`      // Access token for the repository (default: "")
}

type DistDeploymentStrategy struct {
	FileName string `json:"file_name"` // Name of the file just uploaded
}

type EmptyDeploymentStrategy struct{}

type DockerDeploymentStrategy struct {
	// ToDo: Add Docker deployment strategy fields
}

type DeploymentRequest struct {
	StrategyType    string          `json:"strategy_type"`     // Deployment strategy type (default: github)
	Strategy        json.RawMessage `json:"strategy"`          // Deployment strategy (default: EmptyDeploymentStrategy)
	PostDeployCmds  []string        `json:"post_deploy_cmds"`  // Commands to run after deployment (default: [])
	PreDeployCmds   []string        `json:"pre_deploy_cmds"`   // Commands to run before deployment (default: [])
	FailOnError     bool            `json:"fail_on_error"`     // Stop deployment if any of the command fails (default: true)
	LocalUser       string          `json:"local_user"`        // User to execute deployment [LocalParentPath is relative to user's home-path] (default: root)
	LocalParentPath string          `json:"local_parent_path"` // Path where the repository will be cloned and post-deployment commands executed (relative to user's home directory)
}

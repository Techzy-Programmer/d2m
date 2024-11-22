package handler

import (
	"io"
	"slices"

	"github.com/Techzy-Programmer/d2m/config/flow"
	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/Techzy-Programmer/d2m/config/types"
	"github.com/Techzy-Programmer/d2m/config/vars"
	"github.com/gin-gonic/gin"
)

func HandleDeployment(c *gin.Context) {
	// Following check ensures that request only gets through if coming from GitHub actions workflow
	// ToDo: Support custom CD servers and/or workflow runners
	allowedIps := slices.Concat(vars.GHActionIps, vars.LocalIPs)
	if !slices.Contains(allowedIps, c.ClientIP()) {
		c.JSON(403, gin.H{
			"message": "You are not allowed to access this resource",
			"ok":      false,
		})
		return
	}

	body, readErr := io.ReadAll(c.Request.Body)
	if readErr != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request body",
			"ok":      false,
		})
		return
	}

	defer c.Request.Body.Close()

	// Deserialize the decrypted request body
	req, strategy, jsonErr := unmarshalDeploymentRequest(body)
	if jsonErr != nil {
		paint.Error("Error deserializing request body: ", jsonErr)
		c.JSON(400, gin.H{
			"message": "Invalid request body",
			"ok":      false,
		})
		return
	}

	// ToDo: Make it concurrent
	switch s := strategy.(type) {
	case *types.RepoDeploymentStrategy:
		flow.StartRepoDeployment(req, s)

	case *types.DistDeploymentStrategy:
		flow.StartDistDeployment(req, s)

	case *types.EmptyDeploymentStrategy:
		flow.StartEmptyDeployment(req, s)

	case *types.DockerDeploymentStrategy:
		// ToDo: Implement Docker deployment flow
	}

	c.JSON(200, gin.H{
		"message": "Deployment triggered successfully",
		"ok":      true,
	})
}

package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Techzy-Programmer/d2m/config/types"
	"github.com/golang-jwt/jwt/v5"
)

func getRelativeDuration(startTime int64) string {
	now := time.Now().Unix()
	duration := now - startTime

	seconds := duration
	minutes := seconds / 60
	hours := minutes / 60
	days := hours / 24
	weeks := days / 7
	months := days / 30
	years := days / 365

	switch {
	case years > 0:
		return fmt.Sprintf("%d years", years)
	case months > 0:
		return fmt.Sprintf("%d months", months)
	case weeks > 0:
		return fmt.Sprintf("%d weeks", weeks)
	case days > 0:
		return fmt.Sprintf("%d days", days)
	case hours > 0:
		return fmt.Sprintf("%d hours", hours)
	case minutes > 0:
		return fmt.Sprintf("%d minutes", minutes)
	default:
		return fmt.Sprintf("%d seconds", seconds)
	}
}

func generateJWTToken(secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.Claims(jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "d2m-cli",
	}))

	return token.SignedString([]byte(secret))
}

func verifyJWTToken(tokenString, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func unmarshalDeploymentRequest(data []byte) (*types.DeploymentRequest, interface{}, error) {
	var aux types.DeploymentRequest
	if err := json.Unmarshal(data, &aux); err != nil {
		return nil, nil, errors.New("failed to unmarshal deployment request")
	}

	switch aux.StrategyType {
	case "repo":
		var repoStrategy = &types.RepoDeploymentStrategy{}
		if err := json.Unmarshal(aux.Strategy, &repoStrategy); err != nil {
			return nil, nil, errors.New("failed to unmarshal repo strategy")
		}

		return &aux, repoStrategy, nil

	case "dist":
		var distStrategy = &types.DistDeploymentStrategy{}
		if err := json.Unmarshal(aux.Strategy, &distStrategy); err != nil {
			return nil, nil, errors.New("failed to unmarshal dist strategy")
		}

		return &aux, distStrategy, nil

	case "docker":
		var dockerStrategy = &types.DockerDeploymentStrategy{}
		if err := json.Unmarshal(aux.Strategy, &dockerStrategy); err != nil {
			return nil, nil, errors.New("failed to unmarshal docker strategy")
		}

		return &aux, dockerStrategy, nil

	case "empty", "":
		return &aux, &types.EmptyDeploymentStrategy{}, nil
	}

	return nil, nil, errors.New("invalid deployment strategy type")
}

package db

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

func getUserConfigPath(appName string) (string, error) {
	var configPath string

	switch runtime.GOOS {
	case "linux":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configPath = filepath.Join(homeDir, ".config", appName)

	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configPath = filepath.Join(homeDir, "Library", "Application Support", appName)

	case "windows":
		localAppData := os.Getenv("LocalAppData")
		if localAppData == "" {
			return "", errors.New("LocalAppData environment variable is not set")
		}
		configPath = filepath.Join(localAppData, appName)

	default:
		return "", errors.New("unsupported operating system: " + runtime.GOOS)
	}

	return configPath, nil
}

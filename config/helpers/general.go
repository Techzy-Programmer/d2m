package helpers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/Techzy-Programmer/d2m/config/vars"
)

func GetUserConfigPath(appName string) (string, error) {
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

func ScheduleGHActionIPFetch() {
	tries := 0

	for {
		ips, err := getActionsIP()
		if err != nil {
			time.Sleep(time.Duration(tries+1) * 10 * time.Second) // Linear backoff

			if tries < 3 {
				tries++
				continue
			}

			// Network error
			paint.Error("Failed to fetch GitHub Actions IPs:", err)
		}

		tries = 0
		vars.GHActionIps = ips
		time.Sleep(6 * time.Hour)
	}
}

type GitHubMeta struct {
	Actions []string `json:"actions"`
}

func BodyAsText(req *http.Request) string {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return ""
	}

	return string(body)
}

func GetRelativeDuration(startTime int64) string {
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

package helpers

import (
	"errors"
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

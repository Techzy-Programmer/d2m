package univ

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Techzy-Programmer/d2m/config/paint"
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
		GHActionIps = ips
		time.Sleep(6 * time.Hour)
	}
}

type GitHubMeta struct {
	Actions []string `json:"actions"`
}

func getActionsIP() ([]string, error) {
	res, err := http.Get("https://api.github.com/meta")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var meta GitHubMeta
	err = json.Unmarshal(body, &meta)
	if err != nil {
		return nil, err
	}

	return meta.Actions, nil
}

func GetPrivateKey(privText string) (*rsa.PrivateKey, error) {
	// Decode the PEM block
	block, _ := pem.Decode([]byte(privText))
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	// Parse the private key
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("error parsing private key: " + (err.Error()))
	}

	// Type assertion to convert to *rsa.PrivateKey
	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not an RSA private key")
	}

	return rsaPrivateKey, nil
}

func RSADecryptWithPrivateKey(encryptedData string, privateKey *rsa.PrivateKey) (string, error) {
	decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, []byte(encryptedData))
	if err != nil {
		return "", errors.New("error decrypting data: " + err.Error())
	}

	return string(decryptedData), nil
}

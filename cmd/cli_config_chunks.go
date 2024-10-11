package cmd

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/Techzy-Programmer/d2m/config/helpers"
	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/erikgeiser/promptkit/textinput"
	"golang.org/x/crypto/bcrypt"
)

func getWebPort() (string, error) {
	paint.Notice("\nPlease make sure the port you provide is not in use by any other service.")
	paint.Notice("D2M will use this port to serve the Panel UI and webhooks API.")
	paint.Notice("Make sure the port is accessible from the public internet.")
	webPortIn := textinput.New("Configure a port for Web Server: ")
	webPortIn.InitialValue = "8080"
	webPortIn.Placeholder = "8080"

	webPortIn.Validate = func(input string) error {
		if input == "" {
			return errors.New("port number is required")
		}

		port, err := strconv.Atoi(input)
		if err != nil {
			return errors.New("port number must be a valid integer")
		}

		if port < 80 || port > 65535 {
			return errors.New("invalid port number range")
		}

		return nil
	}

	webPort, err := webPortIn.RunPrompt()
	if err != nil {
		return "", err
	}

	return webPort, nil
}

func getAccessPwd() ([]byte, error) {
	paint.Notice("\nYou'll be required to enter this access password at time of web panel login.")
	paint.Notice("You also have to set it up in your GitHub Actions runner.")
	accessPwdIn := textinput.New("Set a new Access Password: ")
	accessPwdIn.Placeholder = "$0mEThiNg$TR0ng&S3cU#e"
	accessPwdIn.Hidden = true

	accessPwdIn.Validate = func(input string) error {
		if !isValidPassword(input) {
			return errors.New("password is too weak")
		}

		return nil
	}

	accessPwd, err := accessPwdIn.RunPrompt()
	if err != nil {
		return nil, err
	}

	cryptPwdBytes, cryptErr := bcrypt.GenerateFromPassword([]byte(accessPwd), bcrypt.DefaultCost)
	if cryptErr != nil {
		return nil, cryptErr
	}

	return cryptPwdBytes, nil
}

func getGHPat() (string, error) {
	paint.Notice("\nTo clone private repositories, D2M requires a GitHub Personal Access Token (PAT).")
	paint.Notice("Leave this field empty if you don't want to include private repos in your pipeline.")
	ghPATIn := textinput.New("Enter GitHub Personal Access Token (with repo clone permission): ")
	gitRegex := `^(gh[ps]_[a-zA-Z0-9]{36}|github_pat_[a-zA-Z0-9]{22}_[a-zA-Z0-9]{59}|)$`
	ghPATIn.Validate = createRegexValidator(gitRegex, "Invalid GitHub PAT")
	ghPATIn.Placeholder = "ghp_XXXXXXXXXXXXXX"
	ghPATIn.Hidden = true

	ghPAT, err := ghPATIn.RunPrompt()
	if err != nil {
		return "", err
	}

	return ghPAT, nil
}

func getGHUsername() (string, error) {
	paint.Notice("\nGitHub username is required to fetch repositories with PAT.")
	ghUsernameIn := textinput.New("Enter your GitHub Username: ")
	ghUsernameIn.Placeholder = "github-username"
	var ghErr error

	ghUsername, ghErr := ghUsernameIn.RunPrompt()
	if ghErr != nil {
		return "", ghErr
	}

	return ghUsername, nil
}

func getPrivateKey() ([]byte, error) {
	paint.Notice("\nD2M requires private key corresponding to public key of your GH-Actions runner to decrypt the webhook payload")
	privDecryptKeyIn := textinput.New("Enter the path to your private key: ")
	privDecryptKeyIn.Placeholder = "path/to/private-key.pem"

	privDecryptKeyIn.Validate = func(input string) error {
		if input == "" || !strings.HasSuffix(input, ".pem") {
			return errors.New("private key path is required")
		}

		if _, err := os.Stat(input); os.IsNotExist(err) {
			return errors.New("private key file not found")
		}

		privKey, err := os.ReadFile(input)
		if err != nil {
			return err
		}

		_, err = helpers.GetPrivateKey(string(privKey))
		if err != nil {
			return errors.New("invalid private key")
		}

		return nil
	}

	privKeyPathIn, err := privDecryptKeyIn.RunPrompt()
	if err != nil {
		return nil, err
	}

	privKey, err := os.ReadFile(privKeyPathIn)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}

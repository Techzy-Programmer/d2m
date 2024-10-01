package cmd

import (
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Techzy-Programmer/d2m/config/db"
	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/erikgeiser/promptkit/textinput"
	"github.com/urfave/cli/v2"
)

func HandleInitCMD(*cli.Context) error {
	hasConfig := db.GetConfig[bool]("user.HasConfig")

	if !hasConfig {
		requestConfig()
	} else {
		getHelp()
	}

	return nil
}

func requestConfig() {
	paint.Info("D2M is not configured yet.\nPlease provide the following details to setup D2M.\n")

	portValidator := func(input string) error {
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

	paint.Notice("\nPlease make sure the port you provide is not in use by any other service.")
	paint.Notice("D2M will use this port to serve the Panel UI and webhooks API.")
	paint.Notice("Make sure the port is accessible from the public internet.")
	webPortIn := textinput.New("Configure a port for Web Server: ")
	webPortIn.Validate = portValidator
	webPortIn.InitialValue = "8080"
	webPortIn.Placeholder = "8080"

	webPort, err := webPortIn.RunPrompt()
	if err != nil {
		paint.Error("Error: ", err)
		return
	}

	ghUsernameIn := textinput.New("Enter your GitHub Username: ")
	ghUsernameIn.Placeholder = "github-username"

	ghUsername, err := ghUsernameIn.RunPrompt()
	if err != nil {
		paint.Error("Error: ", err)
		return
	}

	paint.Notice("\nTo clone private repositories, D2M requires a GitHub Personal Access Token (PAT).")
	paint.Notice("Leave this field empty if you don't want to include private repos in your pipeline.")
	ghPATIn := textinput.New("Enter GitHub Personal Access Token (with repo clone permission): ")
	ghPATIn.Placeholder = "ghp_XXXXXXXXXXXXXX"
	ghPATIn.Hidden = true

	ghPATIn.Validate = func(input string) error {
		if input == "" {
			return nil
		}

		pattern := `^(gh[ps]_[a-zA-Z0-9]{36}|github_pat_[a-zA-Z0-9]{22}_[a-zA-Z0-9]{59})$`

		regex, err := regexp.Compile(pattern)
		if err != nil {
			return err
		}

		if !regex.MatchString(input) {
			return errors.New("invalid GitHub PAT")
		}

		return nil
	}

	ghPAT, err := ghPATIn.RunPrompt()
	if err != nil {
		paint.Error("Error: ", err)
		return
	}

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

		return nil
	}

	privKeyPathIn, err := privDecryptKeyIn.RunPrompt()
	if err != nil {
		paint.Error("Error: ", err)
		return
	}

	privKey, err := os.ReadFile(privKeyPathIn)
	if err != nil {
		paint.Error("Error: ", err)
		return
	}

	db.SetConfig("user.HasConfig", true)
	db.SetConfig("user.WebPort", webPort)
	db.SetConfig("user.GHUsername", ghUsername)
	db.SetConfig("user.GHPAT", ghPAT)
	db.SetConfig("user.PrivateKey", string(privKey))
}

func getHelp() {
	paint.Info("D2M is configured and running. Please execute 'd2m h' to see available commands.")
}

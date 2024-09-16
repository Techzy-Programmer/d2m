package cmd

import (
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Techzy-Programmer/d2m/config"
	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/erikgeiser/promptkit/textinput"
	"github.com/urfave/cli/v2"
)

func HandleInitCMD(*cli.Context) error {
	hasConfig := config.GetData[bool]("user.HasConfig")

	if !hasConfig {
		requestConfig()
	} else {
		askForService()
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

	uiPortIn := textinput.New("Enter port for web panel: ")
	uiPortIn.Validate = portValidator
	uiPortIn.InitialValue = "8000"
	uiPortIn.Placeholder = "8000"
	
	uiPort, err := uiPortIn.RunPrompt()
	if err != nil {
		paint.Error("Error: ", err)
		return;
	}

	apiPortIn := textinput.New("Enter port for webhook API: ")
	apiPortIn.Validate = portValidator
	apiPortIn.InitialValue = "8080"
	apiPortIn.Placeholder = "8080"

	apiPort, err := apiPortIn.RunPrompt()
	if err != nil {
		paint.Error("Error: ", err)
		return;
	}

	ghUsernameIn := textinput.New("Enter your GitHub Username: ")
	ghUsernameIn.Placeholder = "github-username"

	ghUsername, err := ghUsernameIn.RunPrompt()
	if err != nil {
		paint.Error("Error: ", err)
		return;
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
		return;
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

	config.SetData("user.HasConfig", true)
	config.SetData("user.UIPort", uiPort)
	config.SetData("user.APIPort", apiPort)
	config.SetData("user.GHUsername", ghUsername)
	config.SetData("user.GHPAT", ghPAT)
	config.SetData("user.PrivateKey", string(privKey))
}

func askForService() {
}

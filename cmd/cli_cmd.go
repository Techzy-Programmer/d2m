package cmd

import (
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/Techzy-Programmer/d2m/config/db"
	"github.com/Techzy-Programmer/d2m/config/helpers"
	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/erikgeiser/promptkit/textinput"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/bcrypt"
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
	paint.Info("D2M is not configured yet.\nPlease provide the following details to setup D2M.")
	gitRegex := `^(gh[ps]_[a-zA-Z0-9]{36}|github_pat_[a-zA-Z0-9]{22}_[a-zA-Z0-9]{59}|)$`

	createRegexValidator := func(pattern string, emsg string) func(input string) error {
		return func(input string) error {
			regex, err := regexp.Compile(pattern)
			if err != nil {
				return err
			}

			if !regex.MatchString(input) {
				return errors.New(emsg)
			}

			return nil
		}
	}

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
		paint.Error("Error: ", err)
		return
	}

	paint.Notice("\nTo clone private repositories, D2M requires a GitHub Personal Access Token (PAT).")
	paint.Notice("Leave this field empty if you don't want to include private repos in your pipeline.")
	ghPATIn := textinput.New("Enter GitHub Personal Access Token (with repo clone permission): ")
	ghPATIn.Validate = createRegexValidator(gitRegex, "Invalid GitHub PAT")
	ghPATIn.Placeholder = "ghp_XXXXXXXXXXXXXX"
	ghPATIn.Hidden = true

	ghPAT, err := ghPATIn.RunPrompt()
	if err != nil {
		paint.Error("Error: ", err)
		return
	}

	var ghUsername string

	if ghPAT != "" {
		paint.Notice("\nGitHub username is required to fetch repositories with PAT.")
		ghUsernameIn := textinput.New("Enter your GitHub Username: ")
		ghUsernameIn.Placeholder = "github-username"
		var ghErr error

		ghUsername, ghErr = ghUsernameIn.RunPrompt()
		if ghErr != nil {
			paint.Error("Error: ", err)
			return
		}
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

		privKey, err := os.ReadFile(input)
		if err != nil {
			paint.Error("Error: ", err)
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
		paint.Error("Error: ", err)
		return
	}

	privKey, err := os.ReadFile(privKeyPathIn)
	if err != nil {
		paint.Error("Error: ", err)
		return
	}

	cryptPwdBytes, cryptErr := bcrypt.GenerateFromPassword([]byte(accessPwd), bcrypt.DefaultCost)
	if cryptErr != nil {
		paint.Error("Error: ", cryptErr)
		return
	}

	jwtSecret, jwtErr := generateSecureRandomString(32)
	if jwtErr != nil {
		paint.Error("Error: ", jwtErr)
		return
	}

	db.SetConfig("user.GHPAT", ghPAT)
	db.SetConfig("user.HasConfig", true)
	db.SetConfig("user.WebPort", webPort)
	db.SetConfig("app.JWTSecret", jwtSecret)
	db.SetConfig("user.GHUsername", ghUsername)
	db.SetConfig("user.PrivateKey", string(privKey))
	db.SetConfig("user.AccessPwd", string(cryptPwdBytes))
}

func getHelp() {
	paint.Info("D2M is configured and running. Please execute 'd2m h' to see available commands.")
}

func isValidPassword(password string) bool {
	if len(password) < 10 {
		return false
	}

	var (
		hasLower   = false
		hasUpper   = false
		hasDigit   = false
		hasSpecial = false
	)

	specialChars := "!@#$%^&*()_+-=[]{};\\':\"|,./<>?"
	var lastChar rune
	var repeatCount int

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		case strings.ContainsRune(specialChars, char):
			hasSpecial = true
		}

		if char == lastChar {
			repeatCount++
			if repeatCount == 3 {
				return false
			}
		} else {
			repeatCount = 1
		}
		lastChar = char
	}

	return hasLower && hasUpper && hasDigit && hasSpecial
}

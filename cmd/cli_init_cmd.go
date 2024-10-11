package cmd

import (
	"github.com/Techzy-Programmer/d2m/config/db"
	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/urfave/cli/v2"
)

// ToDo: Restart the daemon process after initialization

var InitCmd = &cli.Command{
	Name:    "init",
	Aliases: []string{"i"},
	Usage:   "Initialize & setup D2M instance on this machine",
	Action:  handleInitCMD,
}

func handleInitCMD(*cli.Context) error {
	hasConfig := db.GetConfig[bool]("user.HasConfig")

	if !hasConfig {
		requestConfig()
	} else {
		getHelp()
	}

	return nil
}

func requestConfig() {
	paint.Info("Please provide the following details to initialize D2M.")

	webPort, wpErr := getWebPort()
	if wpErr != nil {
		paint.Error("Error: ", wpErr)
		return
	}

	cryptPwdBytes, cryptErr := getAccessPwd()
	if cryptErr != nil {
		paint.Error("Error: ", cryptErr)
		return
	}

	ghPAT, ghErr := getGHPat()
	if ghErr != nil {
		paint.Error("Error: ", ghErr)
		return
	}

	var ghUsername string

	if ghPAT != "" {
		var ghUnErr error
		ghUsername, ghUnErr = getGHUsername()
		if ghUnErr != nil {
			paint.Error("Error: ", ghUnErr)
			return
		}
	}

	privKey, keyErr := getPrivateKey()
	if keyErr != nil {
		paint.Error("Error: ", keyErr)
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

	paint.Success("D2M initialized and configurations saved successfully!")
}

func getHelp() {
	paint.Warn("D2M has been initialized already! try `d2m h` for help.")
	paint.Info("Want to change configuarations? try `d2m uc h`")
}

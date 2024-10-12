package cmd

import (
	"errors"
	"os"

	"github.com/Techzy-Programmer/d2m/config/db"
	"github.com/Techzy-Programmer/d2m/config/msg"
	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/Techzy-Programmer/d2m/config/vars"
	"github.com/urfave/cli/v2"
)

var initialWP string

var updateWebPortCmd = &cli.Command{
	Name:    "web-port",
	Aliases: []string{"wp"},
	Usage:   "Update the web port for api and panel access",
	Action:  handleUpdateWebPort,
}

var updateAccessPwdCmd = &cli.Command{
	Name:    "access-pwd",
	Aliases: []string{"ap"},
	Usage:   "Update the access password for authentication",
	Action:  handleUpdateAccessPWD,
}

var updateGHPatCmd = &cli.Command{
	Name:    "gh-pat",
	Aliases: []string{"ghp"},
	Usage:   "Update the GitHub Personal Access Token (PAT) to include your private repositories in deployment pipeline",
	Action:  handleUpdateGHPat,
}

var updateGHUsernameCmd = &cli.Command{
	Name:    "gh-username",
	Aliases: []string{"ghu"},
	Usage:   "Update the GitHub Username, required to fetch repositories along with GitHub PAT",
	Action:  handleUpdateGHUsername,
}

var updatePrivateKeyCmd = &cli.Command{
	Name:    "private-key",
	Aliases: []string{"pk"},
	Usage:   "Update the private key for secure network calls",
	Action:  handleUpdatePrivateKey,
}

var UpdateCmd = &cli.Command{
	Name:    "update-conf",
	Aliases: []string{"uc"},
	Usage:   "Updates the configuration one parameter at a time",
	Before:  checkHasConfig,
	After:   handleConfigUpdate,
	Subcommands: []*cli.Command{
		updateAccessPwdCmd,
		updateWebPortCmd,
		updateGHPatCmd,
		updateGHUsernameCmd,
		updatePrivateKeyCmd,
	},
}

func checkHasConfig(*cli.Context) error {
	hasConfig := db.GetConfig("user.HasConfig", false)

	if !hasConfig {
		paint.Error("Error: D2M is not configured yet. Please run 'd2m i' first")
		return errors.New("config not found")
	}

	return nil
}

func handleConfigUpdate(c *cli.Context) error {
	msg.SendMsg(vars.CLIConn, msg.ConfigUpdateMSG{
		Type:  msg.ConfigUpdateMsgType,
		Which: c.Args().First(),
	})

	return nil
}

func handleUpdateAccessPWD(*cli.Context) error {
	cryptPwdBytes, cryptErr := getAccessPwd()
	if cryptErr != nil {
		paint.Error("Error: ", cryptErr)
		return nil
	}

	changeJWTSecret()
	db.SetConfig("user.AccessPwd", string(cryptPwdBytes))
	paint.Success("Access password updated successfully")
	return nil
}

func handleUpdateWebPort(*cli.Context) error {
	initialWP = db.GetConfig("user.WebPort", "8080")
	webPort, wpErr := getWebPort()
	if wpErr != nil {
		paint.Error("Error: ", wpErr)
		return nil
	}

	if initialWP == webPort {
		paint.Warn("Web port is already set to", webPort)
		paint.Notice("No changes made, exiting...")
		os.Exit(0)
		return nil
	}

	db.SetConfig("user.WebPort", webPort)
	paint.Success("Web port updated successfully")
	return nil
}

func handleUpdateGHPat(*cli.Context) error {
	ghPAT, ghErr := getGHPat()
	if ghErr != nil {
		paint.Error("Error: ", ghErr)
		return nil
	}

	db.SetConfig("user.GHPAT", ghPAT)
	paint.Success("GitHub Personal Access Token updated successfully")
	return nil
}

func handleUpdateGHUsername(*cli.Context) error {
	ghUsername, ghErr := getGHUsername()
	if ghErr != nil {
		paint.Error("Error: ", ghErr)
		return nil
	}

	db.SetConfig("user.GHUsername", ghUsername)
	paint.Success("GitHub username updated successfully")
	return nil
}

func handleUpdatePrivateKey(*cli.Context) error {
	privKey, keyErr := getPrivateKey()
	if keyErr != nil {
		paint.Error("Error: ", keyErr)
		return nil
	}

	changeJWTSecret()
	db.SetConfig("user.PrivateKey", string(privKey))
	paint.Success("Private key updated successfully")
	return nil
}

func changeJWTSecret() {
	jwtSecret, jwtErr := generateSecureRandomString(32)
	if jwtErr != nil {
		paint.Error("Error: ", jwtErr)
		return
	}

	db.SetConfig("app.JWTSecret", jwtSecret)
	paint.Warn("Any active session on the Web panel has been invalidated!")
}

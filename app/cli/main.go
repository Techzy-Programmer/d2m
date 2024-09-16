package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"syscall"
	"time"

	"github.com/Techzy-Programmer/d2m/app/cli/mcon"
	"github.com/Techzy-Programmer/d2m/app/daemon"
	"github.com/Techzy-Programmer/d2m/cmd"
	"github.com/Techzy-Programmer/d2m/config"
	"github.com/urfave/cli/v2"
)

func main() {
	daemonFlag := flag.Bool("daemon", false, "Run as daemon")
	flag.Parse()

	if *daemonFlag {
		daemon.LaunchDaemon()
		return
	}

	// Ensure daemon is running before handling CLI commands
	pid := config.GetData[float64]("daemon.PID")

	if !isProcessRunning(pid) {
		ensureDaemonRunning()
	}

	mcon.Initialize()

	app := &cli.App{
		Name:  "d2m",
		Usage: "Managr your deployments with ease",
		Action: cmd.HandleInitCMD,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	time.Sleep(60 * time.Second)
}

func isProcessRunning(pid float64) bool {
	// ToDo: Implement TCP based ping to check if daemon is running
	process, err := os.FindProcess(int(pid))

	if err != nil {
		return false
	}

	if runtime.GOOS == "windows" {
		return true
	}

	sigErr := process.Signal(syscall.Signal(0))
	return sigErr == nil
}

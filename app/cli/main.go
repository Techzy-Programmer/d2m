package main

import (
	"flag"
	"log"
	"net"
	"os"
	"runtime"
	"syscall"
	"time"

	"github.com/Techzy-Programmer/d2m/app/daemon"
	"github.com/Techzy-Programmer/d2m/cmd"
	"github.com/Techzy-Programmer/d2m/config/db"
	"github.com/Techzy-Programmer/d2m/config/vars"
	"github.com/Techzy-Programmer/d2m/internal/daemonizer"
	"github.com/Techzy-Programmer/d2m/internal/ipc"
	"github.com/urfave/cli/v2"
)

var Release string

func init() {
	vars.IsProd = (Release == "prod")
}

func main() {
	if !vars.IsProd {
		startDebug()
	}

	daemonFlag := flag.Bool("daemon", false, "Run as daemon")
	flag.Parse()

	if *daemonFlag {
		vars.IsDaemon = true
		daemon.LaunchDaemon()
		return
	}

	// Ensure daemon is running before handling CLI commands
	pid := db.GetConfig[float64]("daemon.PID")

	if !isProcessRunning(pid) {
		daemonizer.EnsureDaemonRunning()

		if !connectToDaemon() {
			panic("Unable to connect with daemon over TCP")
		}

		<-vars.AliveChannel // Wait for daemon to become responsive
	}

	app := &cli.App{
		Name:  "d2m",
		Usage: "continuous Delivery & Deployment Manager (D2M)\nManage your deployments with ease",
		Commands: []*cli.Command{
			cmd.InitCmd,
			cmd.UpdateCmd,
			cmd.LogsCmd,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	time.Sleep(500 * time.Millisecond)
}

func isProcessRunning(pid float64) bool {
	process, err := os.FindProcess(int(pid))

	if err != nil {
		return false
	}

	if runtime.GOOS == "windows" {
		return checkTCPAlive()
	}

	sigErr := process.Signal(syscall.Signal(0))
	if sigErr == nil {
		return checkTCPAlive()
	}

	return false
}

// Function to check if it's really our own daemon service that's running with the given PID
func checkTCPAlive() bool {
	if !connectToDaemon() {
		return false
	}

	select {
	case <-vars.AliveChannel:
		return true

	case <-time.After(5 * time.Second):
		vars.CLIConn.Close()
		return false
	}
}

func connectToDaemon() bool {
	port := db.GetConfig[string]("daemon.Port")
	if port == "" {
		return false
	}

	if vars.CLIConn != nil {
		vars.CLIConn.Close()
		vars.CLIConn = nil
	}

	conn, err := net.Dial("tcp", ":"+port)
	if err != nil {
		return false
	}

	go ipc.HandleConnection(conn)
	vars.CLIConn = conn
	return true
}

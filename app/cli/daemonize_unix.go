//go:build !windows

package main

// ! It's important to note that this is platform-specific code.
// ! This code will only run on Unix systems like Linux, macOS, Ubuntu, etc.

import (
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/Techzy-Programmer/d2m/config/paint"
)

func ensureDaemonRunning() {
	paint.Info("Daemonizing process...")
	executable, exErr := os.Executable()

	if exErr != nil {
		log.Fatalf("Failed to get executable path")
	}

	cmd := exec.Command("nohup", executable, "--daemon", "&")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	sErr := cmd.Start()

	if sErr != nil {
		log.Fatalf("Failed to start daemon: %v", sErr)
	}

	dErr := cmd.Process.Release()

	if dErr != nil {
		log.Fatalf("Failed to detach daemon process: %v", dErr)
	}

	paint.Success("Daemon started in background!\n")
	time.Sleep(4 * time.Second)
}

//go:build windows

package main

// ! It's important to note that this is platform-specific code.
// ! This code will only run on Windows systems.

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
	cmd := exec.Command(os.Args[0], "--daemon")
	cmd.Stdout = nil
	cmd.Stderr = nil

	cmd.SysProcAttr = &syscall.SysProcAttr{
		// https://learn.microsoft.com/en-us/windows/win32/procthread/process-creation-flags
		// CREATE_NEW_PROCESS_GROUP | DETACHED_PROCESS | CREATE_NO_WINDOW
		CreationFlags: 0x00000200 | 0x00000008 | 0x08000000,
		HideWindow:    true,
		ParentProcess: 0,
	}

	err := cmd.Start()

	if err != nil {
		log.Fatalf("Failed to start daemon: %v", err)
	}

	cmd.Process.Release()
	paint.Success("Daemon started in background!\n")
	time.Sleep(4 * time.Second)
}

package daemon

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Techzy-Programmer/d2m/config"
	"github.com/Techzy-Programmer/d2m/config/msg"
	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/Techzy-Programmer/d2m/config/univ"
	"github.com/Techzy-Programmer/d2m/internal/api"
	"github.com/Techzy-Programmer/d2m/internal/ipc"
	"github.com/Techzy-Programmer/d2m/internal/ui"
)

type daemonConfig struct {
	apiPort string
	uiPort  string
}

var dc daemonConfig

func init() {
	dc.apiPort = config.GetData("user.APIPort", "8080")
	dc.uiPort = config.GetData("user.UIPort", "8000")
}

func LaunchDaemon() {
	fmt.Println("Spinning up the daemon process...")

	go api.StartAPIServer(dc.apiPort)
	go ui.StartUIServer(dc.uiPort)
	go startDaemonTCPServer()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	<-sigChan // Wait for a signal
	config.SetData("daemon.PID", -1)
}

func startDaemonTCPServer() {
	listener, err := net.Listen("tcp", ":"+univ.DaemonPort)

	if err != nil {
		log.Fatalf("Failed to start daemon tcp server: %v", err)
	}

	defer listener.Close()
	paint.Info("Daemon server started at :" + univ.DaemonPort)
	config.SetData("daemon.PID", os.Getpid())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		msg.SendMsg(conn, msg.PingMSG{Type: msg.PingMsgType})
		go ipc.HandleConnection(conn)
	}
}

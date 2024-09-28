package daemon

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
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

	go univ.ScheduleGHActionIPFetch()
	go api.StartAPIServer(dc.apiPort)
	go ui.StartUIServer(dc.uiPort)
	go startDaemonTCPServer()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	<-sigChan // Wait for a signal
	config.SetData("daemon.PID", -1)
	config.SetData("daemon.Port", "")
}

func startDaemonTCPServer() {
	// Following will bind to a random port assigned by the OS
	listener, err := net.Listen("tcp", ":0")

	if err != nil {
		log.Fatalf("Failed to start daemon tcp server: %v", err)
	}

	defer listener.Close()
	addrSegs := strings.Split(listener.Addr().String(), ":")
	asgPort := addrSegs[len(addrSegs) - 1]
	paint.Info("Daemon server started at :" + asgPort)
	config.SetData("daemon.PID", os.Getpid())
	config.SetData("daemon.Port", asgPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		msg.SendMsg(conn, msg.PingMSG{Type: msg.PingMsgType, IsWelcome: true})
		go ipc.HandleConnection(conn)
	}
}

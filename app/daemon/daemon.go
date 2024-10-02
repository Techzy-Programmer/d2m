package daemon

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Techzy-Programmer/d2m/config/db"
	"github.com/Techzy-Programmer/d2m/config/helpers"
	"github.com/Techzy-Programmer/d2m/config/msg"
	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/Techzy-Programmer/d2m/config/vars"
	"github.com/Techzy-Programmer/d2m/internal/ipc"
	"github.com/Techzy-Programmer/d2m/internal/server"
	"github.com/gin-gonic/gin"
)

type daemonConfig struct {
	webPort string
}

var dc daemonConfig

// Synthetic init function
func synInit() {
	if vars.IsProd {
		gin.SetMode(gin.ReleaseMode)
	}

	keyStr := db.GetConfig("user.PrivateKey", "")
	if keyStr != "" {
		vars.PrivKey, _ = helpers.GetPrivateKey(keyStr)
	}

	dc.webPort = db.GetConfig("user.WebPort", "8080")
}

func LaunchDaemon() {
	synInit()
	fmt.Println("Spinning up the daemon process...")

	go helpers.ScheduleGHActionIPFetch()
	go server.StartWebServer(dc.webPort)
	go startDaemonTCPServer()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	<-sigChan // Wait for a signal
	db.DeleteConfig("daemon.PID")
	db.DeleteConfig("daemon.Port")
}

func startDaemonTCPServer() {
	// Following will bind to a random port assigned by the OS
	listener, err := net.Listen("tcp", ":0")

	if err != nil {
		log.Fatalf("Failed to start daemon tcp server: %v", err)
	}

	defer listener.Close()
	addrSegs := strings.Split(listener.Addr().String(), ":")
	asgPort := addrSegs[len(addrSegs)-1]
	paint.Info("Daemon server started at :" + asgPort)
	db.SetConfig("daemon.PID", os.Getpid())
	db.SetConfig("daemon.Port", asgPort)

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

package ipc

import (
	"bufio"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Techzy-Programmer/d2m/config/db"
	"github.com/Techzy-Programmer/d2m/config/msg"
	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/Techzy-Programmer/d2m/config/vars"
	"github.com/Techzy-Programmer/d2m/internal/daemonizer"
	"github.com/Techzy-Programmer/d2m/internal/server"
)

var locks = map[string]*sync.Mutex{
	"wp": {},
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		data, err := reader.ReadString('\n') // Read until newline, well it's a delimiter
		if err != nil {
			// "Error reading from connection:", err
			return
		}

		data = strings.TrimSpace(data)
		msg := msg.DeserializeMSG([]byte(data))
		go processMsg(msg, conn) // Fix: Process message in a separate goroutine to keep receiver thread free and avoid deadlocks
	}
}

func processMsg(message msg.MSG, conn net.Conn) {
	switch m := message.(type) {
	case *msg.PingMSG:
		if m.IsWelcome {
			// Notify CLI thread that daemon is alive
			vars.AliveChannel <- true
			close(vars.AliveChannel)
		}

		time.Sleep(10 * time.Second)
		msg.SendMsg(conn, msg.PingMSG{Type: msg.PingMsgType})
		var _ = m.Type

	case *msg.HaltMSG:
		if m.Ack { // Handled by CLI process
			time.Sleep(500 * time.Millisecond)
			vars.CLIConn = nil // Reset CLI connection as daemon has died
			daemonizer.EnsureDaemonRunning()
			return
		}

		msg.SendMsg(conn, msg.HaltMSG{Type: msg.HaltMsgType, Ack: true})
		time.Sleep(1 * time.Second) // Let daemon chill before dying ;o
		os.Exit(0)

	case *msg.ConfigUpdateMSG:
		switch m.Which {
		case "web-port", "wp":
			locks["wp"].Lock()
			defer locks["wp"].Unlock()
			server.StopWebServer <- true
			webPort := db.GetConfig("user.WebPort", "8080")
			time.Sleep(2 * time.Second) // Account for server shutdown
			server.StartWebServer(webPort)
		}

	default:
		paint.Error("Unknown message received")
	}
}

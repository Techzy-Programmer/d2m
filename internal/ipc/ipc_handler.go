package ipc

import (
	"bufio"
	"net"
	"strings"
	"time"

	"github.com/Techzy-Programmer/d2m/config/msg"
	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/Techzy-Programmer/d2m/config/vars"
)

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
		processMsg(msg, conn)
	}
}

func processMsg(message msg.MSG, conn net.Conn) {
	switch m := message.(type) {
	case *msg.PingMSG:
		if m.IsWelcome {
			// Notify CLI thread that daemon is alive
			vars.AliveChannel <- true
			close(vars.AliveChannel)
			return
		}

		time.Sleep(10 * time.Second)
		msg.SendMsg(conn, msg.PingMSG{Type: msg.PingMsgType})
		var _ = m.Type

	default:
		paint.Error("Unknown message received")
	}
}

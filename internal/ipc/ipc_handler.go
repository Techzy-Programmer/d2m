package ipc

import (
	"bufio"
	"net"
	"strings"
	"time"

	"github.com/Techzy-Programmer/d2m/config/msg"
	"github.com/Techzy-Programmer/d2m/config/paint"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		data, err := reader.ReadString('\n') // Read until newline, well it's a delimiter
		if err != nil {
			paint.Error("Error reading from connection:", err)
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
		paint.Info("Received Ping")
		time.Sleep(10 * time.Second)
		n := msg.SendMsg(conn, msg.PingMSG{Type: msg.PingMsgType})
		paint.SuccessF("Sent > %d bytes : Type > "+m.Type, n)

	default:
		paint.Error("Unknown message received")
	}
}

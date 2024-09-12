package mcon

import (
	"log"
	"net"

	"github.com/Techzy-Programmer/d2m/config/univ"
	"github.com/Techzy-Programmer/d2m/internal/ipc"
)

var CLIConn net.Conn

func Initialize() {
	var err error
	CLIConn, err = net.Dial("tcp", ":"+univ.DaemonPort)

	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}

	go ipc.HandleConnection(CLIConn)
}

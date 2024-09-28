package univ

import "net"

const (
	Version = "0.1.0"
)

var (
	AliveChannel     = make(chan bool)
	GHActionIps      = []string{}
	CLIConn net.Conn
)

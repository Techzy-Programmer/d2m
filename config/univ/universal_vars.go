package univ

import "net"

const (
	Version = "0.1.0"
)

var (
	IsProd           = false
	AliveChannel     = make(chan bool)
	GHActionIps      = []string{}
	CLIConn net.Conn
)

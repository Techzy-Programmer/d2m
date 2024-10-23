package vars

import (
	"crypto/rsa"
	"net"
)

var (
	IsDaemon     = false
	StartedAt    int64
	IsProd       = false
	AliveChannel = make(chan bool)
	LocalIPs     = []string{"127.0.0.1", "localhost", "::1", "::", ""}
	GHActionIps  = []string{}
	CLIConn      net.Conn
	PrivKey      *rsa.PrivateKey
)

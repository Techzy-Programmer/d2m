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
	GHActionIps  = []string{}
	CLIConn      net.Conn
	PrivKey      *rsa.PrivateKey
)

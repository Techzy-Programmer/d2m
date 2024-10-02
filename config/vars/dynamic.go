package vars

import (
	"crypto/rsa"
	"net"
)

var (
	IsProd       = false
	AliveChannel = make(chan bool)
	GHActionIps  = []string{}
	CLIConn      net.Conn
	PrivKey      *rsa.PrivateKey
)

package univ

import (
	"crypto/rsa"
	"net"
)

const (
	Version = "0.1.0"
)

var (
	IsProd       = false
	AliveChannel = make(chan bool)
	GHActionIps  = []string{}
	CLIConn      net.Conn
	PrivKey      *rsa.PrivateKey
)

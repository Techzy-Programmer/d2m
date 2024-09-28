package msg

type MSG interface { }

const (
	PingMsgType = "ping"
)

type PingMSG struct {
	Type string
	IsWelcome bool // Used by daemon on first connection
}

var typeRegistry = map[string] func() interface {} {
	PingMsgType: func() interface {} { return &PingMSG{} },
}

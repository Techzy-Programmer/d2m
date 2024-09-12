package msg

type MSG interface { }

const (
	PingMsgType = "ping"
)

type PingMSG struct {
	Type string
}

var typeRegistry = map[string] func() interface {} {
	PingMsgType: func() interface {} { return &PingMSG{} },
}

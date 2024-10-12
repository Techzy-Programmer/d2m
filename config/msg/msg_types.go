package msg

type MSG interface{}

const (
	PingMsgType         = "ping"
	HaltMsgType         = "halt"
	ConfigUpdateMsgType = "config-update"
)

type PingMSG struct {
	Type      string
	IsWelcome bool // Used by daemon on first connection
}

type HaltMSG struct {
	Type string
	Ack  bool // Used by daemon to acknowledge the halt
}

type ConfigUpdateMSG struct {
	Type  string
	Which string
}

var typeRegistry = map[string]func() interface{}{
	PingMsgType:         func() interface{} { return &PingMSG{} },
	HaltMsgType:         func() interface{} { return &HaltMSG{} },
	ConfigUpdateMsgType: func() interface{} { return &ConfigUpdateMSG{} },
}

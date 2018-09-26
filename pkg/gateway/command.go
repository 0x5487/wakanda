package gateway

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Command struct {
	OP   string          `json:"op"`
	Data json.RawMessage `json:"data"`
}

func CreateCommand(buf []byte) (*Command, error) {
	var command Command

	err := json.Unmarshal(buf, &command)
	if err != nil {
		return nil, err
	}
	return &command, nil
}

func (c *Command) ToWSMessage() (*WSMessage, error) {
	buf, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	msg := &WSMessage{websocket.TextMessage, buf}
	return msg, nil
}

package nats

import (
	"encoding/json"

	"github.com/jasonsoft/wakanda/pkg/messenger"
	"github.com/nats-io/go-nats-streaming"
)

type MessagePublisher struct {
	natsConn stan.Conn
}

func NewMessagePublisher(natsConn stan.Conn) *MessagePublisher {
	return &MessagePublisher{
		natsConn: natsConn,
	}
}

func (pub *MessagePublisher) PublishToDeliveryChannel(msg *messenger.Message) error {
	// to byte
	bytes, err := json.Marshal(msg)
	if err != nil {

	}

	// send messages to delivery subject
	return pub.natsConn.Publish("delivery", bytes)
}

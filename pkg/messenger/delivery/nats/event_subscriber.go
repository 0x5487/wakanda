package nats

import (
	"log"

	"github.com/nats-io/go-nats-streaming"
)

type EventSubscriber struct {
	natsConn stan.Conn
}

func NewEventSubscriber(natsConn stan.Conn) *EventSubscriber {
	return &EventSubscriber{
		natsConn: natsConn,
	}
}

func (evt *EventSubscriber) SubscribeDeliverySubject() {
	evt.natsConn.QueueSubscribe("delivery", "worker1", func(m *stan.Msg) {
		m.Ack()
		log.Printf("Received a message: %s\n ", string(m.Data))
	}, stan.SetManualAckMode(), stan.DurableName("delivery-remember"))
}

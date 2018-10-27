package nats

import (
	"context"
	"encoding/json"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/wakanda/pkg/messenger"
	"github.com/nats-io/go-nats-streaming"
)

type EventSubscriber struct {
	natsConn    stan.Conn
	deliverySvc messenger.DeliveryServicer
}

func NewEventSubscriber(natsConn stan.Conn) *EventSubscriber {
	return &EventSubscriber{
		natsConn: natsConn,
	}
}

func (evt *EventSubscriber) SubscribeDeliverySubject(ctx context.Context) {
	evt.natsConn.QueueSubscribe("delivery", "worker1", func(m *stan.Msg) {
		log.Debugf("delivery: received a message: %s", string(m.Data))

		messages := []*messenger.Message{}
		err := json.Unmarshal(m.Data, &messages)
		if err != nil {
			log.Errorf("delivery: receive message failed: %v", err)
			return
		}

		err = evt.deliverySvc.DeliveryMessage(ctx, messages)
		if err != nil {
			return
		}

		m.Ack()
	}, stan.SetManualAckMode(), stan.DurableName("delivery-remember"))
	log.Info("delivery: delivery subject was subscribed")
}

func (evt *EventSubscriber) Shutdown(ctx context.Context) error {
	return nil
}

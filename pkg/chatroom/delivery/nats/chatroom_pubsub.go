package nats

import (
	"context"
	"encoding/json"

	"github.com/gogo/protobuf/proto"

	"github.com/jasonsoft/log"
	chatroomProto "github.com/jasonsoft/wakanda/pkg/chatroom/proto"
	deliveryProto "github.com/jasonsoft/wakanda/pkg/delivery/proto"
	"github.com/jasonsoft/wakanda/pkg/gateway"
	"github.com/jasonsoft/wakanda/pkg/messenger"
	"github.com/nats-io/go-nats-streaming"
)

type ChatroomPubSub struct {
	natsConn stan.Conn
}

func NewChatroomPubSub(natsConn stan.Conn) *ChatroomPubSub {
	return &ChatroomPubSub{
		natsConn: natsConn,
	}
}

func (ps *ChatroomPubSub) SubscribeMSGChatroom(ctx context.Context) {
	ps.natsConn.QueueSubscribe("msg-chatroom", "worker1", func(m *stan.Msg) {
		defer func() {
			m.Ack()
		}()

		cmdChatroom := chatroomProto.ChatroomCommand{}
		err := proto.Unmarshal(m.Data, &cmdChatroom)
		if err != nil {
			log.Errorf("chatroom: proto unmarshal message failed: %v", err)
			return
		}

		log.Debugf("chatroom: msg from sender_id: %s", cmdChatroom.SenderID)

		msg := messenger.Message{}
		err = json.Unmarshal(cmdChatroom.Data, &msg)
		if err != nil {
			log.Errorf("delivery: json marshal command failed: %v", err)
			return
		}

		// filter message
		msg.Content = "filter~" + msg.Content

		// prepare message and send to delivery service
		bytes, err := json.Marshal(msg)

		command := &gateway.Command{
			OP:   "MSG",
			Data: bytes,
		}
		data, err := json.Marshal(command)
		if err != nil {
			log.Errorf("delivery: json marshal command failed: %v", err)
			return
		}

		// send messages to "delivery-chatroom" subject
		cmdDelivery := deliveryProto.DeliveryChatroomMessageCommand{
			RoomID:          cmdChatroom.RoomID,
			SenderID:        cmdChatroom.SenderID,
			SenderFirstName: cmdChatroom.SenderFirstName,
			SenderLastName:  cmdChatroom.SenderLastName,
			Data:            data,
		}

		bytes, err = proto.Marshal(&cmdDelivery)
		ps.natsConn.Publish("delivery-chatroom", bytes)

	}, stan.SetManualAckMode(), stan.DurableName("msg-chatroom-remember"))
	log.Info("chatroom: msg-chatroom subject was subscribed")
}

func (sub *ChatroomPubSub) Shutdown(ctx context.Context) error {
	return nil
}

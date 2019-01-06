package nats

import (
	"context"

	"github.com/jasonsoft/log"

	"github.com/golang/protobuf/proto"
	chatroomProto "github.com/jasonsoft/wakanda/pkg/chatroom/proto"
	"github.com/nats-io/go-nats-streaming"
)

type DispatcherPub struct {
	natsConn stan.Conn
}

func NewDispatcherPub(natsConn stan.Conn) *DispatcherPub {
	return &DispatcherPub{
		natsConn: natsConn,
	}
}

func (pub *DispatcherPub) PublishToChatMessageChannel(ctx context.Context, cmd *chatroomProto.ChatroomCommand) error {
	bytes, err := proto.Marshal(cmd)
	if err != nil {
		log.Errorf("dispatcher: proto marshal message failed: %v", err)
		return err
	}

	// send messages to delivery subject
	return pub.natsConn.Publish("msg-chatroom", bytes)
}

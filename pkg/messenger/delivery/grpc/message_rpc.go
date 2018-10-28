package grpc

import (
	"context"
	"encoding/json"

	"github.com/jasonsoft/wakanda/pkg/messenger"

	"google.golang.org/grpc/metadata"

	"github.com/jasonsoft/log"
	messengerNats "github.com/jasonsoft/wakanda/pkg/messenger/delivery/nats"
	"github.com/jasonsoft/wakanda/pkg/messenger/proto"
)

var (
	_emptyReply = &proto.EmptyReply{}
)

type MessageServer struct {
	messageSvc messenger.MessageServicer
	messagePub *messengerNats.MessagePublisher
}

func NewMessageServer(messageservice messenger.MessageServicer, messagePub *messengerNats.MessagePublisher) *MessageServer {
	return &MessageServer{
		messageSvc: messageservice,
		messagePub: messagePub,
	}
}

func (s *MessageServer) CreateMessage(ctx context.Context, in *proto.CreateMessageRequest) (*proto.EmptyReply, error) {
	reqID := ""
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if len(md["req_id"]) > 0 {
			reqID = md["req_id"][0]
		}
	}
	customFields := log.Fields{
		"req_id": reqID,
	}
	logger := log.WithFields(customFields)

	msg := &messenger.Message{}
	err := json.Unmarshal(in.Data, msg)
	if err != nil {
		return nil, err
	}
	logger.Debugf("messenger: MSG data: %s", msg.Content)
	msg.RequestID = reqID

	err = s.messageSvc.CreateMessage(ctx, msg)
	if err != nil {
		return nil, err
	}

	err = s.messagePub.PublishToDeliveryChannel(msg)
	if err != nil {
		return nil, err
	}

	return _emptyReply, nil
}

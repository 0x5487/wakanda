package grpc

import (
	"context"
	"encoding/json"

	"github.com/jasonsoft/wakanda/pkg/messenger"

	"google.golang.org/grpc/metadata"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/wakanda/pkg/messenger/proto"
)

type MessageServer struct {
	messageSvc messenger.MessageServicer
}

func NewMessageServer(messageservice messenger.MessageServicer) *MessageServer {
	return &MessageServer{
		messageSvc: messageservice,
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

	content := ""
	err := json.Unmarshal(in.Data, &content)
	if err != nil {
		return nil, err
	}
	logger.Debugf("messenger: MSG data: %s", content)

	msg := &messenger.Message{
		RequestID: reqID,
	}
	s.messageSvc.CreateMessage(ctx, msg)

	return new(proto.EmptyReply), nil
}

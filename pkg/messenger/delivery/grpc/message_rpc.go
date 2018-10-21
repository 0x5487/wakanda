package grpc

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc/metadata"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/wakanda/pkg/messenger/proto"
)

type MessageServer struct {
}

func NewMessageServer() *MessageServer {
	return &MessageServer{}
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

	msg := ""
	err := json.Unmarshal(in.Data, &msg)
	if err != nil {
		return nil, err
	}
	logger.Debugf("messenger: MSG data: %s", msg)

	return new(proto.EmptyReply), nil
}

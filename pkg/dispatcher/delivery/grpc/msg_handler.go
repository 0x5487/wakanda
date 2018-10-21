package grpc

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/wakanda/pkg/dispatcher/proto"
	messageProto "github.com/jasonsoft/wakanda/pkg/messenger/proto"
)

func (svc *DispatcherServer) handleMSG(ctx context.Context, in *proto.CommandRequest) (*proto.CommandReply, error) {
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

	logger.Info("dispatcher: msg is received")

	// pass message to messenger service
	createMsgCmd := &messageProto.CreateMessageRequest{
		Data: in.Data,
	}
	ctx = metadata.NewOutgoingContext(ctx, md)
	_, err := svc.messageRPCClient.CreateMessage(ctx, createMsgCmd)
	if err != nil {
		logger.Errorf("dispatcher: err from message RPC: %v", err)
		return nil, err
	}

	reply := &proto.CommandReply{
		ReqID: reqID,
		OP:    "ACK",
	}

	return reply, nil
}

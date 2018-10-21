package grpc

import (
	"context"
	"errors"

	"github.com/jasonsoft/log"

	"github.com/jasonsoft/wakanda/pkg/dispatcher/proto"
	messageProto "github.com/jasonsoft/wakanda/pkg/messenger/proto"
)

type DispatcherServer struct {
	messageRPCClient messageProto.MessageServiceClient
}

func NewDispatchServer(messageRPCClient messageProto.MessageServiceClient) *DispatcherServer {
	return &DispatcherServer{
		messageRPCClient: messageRPCClient,
	}
}

func (svc *DispatcherServer) HandleCommand(ctx context.Context, in *proto.CommandRequest) (*proto.CommandReply, error) {
	log.Debugf("dispatcher: receiver op: %s", in.OP)

	switch in.OP {
	case "MSG":
		return svc.handleMSG(ctx, in)
	default:
		log.Warnf("dispatcher: unknown command: %s", in.OP)
		reply := &proto.CommandReply{}
		return reply, errors.New("unknown command")
	}

}

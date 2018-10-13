package grpc

import (
	"context"
	"errors"

	"github.com/jasonsoft/log"

	"github.com/jasonsoft/wakanda/pkg/messenger/proto"
)

type DispatcherServer struct {
}

func NewDispatchServer() *DispatcherServer {
	return &DispatcherServer{}
}

func (s *DispatcherServer) HandleCommand(ctx context.Context, in *proto.HandleCommandRequest) (*proto.HandleCommandReply, error) {
	log.Debugf("grpc: receiver op: %s", in.OP)

	switch in.OP {
	case "MSG":
		return handleMSG(ctx, in)
	default:
		log.Warnf("gateway: unknown command: %s", in.OP)
		reply := &proto.HandleCommandReply{}
		return reply, errors.New("unknown command")
	}

}

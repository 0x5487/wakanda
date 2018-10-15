package grpc

import (
	"context"
	"errors"

	"github.com/jasonsoft/log"

	"github.com/jasonsoft/wakanda/pkg/dispatcher/proto"
)

type DispatcherServer struct {
}

func NewDispatchServer() *DispatcherServer {
	return &DispatcherServer{}
}

func (s *DispatcherServer) HandleCommand(ctx context.Context, in *proto.CommandRequest) (*proto.CommandReply, error) {
	log.Debugf("dispatcher: receiver op: %s", in.OP)

	switch in.OP {
	case "MSG":
		return handleMSG(ctx, in)
	default:
		log.Warnf("dispatcher: unknown command: %s", in.OP)
		reply := &proto.CommandReply{}
		return reply, errors.New("unknown command")
	}

}

package grpc

import (
	"context"
	"errors"

	"github.com/jasonsoft/log"

	deliveryNats "github.com/jasonsoft/wakanda/pkg/dispatcher/delivery/nats"
	deliveryProto "github.com/jasonsoft/wakanda/pkg/dispatcher/proto"
)

var (
	_emptyDispatcherCommandReply = deliveryProto.DispatcherCommandReply{}
)

type DispatcherServer struct {
	dispatcherPub *deliveryNats.DispatcherPub
}

func NewDispatcherServer(nats *deliveryNats.DispatcherPub) *DispatcherServer {
	return &DispatcherServer{
		dispatcherPub: nats,
	}
}

func (svc *DispatcherServer) HandleCommand(ctx context.Context, in *deliveryProto.DispatcherCommandRequest) (*deliveryProto.DispatcherCommandReply, error) {
	log.Debugf("dispatcher: receiver op: %s", in.OP)

	switch in.OP {
	case "MSGRM":
		return svc.handleMSGRM(ctx, in)
	default:
		log.Warnf("dispatcher: unknown command: %s", in.OP)
		reply := &deliveryProto.DispatcherCommandReply{}
		return reply, errors.New("unknown command")
	}

}

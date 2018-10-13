package grpc

import (
	"context"
	"encoding/json"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/wakanda/pkg/messenger/proto"
)

func handleMSG(ctx context.Context, in *proto.HandleCommandRequest) (*proto.HandleCommandReply, error) {
	msg := ""
	err := json.Unmarshal(in.Data, &msg)
	if err != nil {
		return nil, err
	}
	log.Debugf("messenger: MSG data: %s", msg)
	// TODO: filter content

	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	reply := &proto.HandleCommandReply{
		OP:   "MSG",
		Data: data,
	}

	return reply, nil
}

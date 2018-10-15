package grpc

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc/metadata"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/wakanda/pkg/dispatcher/proto"
)

func handleMSG(ctx context.Context, in *proto.CommandRequest) (*proto.CommandReply, error) {
	msg := ""
	err := json.Unmarshal(in.Data, &msg)
	if err != nil {
		return nil, err
	}
	log.Debugf("dispatcher: req_id: %s, MSG data: %s", in.ReqID, msg)
	// TODO: filter content

	reqID := ""
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		reqID = md["req_id"][0]
	}

	reply := &proto.CommandReply{
		ReqID: reqID,
		OP:    "ACK",
	}

	return reply, nil
}

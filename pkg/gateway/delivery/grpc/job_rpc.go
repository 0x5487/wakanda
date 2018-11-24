package grpc

import (
	"context"
	"encoding/json"

	"github.com/jasonsoft/log"

	"github.com/jasonsoft/wakanda/pkg/gateway"

	gatewayProto "github.com/jasonsoft/wakanda/pkg/gateway/proto"
)

var (
	_emptyReply = &gatewayProto.EmptyReply{}
)

type JobServer struct {
	manager *gateway.Manager
}

func NewJobServer(manager *gateway.Manager) *JobServer {
	return &JobServer{
		manager: manager,
	}
}

func (s *JobServer) SendJobs(ctx context.Context, in *gatewayProto.SendJobRequest) (*gatewayProto.EmptyReply, error) {
	log.Debugf("gateway: jobs count: %d", len(in.Jobs))

	var err error
	for _, rpcJob := range in.Jobs {
		command := &gateway.Command{}
		err = json.Unmarshal(rpcJob.Data, command)
		if err != nil {
			log.Errorf("gateway: json unmarshal command failed: %v ", err)
			continue
		}

		switch rpcJob.Type {
		case "S":
			s.SendCommandToSession(context.Background(), rpcJob.TargetID, command)
		}
	}

	return _emptyReply, nil
}

func (s *JobServer) SendCommandToSession(ctx context.Context, sessionID string, command *gateway.Command) {
	s.manager.Push(sessionID, command)
}

package grpc

import (
	"context"
	"encoding/json"

	"github.com/jasonsoft/wakanda/pkg/gateway"

	gatewayProto "github.com/jasonsoft/wakanda/pkg/gateway/proto"
)

type GatewayServer struct {
	manager *gateway.Manager
}

func NewGatewayServer(manager *gateway.Manager) *GatewayServer {
	return &GatewayServer{
		manager: manager,
	}
}

func (s *GatewayServer) SendJobs(ctx context.Context, in *gatewayProto.SendJobRequest) (*gatewayProto.EmptyReply, error) {

	for _, rpcJob := range in.Jobs {
		command := &gateway.Command{}
		err := json.Unmarshal(rpcJob.Data, command)
		if err != nil {
			continue
		}

		switch rpcJob.Type {
		case "S":
			s.SendCommandToSession(context.Background(), rpcJob.TargetID, command)
		}
	}

	return nil, nil
}

func (s *GatewayServer) SendCommandToSession(ctx context.Context, sessionID string, command *gateway.Command) {
	s.manager.Push(sessionID, command)
}

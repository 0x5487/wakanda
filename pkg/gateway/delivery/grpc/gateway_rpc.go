package grpc

import (
	"context"

	gatewayProto "github.com/jasonsoft/wakanda/pkg/gateway/proto"
)

type GatewayServer struct {
	gatewayClient gatewayProto.GatewayServiceClient
}

func NewGatewayServer(gatewayClient gatewayProto.GatewayServiceClient) *GatewayServer {
	return &GatewayServer{
		gatewayClient: gatewayClient,
	}
}

func (s *GatewayServer) SendCommand(ctx context.Context, in *gatewayProto.SendCommandRequest) (*gatewayProto.EmptyReply, error) {

	for _, command := range in.Commands {
		command.
	}

	return nil, nil
}

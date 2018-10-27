package grpc

import (
	"context"

	"github.com/jasonsoft/wakanda/pkg/router/proto"
	"google.golang.org/grpc"
)

type RouterServer struct {
}

func (s *RouterServer) Routes(ctx context.Context, in *proto.RouteRequest, opts ...grpc.CallOption) (*proto.RouteReply, error) {
	panic("not implemented")
}

func (s *RouterServer) CreateOrUpdateRoute(ctx context.Context, in *proto.CreateOrUpdateRouteRequest, opts ...grpc.CallOption) (*proto.EmptyReply, error) {
	panic("not implemented")
}

func (s *RouterServer) DeleteSession(ctx context.Context, in *proto.DeleteSessionRequest, opts ...grpc.CallOption) (*proto.EmptyReply, error) {
	panic("not implemented")
}

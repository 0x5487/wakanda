package grpc

import (
	"context"

	"github.com/go-redis/redis"

	"github.com/golang/protobuf/ptypes"
	"github.com/jasonsoft/wakanda/pkg/router/proto"
)

var (
	_emptyReply = &proto.EmptyReply{}
)

type RouterServer struct {
	redisClient *redis.Client
}

func NewRouterServer(redisClient *redis.Client) *RouterServer {
	return &RouterServer{
		redisClient: redisClient,
	}
}

func (s *RouterServer) Routes(ctx context.Context, in *proto.RouteRequest) (*proto.RouteReply, error) {

	result := proto.RouteReply{}

	for _, memberID := range in.MemberIDs {
		// get member key
		memberKey := "m:" + memberID
		sessionIDs, err := s.redisClient.SMembers(memberKey).Result()
		if err != nil {
			continue
		}

		for _, sessionID := range sessionIDs {
			sessionKey := "s:" + sessionID
			routeinfo, err := s.redisClient.HGetAll(sessionKey).Result()
			if err != nil {
				continue
			}

			route := &proto.Route{
				SessionID:   sessionID,
				MemberID:    memberID,
				GatewayAddr: routeinfo["gateway_addr"],
			}
			result.Routes = append(result.Routes, route)
		}
	}

	return &result, nil

}

func (s *RouterServer) CreateOrUpdateRoute(ctx context.Context, in *proto.CreateOrUpdateRouteRequest) (*proto.EmptyReply, error) {
	// create session key
	m := map[string]interface{}{
		"member_id":    in.GetMemberID(),
		"gateway_addr": in.GatewayAddr,
		"last_seen":    ptypes.TimestampNow(),
	}

	sessionKey := "s:" + in.SessionID
	_, err := s.redisClient.HMSet(sessionKey, m).Result()
	if err != nil {
		return nil, err
	}

	// create member key
	memberKey := "m:" + in.MemberID
	_, err = s.redisClient.SAdd(memberKey, in.SessionID).Result()
	if err != nil {
		return nil, err
	}

	return _emptyReply, nil
}

func (s *RouterServer) DeleteSession(ctx context.Context, in *proto.DeleteSessionRequest) (*proto.EmptyReply, error) {
	panic("not implemented")
}

func (s *RouterServer) Shutdown(ctx context.Context) error {
	return nil
}

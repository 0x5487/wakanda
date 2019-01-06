package grpc

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/jasonsoft/log"
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
	log.Debug("router: === Begin Routes ===")

	result := proto.RouteReply{}

	validLastSeenTimestamp := int(time.Now().Add(-5 * time.Minute).Unix())

	for _, memberID := range in.MemberIDs {
		// get member key
		memberKey := "m:" + memberID
		log.Debugf("router: memberKey: %s", memberKey)

		sessionIDs, err := s.redisClient.SMembers(memberKey).Result()
		if err != nil {
			log.Errorf("router: redis SMembers command failed: %v", err)
			continue
		}

		if len(sessionIDs) == 0 {
			log.Debugf("router: no session id was found by memberKey: %s", memberKey)
			continue
		}

		for _, sessionID := range sessionIDs {
			sessionKey := "s:" + sessionID
			routeinfo, err := s.redisClient.HGetAll(sessionKey).Result()
			if err != nil {
				log.Errorf("router: redis HGetAll command failed: %v", err)
				continue
			}

			// if the session's last_seen is 5 mins ago (invalid session), delete the session.
			lastseen, err := strconv.Atoi(routeinfo["last_seen"])
			if err != nil {
				log.Errorf("router: last seen cast failed: %v", err)
				continue
			}

			if lastseen < validLastSeenTimestamp {
				// delete session
				_, err = s.redisClient.HDel(sessionKey).Result()
				if err != nil {
					log.Errorf("router: redis HGetAll command failed: %v", err)
					continue
				}

				// delete memberKey
				_, err = s.redisClient.SRem(memberKey, sessionID).Result()
				if err != nil {
					log.Errorf("router: redis SRem command failed: %v", err)
					continue
				}
			}

			route := &proto.Route{
				SessionID:   sessionID,
				MemberID:    memberID,
				GatewayAddr: routeinfo["gateway_addr"],
			}
			result.Routes = append(result.Routes, route)
		}
	}

	log.Infof("router: number of routes: %d", len(result.Routes))

	return &result, nil

}

func (s *RouterServer) CreateOrUpdateRoute(ctx context.Context, in *proto.CreateOrUpdateRouteRequest) (*proto.EmptyReply, error) {
	log.Debug("router: === Begin CreateOrUpdateRoute ===")

	// create session key
	m := map[string]interface{}{
		"member_id":    in.GetMemberID(),
		"gateway_addr": in.GatewayAddr,
		"last_seen":    int32(time.Now().Unix()),
	}

	sessionKey := "s:" + in.SessionID
	_, err := s.redisClient.HMSet(sessionKey, m).Result()
	if err != nil {
		log.Errorf("router: redis HMSET failed: %v", err)
		return nil, err
	}

	// create member key
	memberKey := "m:" + in.MemberID
	_, err = s.redisClient.SAdd(memberKey, in.SessionID).Result()
	if err != nil {
		log.Errorf("router: redis SADD failed: %v", err)
		return nil, err
	}

	//log.Debug("router: === End CreateOrUpdateRoute ===")
	return _emptyReply, nil
}

func (s *RouterServer) DeleteSession(ctx context.Context, in *proto.DeleteSessionRequest) (*proto.EmptyReply, error) {
	panic("not implemented")
}

func (s *RouterServer) Shutdown(ctx context.Context) error {
	return nil
}

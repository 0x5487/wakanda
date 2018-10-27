package redis

import (
	"github.com/go-redis/redis"
	"github.com/jasonsoft/wakanda/pkg/router/proto"
)

type RouteRepo struct {
	redisClient *redis.Client
}

func (repo *RouteRepo) CreateOrUpdateRoute(route *proto.CreateOrUpdateRouteRequest) error {
	
}

package router

import "github.com/jasonsoft/wakanda/pkg/router/proto"

type RouterRepository interface {
	CreateOrUpdateRoute(route *proto.CreateOrUpdateRouteRequest) error
}

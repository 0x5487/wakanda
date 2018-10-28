package main

import (
	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/console"
	"github.com/jasonsoft/log/handlers/gelf"
	"github.com/jasonsoft/napnap"
	"github.com/jasonsoft/wakanda/internal/config"
	"github.com/jasonsoft/wakanda/internal/identity"
	"github.com/jasonsoft/wakanda/internal/middleware"
	"github.com/jasonsoft/wakanda/pkg/dispatcher/proto"
	"github.com/jasonsoft/wakanda/pkg/gateway"
	gatewayHttp "github.com/jasonsoft/wakanda/pkg/gateway/delivery/http"
	routerProto "github.com/jasonsoft/wakanda/pkg/router/proto"
	"google.golang.org/grpc"
)

var (
	_manager          *gateway.Manager
	_dispatcherClient proto.DispatcherClient
	_routerClient     routerProto.RouterServiceClient
)

func initialize(config *config.Configuration) {
	initLogger("gateway", config)

	_manager = gateway.NewManager()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(config.Dispatcher.AdvertiseAddr, opts...)
	if err != nil {
		log.Fatalf("gateway: can't connect to messenger grpc service: %v", err)
	}
	log.Info("gateway: dispatcher service was connected")
	_dispatcherClient = proto.NewDispatcherClient(conn)

	routerConn, err := grpc.Dial(config.Router.AdvertiseAddr, opts...)
	if err != nil {
		log.Fatalf("gateway: can't connect to router grpc service: %v", err)
	}
	log.Info("gateway: router service was connected")
	_routerClient = routerProto.NewRouterServiceClient(routerConn)

}

func initLogger(appID string, config *config.Configuration) {
	// set up log target
	log.SetAppID(appID)
	for _, target := range config.Logs {
		switch target.Type {
		case "console":
			clog := console.New()
			levels := log.GetLevelsFromMinLevel(target.MinLevel)
			log.RegisterHandler(clog, levels...)
		case "gelf":
			graylog := gelf.New(target.ConnectionString)
			levels := log.GetLevelsFromMinLevel(target.MinLevel)
			log.RegisterHandler(graylog, levels...)
		}
	}
}

func napWithMiddlewares() *napnap.NapNap {
	nap := napnap.New()
	corsOpts := napnap.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"*"},
	}
	nap.Use(napnap.NewCors(corsOpts))
	nap.Use(napnap.NewHealth())
	nap.Use(middleware.NewErrorHandingMiddleware())
	nap.Use(identity.NewMiddleware())
	httpHandler := gatewayHttp.NewGatewayHttpHandler(_manager, _dispatcherClient, _routerClient)
	nap.Use(gatewayHttp.NewGatewayRouter(httpHandler))
	return nap
}

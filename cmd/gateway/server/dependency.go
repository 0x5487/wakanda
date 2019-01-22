package main

import (
	"os"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/console"
	"github.com/jasonsoft/log/handlers/gelf"
	"github.com/jasonsoft/napnap"
	"github.com/jasonsoft/wakanda/internal/config"
	"github.com/jasonsoft/wakanda/internal/middleware"
	dispatcherProto "github.com/jasonsoft/wakanda/pkg/dispatcher/proto"
	"github.com/jasonsoft/wakanda/pkg/gateway"
	gatewayHttp "github.com/jasonsoft/wakanda/pkg/gateway/delivery/http"
	identityHttp "github.com/jasonsoft/wakanda/pkg/identity/delivery/http"
	routerProto "github.com/jasonsoft/wakanda/pkg/router/proto"
	"google.golang.org/grpc"
)

var (
	_manager          *gateway.Manager
	_dispatcherClient dispatcherProto.DispatcherServiceClient
	_routerClient     routerProto.RouterServiceClient
)

func initialize(config *config.Configuration) {
	initLogger("gateway", config)

	// setup manager

	_manager = gateway.NewManager()
	gatewayAddr := os.Getenv("gateway_addr")
	if len(gatewayAddr) == 0 {
		gatewayAddr, _ = os.Hostname()
		gatewayAddr += ":19997" // TODO: can't hard code here
	}
	_manager.SetGatewayAddr(gatewayAddr)

	// setup displatcher client
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(config.Dispatcher.AdvertiseAddr, opts...)
	if err != nil {
		log.Fatalf("gateway: can't connect to dispatcher grpc service: %v", err)
	}
	log.Info("gateway: dispatcher service was connected")
	_dispatcherClient = dispatcherProto.NewDispatcherServiceClient(conn)

	// setup router client
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

func napWithMiddlewares(config *config.Configuration) *napnap.NapNap {
	nap := napnap.New()
	corsOpts := napnap.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"*"},
	}
	nap.Use(napnap.NewCors(corsOpts))
	nap.Use(napnap.NewHealth())
	nap.Use(middleware.NewErrorHandingMiddleware())
	nap.Use(identityHttp.NewAuthMiddleware(config))
	nap.Use(gatewayHttp.NewPrometheusMiddleware(_manager))
	httpHandler := gatewayHttp.NewGatewayHttpHandler(_manager, _dispatcherClient, _routerClient)
	nap.Use(gatewayHttp.NewGatewayRouter(httpHandler))
	nap.Use(gatewayHttp.NewPrometheusRouter())
	return nap
}

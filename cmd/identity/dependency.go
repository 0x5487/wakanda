package main

import (
	"os"

	"github.com/jasonsoft/napnap"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/console"
	"github.com/jasonsoft/log/handlers/gelf"
	"github.com/jasonsoft/wakanda/internal/config"
	"github.com/jasonsoft/wakanda/internal/middleware"
	dispatcherGRPC "github.com/jasonsoft/wakanda/pkg/dispatcher/delivery/grpc"
	dispatcherNats "github.com/jasonsoft/wakanda/pkg/dispatcher/delivery/nats"
	identityHttp "github.com/jasonsoft/wakanda/pkg/identity/delivery/http"
	"github.com/nats-io/go-nats-streaming"
)

var (
	// grpc servers
	_dispatcherServer *dispatcherGRPC.DispatcherServer
)

func initialize(config *config.Configuration) error {
	initLogger("dispatcher", config)

	natsConn, err := setupNatsConn(config)
	if err != nil {
		return err
	}

	dispatcherPub := dispatcherNats.NewDispatcherPub(natsConn)
	_dispatcherServer = dispatcherGRPC.NewDispatcherServer(dispatcherPub)

	return nil
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

func setupNatsConn(config *config.Configuration) (stan.Conn, error) {
	hostname, _ := os.Hostname()
	clientID := "dispatcher-" + hostname
	natsConn, err := stan.Connect(config.Nats.ClusterID, clientID, stan.NatsURL("nats://"+config.Nats.Address))
	if err != nil {
		return nil, err
	}
	return natsConn, nil
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
	return nap
}

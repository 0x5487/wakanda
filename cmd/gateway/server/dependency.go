package main

import (
	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/console"
	"github.com/jasonsoft/log/handlers/gelf"
	"github.com/jasonsoft/napnap"
	"github.com/jasonsoft/wakanda/internal/config"
	"github.com/jasonsoft/wakanda/internal/middleware"
	gatewayHttp "github.com/jasonsoft/wakanda/pkg/gateway/delivery/http"
)

func initialize(config *config.Configuration) {
	initLogger("gateway", config)

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
	nap.Use(gatewayHttp.NewGatewayRouter())
	return nap
}

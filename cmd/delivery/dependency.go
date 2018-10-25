package main

import (
	"os"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/console"
	"github.com/jasonsoft/log/handlers/gelf"
	"github.com/jasonsoft/wakanda/internal/config"
	"github.com/nats-io/go-nats-streaming"
)

var (
	_natsConn stan.Conn
)

func initialize(config *config.Configuration) error {
	initLogger("delivery", config)

	hostname, _ := os.Hostname()

	var err error
	_natsConn, err = stan.Connect(config.Nats.ClusterID, hostname, stan.NatsURL("nats://"+config.Nats.Address))
	if err != nil {
		return err
	}

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

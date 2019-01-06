package main

import (
	"os"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/console"
	"github.com/jasonsoft/log/handlers/gelf"
	"github.com/jasonsoft/wakanda/internal/config"
	chatroomNats "github.com/jasonsoft/wakanda/pkg/chatroom/delivery/nats"
	"github.com/nats-io/go-nats-streaming"
)

var (
	_chatroomPubSub *chatroomNats.ChatroomPubSub
)

func initialize(config *config.Configuration) error {
	initLogger("chatroom", config)

	natsConn, err := setupNatsConn(config)
	if err != nil {
		return err
	}

	_chatroomPubSub = chatroomNats.NewChatroomPubSub(natsConn)

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
	clientID := "chatroom-" + hostname
	natsConn, err := stan.Connect(config.Nats.ClusterID, clientID, stan.NatsURL("nats://"+config.Nats.Address))
	if err != nil {
		return nil, err
	}
	return natsConn, nil
}

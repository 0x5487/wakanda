package main

import (
	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/console"
	"github.com/jasonsoft/log/handlers/gelf"
	"github.com/jasonsoft/wakanda/internal/config"
	dispatcherGRPC "github.com/jasonsoft/wakanda/pkg/dispatcher/delivery/grpc"
	messengerProto "github.com/jasonsoft/wakanda/pkg/messenger/proto"
	"google.golang.org/grpc"
)

var (
	// grpc servers
	_dispatcherServer *dispatcherGRPC.DispatcherServer

	// grpc clients
	_messageClient messengerProto.MessageServiceClient
)

func initialize(config *config.Configuration) {
	initLogger("dispatcher", config)

	// Set up a connection to the server.
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(config.Messenger.AdvertiseAddr, opts...)
	if err != nil {
		log.Fatalf("dispatcher: can't connect to messenger grpc service: %v", err)
	}
	//defer conn.Close()

	_messageClient = messengerProto.NewMessageServiceClient(conn)

	_dispatcherServer = dispatcherGRPC.NewDispatchServer(_messageClient)

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

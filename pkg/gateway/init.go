package gateway

import (
	"github.com/jasonsoft/wakanda/internal/config"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/wakanda/pkg/dispatcher/proto"
	"google.golang.org/grpc"
)

var (
	_manager          *Manager
	_dispatcherClient proto.DispatcherClient
)

func Initialize(config *config.Configuration) {
	_manager = NewManager()

	// Set up a connection to the server.
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(config.Dispatcher.AdvertiseAddr, opts...)
	if err != nil {
		log.Fatalf("gateway: can't connect to messenger grpc service: %v", err)
	}
	//defer conn.Close()

	_dispatcherClient = proto.NewDispatcherClient(conn)
}

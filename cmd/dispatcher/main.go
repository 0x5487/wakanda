package main

import (
	"fmt"
	"net"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/wakanda/internal/config"
	dispatcherProto "github.com/jasonsoft/wakanda/pkg/dispatcher/proto"
	"google.golang.org/grpc"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			// unknown error
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("unknown error: %v", err)
			}
			log.StackTrace().Error(err)
		}
	}()

	config := config.New("app.yml")
	initialize(config)

	// start grpc servers
	lis, err := net.Listen("tcp", config.Dispatcher.GRPCBind)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	dispatcherProto.RegisterDispatcherServiceServer(s, _dispatcherServer)
	log.Info("dispatcher: grpc service started")
	if err = s.Serve(lis); err != nil {
		log.Fatalf("dispatcher: failed to start dispatcher grpc server: %v", err)
	}

}

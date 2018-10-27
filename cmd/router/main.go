package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/wakanda/internal/config"
	routerGRPC "github.com/jasonsoft/wakanda/pkg/router/delivery/grpc"
	routerProto "github.com/jasonsoft/wakanda/pkg/router/proto"
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

	// start grpc server
	lis, err := net.Listen("tcp", config.Router.GRPCBind)
	if err != nil {
		log.Fatalf("router: failed to listen: %v", err)
	}
	s := grpc.NewServer()
	routerServer := routerGRPC.NewRouterServer(_redisClient)
	routerProto.RegisterRouterServiceServer(s, routerServer)
	go func() {
		log.Info("router: grpc service started")
		if err = s.Serve(lis); err != nil {
			log.Fatalf("router: failed to start router grpc server: %v", err)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
	<-stopChan
	log.Info("router: shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := routerServer.Shutdown(ctx); err != nil {
		log.Errorf("router: service hanlder shutdown error: %v", err)
	} else {
		log.Info("router: gracefully stopped")
	}
}

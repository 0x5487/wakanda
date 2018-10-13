package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	messengerGRPC "github.com/jasonsoft/wakanda/pkg/messenger/delivery/grpc"
	messengerProto "github.com/jasonsoft/wakanda/pkg/messenger/proto"
	"google.golang.org/grpc"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/napnap"
	"github.com/jasonsoft/wakanda/internal/config"
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

	// start grpc
	lis, err := net.Listen("tcp", ":16998")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	dispatcherserver := messengerGRPC.NewDispatchServer()
	messengerProto.RegisterDispatcherServer(s, dispatcherserver)
	go func() {
		log.Info("messenger grpc service started")
		if err = s.Serve(lis); err != nil {
			log.Fatalf("failed to start dispatcher grpc server: %v", err)
		}
	}()

	// start http service
	nap := napWithMiddlewares()
	httpEngine := napnap.NewHttpEngine(config.Messenger.Bind)
	go func() {
		log.Info("messenger http service started")
		err := nap.Run(httpEngine)
		if err != nil {
			log.Error(err)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
	<-stopChan
	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpEngine.Shutdown(ctx); err != nil {
		log.Errorf("hanlder shutdown error: %v", err)
	} else {
		log.Info("gracefully stopped")
	}
}

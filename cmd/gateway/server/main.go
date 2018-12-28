package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/napnap"
	"github.com/jasonsoft/wakanda/internal/config"
	gatewayGRPC "github.com/jasonsoft/wakanda/pkg/gateway/delivery/grpc"
	gatewayProto "github.com/jasonsoft/wakanda/pkg/gateway/proto"
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
	lis, err := net.Listen("tcp", config.Gateway.GRPCBind)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	gatewayGrpcServer := grpc.NewServer()
	gatewayServer := gatewayGRPC.NewGatewayServer(_manager)
	gatewayProto.RegisterGatewayServiceServer(gatewayGrpcServer, gatewayServer)
	go func() {
		log.Info("gateway: grpc service started")
		if err = gatewayGrpcServer.Serve(lis); err != nil {
			log.Fatalf("gateway: failed to start gateway grpc server: %v", err)
		}
	}()

	// start job grpc server
	jobTcpListen, err := net.Listen("tcp", config.Gateway.JobGRPCBind)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	jobGrpcServer := grpc.NewServer()
	jobServer := gatewayGRPC.NewJobServer(_manager)
	gatewayProto.RegisterJobServiceServer(jobGrpcServer, jobServer)
	go func() {
		log.Info("gateway: job grpc service started")
		if err = jobGrpcServer.Serve(jobTcpListen); err != nil {
			log.Fatalf("gateway: failed to start job grpc server: %v", err)
		}
	}()

	// start http service
	nap := napWithMiddlewares()
	httpEngine := napnap.NewHttpEngine(config.Gateway.HTTPBind)
	go func() {
		log.Info("gateway: gateway service started")
		err := nap.Run(httpEngine)
		if err != nil {
			log.Error(err)
		}
	}()

	go func() {
		http.ListenAndServe(":18000", nil)
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
	<-stopChan
	log.Info("gateway: shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpEngine.Shutdown(ctx); err != nil {
		log.Errorf("gateway: hanlder shutdown error: %v", err)
	} else {
		log.Info("gateway: gracefully stopped")
	}
}

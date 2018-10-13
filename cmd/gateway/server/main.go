package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/console"
	"github.com/jasonsoft/napnap"
	"github.com/jasonsoft/wakanda/pkg/gateway"
	gatewayHttp "github.com/jasonsoft/wakanda/pkg/gateway/delivery/http"
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

	log.SetAppID("gateway") // unique id for the app

	clog := console.New()
	log.RegisterHandler(clog, log.AllLevels...)

	gateway.Initialize()

	nap := napnap.New()
	//router := napnap.NewRouter()
	nap.Use(napnap.NewHealth())
	nap.Use(gatewayHttp.NewGatewayRouter())

	httpEngine := napnap.NewHttpEngine(":19999")
	go func() {
		log.Info("gateway started")
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

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
	"github.com/jasonsoft/wakanda/internal/middleware"
	messengerHttp "github.com/jasonsoft/wakanda/pkg/messenger/delivery/http"
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

	log.SetAppID("messenger") // unique id for the app

	clog := console.New()
	log.RegisterHandler(clog, log.AllLevels...)

	nap := napnap.New()
	nap.ForwardRemoteIpAddress = true
	nap.Use(napnap.NewHealth())
	nap.Use(middleware.NewErrorHandingMiddleware())
	corsOpts := napnap.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"*"},
	}
	nap.Use(napnap.NewCors(corsOpts))
	nap.Use(middleware.NewIdentityMiddleware())
	nap.Use(messengerHttp.NewRouter())

	httpEngine := napnap.NewHttpEngine(":16999")
	go func() {
		log.Info("messenger started")
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

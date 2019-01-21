package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// start http service
	nap := napWithMiddlewares(config)
	httpEngine := napnap.NewHttpEngine(config.Identity.HTTPBind)
	go func() {
		log.Info("identity: identity service started")
		err := nap.Run(httpEngine)
		if err != nil {
			log.Error(err)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
	<-stopChan
	log.Info("identity: shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpEngine.Shutdown(ctx); err != nil {
		log.Errorf("identity: hanlder shutdown error: %v", err)
	} else {
		log.Info("identity: gracefully stopped")
	}

}

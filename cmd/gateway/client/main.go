package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.paradise-soft.com.tw/rd/log"
	"gitlab.paradise-soft.com.tw/rd/log/handlers/console"
	"gitlab.paradise-soft.com.tw/rd/napnap"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			// unknown error
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("unknown error: %v", err)
			}
		}
	}()

	log.SetAppID("client") // unique id for the app

	clog := console.New()
	log.RegisterHandler(clog, log.AllLevels...)

	nap := napWithMiddlewares()
	httpEngine := napnap.NewHttpEngine(":8080")
	go func() {
		log.Info("client started")
		err := nap.Run(httpEngine)
		if err != nil {
			log.Error(err)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
	<-stopChan
	log.Info("shutting down chat server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpEngine.Shutdown(ctx); err != nil {
		log.Errorf("chat server shutdown error: %v", err)
	} else {
		log.Info("gracefully stopped")
	}
}

func napWithMiddlewares() *napnap.NapNap {
	nap := napnap.New()
	nap.Use(napnap.NewHealth())
	corsOpts := napnap.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"*"},
	}
	nap.Use(napnap.NewCors(corsOpts))
	nap.SetRender("./templates")

	router := napnap.NewRouter()
	// display client page
	router.Get("/", func(c *napnap.Context) {
		c.Render(200, "chat.html", nil)
	})
	nap.Use(router)
	nap.Use(napnap.NewNotfoundMiddleware())

	return nap
}

package main

import (
	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/console"
	"github.com/jasonsoft/napnap"
	"github.com/jasonsoft/wakanda/internal/config"
	"github.com/jasonsoft/wakanda/internal/middleware"
	messengerHttp "github.com/jasonsoft/wakanda/pkg/messenger/delivery/http"
)

var (
	_messengerHandler *messengerHttp.MessengerHandler
)

func Initialize(config *config.Configuration) {

	// setup log
	log.SetAppID("messenger") // unique id for the app
	clog := console.New()
	log.RegisterHandler(clog, log.AllLevels...)

	//contactRepo := NewContactRepo()

}

func napWithMiddlewares() *napnap.NapNap {
	nap := napnap.New()
	nap.Use(napnap.NewHealth())
	nap.Use(middleware.NewErrorHandingMiddleware())
	corsOpts := napnap.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"*"},
	}
	nap.Use(napnap.NewCors(corsOpts))
	nap.Use(middleware.NewIdentityMiddleware())
	nap.Use(messengerHttp.NewRouter(_messengerHandler))
	return nap
}

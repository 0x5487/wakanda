package main

import (
	"fmt"
	"time"

	"github.com/cenk/backoff"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/console"
	"github.com/jasonsoft/log/handlers/gelf"
	"github.com/jasonsoft/napnap"
	"github.com/jasonsoft/wakanda/internal/config"
	"github.com/jasonsoft/wakanda/internal/identity"
	"github.com/jasonsoft/wakanda/internal/middleware"
	messengerHttp "github.com/jasonsoft/wakanda/pkg/messenger/delivery/http"
	messengerCockroachdb "github.com/jasonsoft/wakanda/pkg/messenger/repository/cockroachdb"
	messengerSvc "github.com/jasonsoft/wakanda/pkg/messenger/service"
	"github.com/jmoiron/sqlx"
)

var (
	_messengerHandler *messengerHttp.MessengerHandler
	_dbx              *sqlx.DB
)

func initialize(config *config.Configuration) {
	initLogger("messenger", config)
	initDatabase(config)

	contactRepo := messengerCockroachdb.NewContactRepo(_dbx)
	groupRepo := messengerCockroachdb.NewGroupRepo(_dbx)
	conversationRepo := messengerCockroachdb.NewConversationRepo(_dbx)

	contactSvc := messengerSvc.NewContactService(contactRepo, groupRepo, conversationRepo)
	groupSvc := messengerSvc.NewGroupService(groupRepo)
	conversationSvc := messengerSvc.NewConverstationService(conversationRepo)

	_messengerHandler = messengerHttp.NewMessengerHandler(contactSvc, groupSvc, conversationSvc)

}

func initLogger(appID string, config *config.Configuration) {
	// set up log target
	log.SetAppID(appID)
	for _, target := range config.Logs {
		switch target.Type {
		case "console":
			clog := console.New()
			levels := log.GetLevelsFromMinLevel(target.MinLevel)
			log.RegisterHandler(clog, levels...)
		case "gelf":
			graylog := gelf.New(target.ConnectionString)
			levels := log.GetLevelsFromMinLevel(target.MinLevel)
			log.RegisterHandler(graylog, levels...)
		}
	}
}

func initDatabase(config *config.Configuration) {
	if len(config.Database.DBName) == 0 {
		return
	}

	var err error
	var connectionString string

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(30) * time.Second

	connectionString = fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=require", config.Database.Username, config.Database.Password, config.Database.Address, config.Database.DBName)
	if err = backoff.Retry(func() error {
		_dbx, err = sqlx.Open("postgres", connectionString)
		if err != nil {
			panic(err)
		}
		err = _dbx.Ping()
		if err != nil {
			log.Errorf("main: cockroachdb ping error: %v", err)
			return err
		}
		return nil
	}, bo); err != nil {
		log.Panicf("cockroachdb connect timeout: %s", err.Error())
	}

	log.Infof("%s ping success", _dbx.DriverName())
	_dbx.SetMaxIdleConns(150)
	_dbx.SetMaxOpenConns(300)
	_dbx.SetConnMaxLifetime(14400 * time.Second)
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
	nap.Use(identity.NewMiddleware())
	nap.Use(messengerHttp.NewRouter(_messengerHandler))
	return nap
}

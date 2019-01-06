package main

import (
	"fmt"
	"os"
	"time"

	"github.com/cenk/backoff"
	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/console"
	"github.com/jasonsoft/log/handlers/gelf"
	"github.com/jasonsoft/napnap"
	"github.com/jasonsoft/wakanda/internal/config"
	"github.com/jasonsoft/wakanda/internal/identity"
	"github.com/jasonsoft/wakanda/internal/middleware"
	messengerGRPC "github.com/jasonsoft/wakanda/pkg/messenger/delivery/grpc"
	messengerHttp "github.com/jasonsoft/wakanda/pkg/messenger/delivery/http"
	messengerNats "github.com/jasonsoft/wakanda/pkg/messenger/delivery/nats"
	messengerCockroachdb "github.com/jasonsoft/wakanda/pkg/messenger/repository/cockroachdb"
	messengerSvc "github.com/jasonsoft/wakanda/pkg/messenger/service"
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/go-nats-streaming"
)

var (
	_messengerHandler *messengerHttp.MessengerHandler
	_messageSvc       *messengerSvc.MessageService
	_messageServer    *messengerGRPC.MessageServer
)

func initialize(config *config.Configuration) error {
	initLogger("messenger", config)

	dbx, err := setupDatabase(config)
	if err != nil {
		return err
	}

	// repository
	contactRepo := messengerCockroachdb.NewContactRepo(dbx)
	groupRepo := messengerCockroachdb.NewGroupRepo(dbx)
	groupMemberRepo := messengerCockroachdb.NewGroupMemberRepo(dbx)
	conversationRepo := messengerCockroachdb.NewConversationRepo(dbx)
	messageRepo := messengerCockroachdb.NewMessageRepo(dbx)

	// services
	contactSvc := messengerSvc.NewContactService(contactRepo, groupRepo, groupMemberRepo, conversationRepo)
	groupSvc := messengerSvc.NewGroupService(groupRepo, groupMemberRepo)
	conversationSvc := messengerSvc.NewConverstationService(conversationRepo)
	messageSvc := messengerSvc.NewMessageService(messageRepo, groupRepo)

	natsConn, err := setupNatsConn(config)
	if err != nil {
		return err
	}

	messagePub := messengerNats.NewMessagePublisher(natsConn)
	_messageServer = messengerGRPC.NewMessageServer(messageSvc, messagePub)
	_messengerHandler = messengerHttp.NewMessengerHandler(contactSvc, groupSvc, conversationSvc)

	return nil
}

func setupNatsConn(config *config.Configuration) (stan.Conn, error) {
	hostname, _ := os.Hostname()
	clientID := "messenger-" + hostname
	natsConn, err := stan.Connect(config.Nats.ClusterID, clientID, stan.NatsURL("nats://"+config.Nats.Address))
	if err != nil {
		return nil, err
	}
	return natsConn, nil
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

func setupDatabase(config *config.Configuration) (*sqlx.DB, error) {
	if len(config.Database.DBName) == 0 {
		return nil, nil
	}

	var dbx *sqlx.DB
	var err error
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(30) * time.Second
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=require", config.Database.Username, config.Database.Password, config.Database.Address, config.Database.DBName)

	if err = backoff.Retry(func() error {
		dbx, err = sqlx.Open("postgres", connectionString)
		if err != nil {
			panic(err)
		}
		err = dbx.Ping()
		if err != nil {
			log.Errorf("main: cockroachdb ping error: %v", err)
			return err
		}
		return nil
	}, bo); err != nil {
		log.Panicf("cockroachdb connect timeout: %s", err.Error())
	}

	log.Infof("%s ping success", dbx.DriverName())
	dbx.SetMaxIdleConns(150)
	dbx.SetMaxOpenConns(300)
	dbx.SetConnMaxLifetime(14400 * time.Second)
	return dbx, nil
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

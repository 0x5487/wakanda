package main

import (
	"fmt"
	"os"
	"time"

	"github.com/cenk/backoff"
	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/console"
	"github.com/jasonsoft/log/handlers/gelf"
	"github.com/jasonsoft/wakanda/internal/config"
	gatewayProto "github.com/jasonsoft/wakanda/pkg/gateway/proto"
	messengerNats "github.com/jasonsoft/wakanda/pkg/messenger/delivery/nats"
	messengerCockroachdb "github.com/jasonsoft/wakanda/pkg/messenger/repository/cockroachdb"
	messengerSvc "github.com/jasonsoft/wakanda/pkg/messenger/service"
	routerProto "github.com/jasonsoft/wakanda/pkg/router/proto"
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/go-nats-streaming"
	"google.golang.org/grpc"
)

var (
	deliverySub *messengerNats.DeliverySubscriber
)

func initialize(config *config.Configuration) error {
	initLogger("delivery", config)

	natsConn, err := setupNatsConn(config)
	if err != nil {
		return err
	}

	dbx, err := setupDatabase(config)
	if err != nil {
		return err
	}

	groupRepo := messengerCockroachdb.NewGroupRepo(dbx)
	groupMemberRepo := messengerCockroachdb.NewGroupMemberRepo(dbx)
	groupSvc := messengerSvc.NewGroupService(groupRepo, groupMemberRepo)

	// setup router client
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	routerConn, err := grpc.Dial(config.Router.AdvertiseAddr, opts...)
	if err != nil {
		log.Fatalf("delivery: can't connect to router grpc service: %v", err)
	}
	log.Info("delivery: router service was connected")
	routerClient := routerProto.NewRouterServiceClient(routerConn)

	// setup gateway client
	gatewayConn, err := grpc.Dial(config.Gateway.AdvertiseAddr, opts...)
	if err != nil {
		log.Fatalf("delivery: can't connect to router grpc service: %v", err)
	}
	log.Info("delivery: router service was connected")
	gatewayClient := gatewayProto.NewGatewayServiceClient(gatewayConn)

	deliverySub = messengerNats.NewDeliverySubscriber(natsConn, groupSvc, routerClient, gatewayClient)

	return nil
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

func setupNatsConn(config *config.Configuration) (stan.Conn, error) {
	hostname, _ := os.Hostname()
	clientID := "delivery-" + hostname
	natsConn, err := stan.Connect(config.Nats.ClusterID, clientID, stan.NatsURL("nats://"+config.Nats.Address))
	if err != nil {
		return nil, err
	}
	return natsConn, nil
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

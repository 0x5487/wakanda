package main

import (
	"fmt"
	"time"

	"github.com/cenk/backoff"
	"github.com/go-redis/redis"
	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/console"
	"github.com/jasonsoft/log/handlers/gelf"
	"github.com/jasonsoft/wakanda/internal/config"
	"github.com/jmoiron/sqlx"
)

var (
	_redisClient *redis.Client
)

func initialize(config *config.Configuration) {
	initLogger("router", config)

	_redisClient = setUpRedis(config)

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



func setUpRedis(config *config.Configuration) *redis.Client {
	var client *redis.Client
	var err error

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(30) * time.Second

	if err = backoff.Retry(func() error {

		client = redis.NewClient(&redis.Options{
			Addr:     config.Redis.Address,
			Password: config.Redis.Password,
			DB:       config.Redis.DB,
		})

		_, err := client.Ping().Result()
		if err != nil {
			return err
		}

		return nil
	}, bo); err != nil {
		log.Panicf("redis connect timeout: %s", err.Error())
	}
	return client
}

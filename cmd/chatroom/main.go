package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jasonsoft/log"
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
	err := initialize(config)
	if err != nil {
		log.Panicf("chatroom: initialize config fail: %v", err)
	}

	// start listen nats
	ctx := context.Background()

	_chatroomPubSub.SubscribeMSGChatroom(ctx)
	log.Info("chatroom: chatroom server started")

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
	<-stopChan
	log.Info("chatroom: chatroom shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := _chatroomPubSub.Shutdown(ctx); err != nil {
		log.Errorf("chatroom: service hanlder shutdown error: %v", err)
	} else {
		log.Info("chatroom: gracefully stopped")
	}
}

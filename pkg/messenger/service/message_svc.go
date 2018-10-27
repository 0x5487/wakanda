package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jasonsoft/wakanda/pkg/messenger"
	"github.com/nats-io/go-nats-streaming"
)

type MessageService struct {
	messageRepo messenger.MessageRepository
	groupRepo   messenger.GroupRepository
	messageChan chan *messenger.Message
	natsConn    stan.Conn
}

func NewMessageService(messageRepo messenger.MessageRepository, groupRepo messenger.GroupRepository, natsConn stan.Conn) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
		groupRepo:   groupRepo,
		messageChan: make(chan *messenger.Message, 5000),
		natsConn:    natsConn,
	}
}

func (svc *MessageService) startTasks() {
	timer := time.NewTicker(1 * time.Second)
	msgs := []*messenger.Message{}
	for {
		select {
		case msg := <-svc.messageChan:
			msgs = append(msgs, msg)
		case <-timer.C:
			if len(msgs) == 0 {
				continue
			}

			svc.messageRepo.BatchInsert(context.Background(), msgs)

			// to byte
			bytes, err := json.Marshal(msgs)
			if err != nil {

			}

			// send messages to delivery subject
			svc.natsConn.Publish("delivery", bytes)
		}
	}
}

func (svc *MessageService) Messages(ctx context.Context, opts *messenger.FindMessagesOptions) ([]*messenger.Message, error) {
	panic("not implemented")
}

func (svc *MessageService) CreateMessage(ctx context.Context, msg *messenger.Message) error {
	// ensure the message's timestamp is valid.
	nowUTC := time.Now().UTC()
	if msg.CreatedAt.Before(nowUTC.AddDate(0, 0, -3)) {
		return messenger.ErrMessageInvalid
	}
	if msg.CreatedAt.After(nowUTC) {
		msg.CreatedAt = &nowUTC
	}
	msg.UpdatedAt = &nowUTC

	// ensure the member in the group
	if svc.groupRepo.IsMemberInGroup(ctx, msg.SenderID, msg.GroupID) == false {
		return messenger.ErrMessageInvalid
	}

	// if group type is p2p, ensure friendship is available

	// save message
	svc.messageChan <- msg
	return nil
}

func (svc *MessageService) AckMessage(ctx context.Context, conversationID string, ackMsgID string, memberID string) error {
	// ensure the conversation belong to the member

	panic("not implemented")
}

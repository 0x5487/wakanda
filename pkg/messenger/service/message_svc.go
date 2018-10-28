package service

import (
	"context"
	"time"

	"github.com/jasonsoft/wakanda/pkg/messenger"
)

type MessageService struct {
	messageRepo messenger.MessageRepository
	groupRepo   messenger.GroupRepository
	messageChan chan *messenger.Message
}

func NewMessageService(messageRepo messenger.MessageRepository, groupRepo messenger.GroupRepository) *MessageService {
	svc := &MessageService{
		messageRepo: messageRepo,
		groupRepo:   groupRepo,
		messageChan: make(chan *messenger.Message, 5000),
	}

	go svc.startTasks()
	return svc
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
		}
	}
}

func (svc *MessageService) Messages(ctx context.Context, opts *messenger.FindMessagesOptions) ([]*messenger.Message, error) {
	panic("not implemented")
}

func (svc *MessageService) CreateMessage(ctx context.Context, msg *messenger.Message) error {
	nowUTC := time.Now().UTC()
	msg.CreatedAt = &nowUTC
	msg.UpdatedAt = &nowUTC

	// ensure the message's timestamp is valid.
	// if msg.CreatedAt.Before(nowUTC.AddDate(0, 0, -3)) {
	// 	return messenger.ErrMessageInvalid
	// }
	// if msg.CreatedAt.After(nowUTC) {
	// 	msg.CreatedAt = &nowUTC
	// }
	// msg.UpdatedAt = &nowUTC

	// ensure the member in the group
	// if svc.groupRepo.IsMemberInGroup(ctx, msg.SenderID, msg.GroupID) == false {
	// 	return messenger.ErrMessageInvalid
	// }

	// if group type is p2p, ensure friendship is available

	// save message
	svc.messageChan <- msg
	return nil
}

func (svc *MessageService) AckMessage(ctx context.Context, conversationID string, ackMsgID string, memberID string) error {
	// ensure the conversation belong to the member

	panic("not implemented")
}

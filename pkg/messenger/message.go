package messenger

import (
	"context"
	"time"
)

type MessageType int

const (
	MessageTypeNotification = 1
	MessageTypeText         = 2
	MessageTypeImage        = 3
)

type MessageState int

const (
	MessageStateNormal = 1
	MessageStateDelete = 2
	MessageStateDone   = 3 // 完成朋友邀請
)

type Message struct {
	ID        string
	GroupID   string
	SenderID  string
	Type      MessageType
	State     MessageState
	Content   string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type FindMessagesOptions struct {
	ID             string
	MemberID       string
	GroupID        string
	StartMessageID int64
	Size           int
}

type Messager interface {
	Messages(ctx context.Context, opts FindMessagesOptions) ([]*Message, error)
	CreateMessage(ctx context.Context, msg *Message) error
	AckMessage(ctx context.Context, ConversationID string, ackMsgID int64, memberID string) error
}

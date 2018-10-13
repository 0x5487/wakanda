package messenger

import (
	"context"
	"time"
)

type MessageType int

const (
	MessageTypeNotification MessageType = 1
	MessageTypeText         MessageType = 2
	MessageTypeImage        MessageType = 3
)

type MessageState int

const (
	MessageStateNormal MessageState = 1
	MessageStateDelete MessageState = 2
	MessageStateDone   MessageState = 3 // 完成朋友邀請
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
	StartMessageID string `db:"start_message_id"`
	Skip           int    `db:"skip"`
	PerPage        int    `db:"per_page"`
}

type Messager interface {
	Messages(ctx context.Context, opts *FindMessagesOptions) ([]*Message, error)
	SendMessage(ctx context.Context, msg *Message) error
	AckMessage(ctx context.Context, conversationID string, ackMsgID string, memberID string) error
}

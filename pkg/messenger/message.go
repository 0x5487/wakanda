package messenger

import "time"

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
	Messages(opts FindMessagesOptions) ([]*Message, error)
	CreateMessage(msg *Message) error
	AckMessage(ConversationID string, ackMsgID int64, memberID string) error
}

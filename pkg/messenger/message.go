package messenger

import (
	"context"
	"time"

	"github.com/jasonsoft/wakanda/internal/types"
	"github.com/jmoiron/sqlx"
)

var (
	ErrMessageInvalid = types.AppError{ErrorCode: "invalid_message", Message: "the message's is invalid"}
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
	ID              string
	RequestID       string
	GroupID         string
	SenderID        string
	SenderFirstName string
	SenderLastName  string
	Type            MessageType
	Content         string
	State           MessageState
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

type FindMessagesOptions struct {
	ID              string
	MemberID        string
	GroupIDs        []string
	AnchorUpdatedAt *time.Time
	Skip            int
	PerPage         int
}

type MessageServicer interface {
	Messages(ctx context.Context, opts *FindMessagesOptions) ([]*Message, error)
	CreateMessage(ctx context.Context, msg *Message) error
	AckMessage(ctx context.Context, conversationID string, ackMsgID string, memberID string) error
}

type MessageRepository interface {
	DB() *sqlx.DB
	BatchInsert(ctx context.Context, targets []*Message) error
}

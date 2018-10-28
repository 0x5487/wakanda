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
	ID              string       `json:"id" db:"id"`
	RequestID       string       `json:"request_id" db:"request_id"`
	GroupID         string       `json:"group_id" db:"group_id"`
	SenderID        string       `json:"sender_id" db:"sender_id"`
	SenderFirstName string       `json:"sender_first_name" db:"sender_first_name"`
	SenderLastName  string       `json:"sender_last_name" db:"sender_last_name"`
	Type            MessageType  `json:"type" db:"type"`
	Content         string       `json:"content" db:"content"`
	State           MessageState `json:"state" db:"state"`
	CreatedAt       *time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt       *time.Time   `json:"updated_at" db:"updated_at"`
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

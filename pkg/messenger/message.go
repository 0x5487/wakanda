package core

import "time"

const (
	// MessageTypeText represent text format of message
	MessageTextType  = 1
	MessageImageType = 2

	MessageNormalState  = 1
	MessageDeletedState = 2
)

type Message struct {
	ID        int64
	GroupID   string
	SenderID  string
	Type      int
	State     int
	Content   string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type FindMessagesOptions struct {
	MemberID      string
	FromMessageID int64
	Take          int
	Skip          int
}

type Messenger interface {
	Messages(opts FindMessagesOptions) ([]*Message, error)
	CreateMessage(msg *Message) error
	AckMessage(groupID string, ackMsgID int64) error
}

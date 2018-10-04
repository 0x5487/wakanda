package messenger

import "time"

type MessageType int

const (
	// MessageTypeText represent text format of message
	MessageTypeText  = 1
	MessageTypeImage = 2
)

type MessageState int

const (
	MessageStateNormal = 1
	MessageStateDelete = 2
)

type Message struct {
	ID             string
	ConversationID string
	SenderID       string
	Type           MessageType
	State          MessageState
	Content        string
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
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

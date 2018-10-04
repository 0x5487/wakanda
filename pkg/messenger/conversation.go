package messenger

import "time"

type ConversationType int

const ()

type Conversation struct {
	ID               string
	GroupID          string
	MemberID         string
	IsMute           bool
	LastAckMessageID string
	CreatedAt        *time.Time
	UpdatedAt        *time.Time
}

type FindConversionOptions struct {
	ID       string
	MemberID string
	SortBy   string
	OrderBy  string
}

type ConversationServicer interface {
	Conversations(opts FindConversionOptions) ([]*Conversation, error)
	CreateConversation(conversation *Conversation) error
	UnreadMessageCount(conversationID string) (int, error)
	MarkAllMessageAsRead(conversationID string) error
	GetConversationMessageCount(conversationID string) (int, error)
}

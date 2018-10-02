package messenger

import "time"

type Conversation struct {
	ID                 string
	GroupID            string
	MemberID           string
	IsMute             bool
	LatestAckMessageID string
	CreatedAt          *time.Time
	UpdatedAt          *time.Time
}

type FindConversionOptions struct {
	ID       string
	MemberID string
}

type ConversationServicer interface {
	Conversations(opts FindConversionOptions) ([]*Conversation, error)
	UnreadMessageCount(conversationID string) (int, error)
	MarkAllMessageAsRead(conversationID string) error
	GetConversationMessageCount() (int, error)
}

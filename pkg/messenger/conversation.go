package messenger

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Conversation struct {
	GroupID          string
	MemberID         string
	IsMute           bool
	LastAckMessageID string
	CreatedAt        *time.Time
	UpdatedAt        *time.Time
}

type FindConversionOptions struct {
	ID              string
	GroupID         string
	MemberID        string
	SortBy          string
	OrderBy         string
	AnchorUpdatedAt *time.Time
	Skip            int
	PerPage         int
}

type ConversationServicer interface {
	Conversations(ctx context.Context, opts *FindConversionOptions) ([]*Conversation, error)
	CreateConversation(ctx context.Context, conversation *Conversation) error
	UnreadMessageCount(ctx context.Context, conversationID string) (int, error)
	MarkAllMessageAsRead(ctx context.Context, conversationID string) error
	GetConversationMessageCount(ctx context.Context, conversationID string) (int, error)
}

type ConversationRepository interface {
	DB() *sqlx.DB
	Insert(ctx context.Context, target *Conversation, tx *sqlx.Tx) error
}

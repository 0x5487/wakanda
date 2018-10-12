package messenger

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Conversation struct {
	ID               string     `json:"id" db:"id"`
	GroupID          string     `json:"group_id" db:"group_id"`
	MemberID         string     `json:"member_id" db:"member_id"`
	IsMute           bool       `json:"is_mute" db:"is_mute"`
	LastAckMessageID string     `json:"last_ack_message_id" db:"last_ack_message_id"`
	CreatedAt        *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at" db:"updated_at"`
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
	InsertTx(ctx context.Context, target *Conversation, tx *sqlx.Tx) error
}

package messenger

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type ConversationState int

const (
	ConversationStateNormal ConversationState = 1
)

type Conversation struct {
	ID               string            `json:"id" db:"id"`
	GroupID          string            `json:"group_id" db:"group_id"`
	MemberID         string            `json:"-" db:"member_id"`
	IsMute           bool              `json:"is_mute" db:"is_mute"`
	LastAckMessageID string            `json:"last_ack_message_id" db:"last_ack_message_id"`
	State            ConversationState `json:"state" db:"state"`
	CreatedAt        *time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt        *time.Time        `json:"updated_at" db:"updated_at"`
}

type FindConversionOptions struct {
	ID              string
	GroupID         string
	MemberID        string `db:"member_id"`
	SortBy          string
	OrderBy         string
	AnchorUpdatedAt *time.Time `db:"anchor_updated_at"`
	Skip            int        `db:"skip"`
	PerPage         int        `db:"per_page"`
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
	Conversations(ctx context.Context, opts *FindConversionOptions) ([]*Conversation, error)
	InsertTx(ctx context.Context, target *Conversation, tx *sqlx.Tx) error
}

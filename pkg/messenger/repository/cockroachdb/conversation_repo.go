package cockroachdb

import (
	"context"

	"github.com/jasonsoft/wakanda/pkg/messenger"
	"github.com/jmoiron/sqlx"
)

type ConversationRepo struct {
	db *sqlx.DB
}

func NewConversationRepo(db *sqlx.DB) *ConversationRepo {
	return &ConversationRepo{
		db: db,
	}
}

func (repo *ConversationRepo) DB() *sqlx.DB {
	return repo.db
}

func (repo *ConversationRepo) Insert(ctx context.Context, target *messenger.Conversation, tx *sqlx.Tx) error {
	return nil
}

package cockroachdb

import (
	"context"

	"github.com/jasonsoft/wakanda/pkg/messenger"
	"github.com/jmoiron/sqlx"
)

type MessageRepo struct {
	db *sqlx.DB
}

func NewMessageRepo(db *sqlx.DB) *MessageRepo {
	return &MessageRepo{
		db: db,
	}
}

func (repo *MessageRepo) DB() *sqlx.DB {
	return repo.db
}

func (repo *MessageRepo) BatchInsert(ctx context.Context, targets []*messenger.Message) error {
	return nil
}

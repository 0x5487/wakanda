package cockroachdb

import (
	"context"

	"github.com/jasonsoft/wakanda/pkg/messenger"

	"github.com/jmoiron/sqlx"
)

type GroupRepo struct {
	db *sqlx.DB
}

func NewGroupRepo(db *sqlx.DB) *GroupRepo {
	return &GroupRepo{
		db: db,
	}
}

func (repo *GroupRepo) DB() *sqlx.DB {
	return repo.db
}

func (repo *GroupRepo) InsertTx(ctx context.Context, target *messenger.Group, tx *sqlx.Tx) error {
	return nil
}

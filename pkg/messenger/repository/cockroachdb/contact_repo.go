package cockroachdb

import (
	"context"

	"github.com/jasonsoft/wakanda/pkg/messenger"
	"github.com/jmoiron/sqlx"
)

type ContactRepo struct {
	db *sqlx.DB
}

func NewContactRepo(db *sqlx.DB) *ContactRepo {
	return &ContactRepo{
		db: db,
	}
}

func (repo *ContactRepo) DB() *sqlx.DB {
	return repo.db
}

func (repo *ContactRepo) Insert(ctx context.Context, target *messenger.Contact, tx *sqlx.Tx) error {
	return nil
}

func (repo *ContactRepo) Select(ctx context.Context, opts *messenger.FindContactOptions) ([]*messenger.Contact, error) {
	return nil, nil
}

func (repo *ContactRepo) Block(ctx context.Context, target *messenger.Contact) error {
	return nil
}

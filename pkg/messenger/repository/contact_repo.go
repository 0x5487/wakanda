package repository

import (
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

func (repo *ContactRepo) Insert(target *messenger.Contact) error {
	return nil
}

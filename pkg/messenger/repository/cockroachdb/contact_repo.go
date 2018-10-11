package cockroachdb

import (
	"context"

	"github.com/jasonsoft/log"
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

const insertContactSQL = `INSERT INTO messenger_contacts
(group_id, member_id_1, member_id_2, state)
VALUES(:group_id, :member_id_1, :member_id_2, :state);`

func (repo *ContactRepo) InsertTx(ctx context.Context, target *messenger.Contact, tx *sqlx.Tx) error {
	logger := log.FromContext(ctx)

	_, err := tx.NamedExec(insertContactSQL, target)
	if err != nil {
		logger.Errorf("cockroachdb: insert contact table failed: %v", err)
		return err
	}

	return nil
}

func (repo *ContactRepo) Contacts(ctx context.Context, opts *messenger.FindContactOptions) ([]*messenger.Contact, error) {
	return nil, nil
}

func (repo *ContactRepo) Block(ctx context.Context, target *messenger.Contact) error {
	return nil
}

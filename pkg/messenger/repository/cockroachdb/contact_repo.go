package cockroachdb

import (
	"context"
	"database/sql"

	"github.com/jasonsoft/wakanda/internal/types"

	"github.com/jasonsoft/wakanda/internal/cockroachdb"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/wakanda/pkg/messenger"
	"github.com/jmoiron/sqlx"
)

var (
	ErrContactExist = types.AppError{ErrorCode: "contact_exist", Message: "the contact already exists"}
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
		if cockroachdb.IsErrDBDuplicate(err) {
			return ErrContactExist
		}
		logger.Errorf("cockroachdb: insert messenger_contacts table failed: %v", err)
		return err
	}

	return nil
}

const listContactSQL = `select group_id, member_id_2 as member_id, state, created_at, updated_at from messenger_contacts where member_id_1 = :member_id and updated_at > :anchor_updated_at
union all
select group_id, member_id_1 as member_id, state, created_at, updated_at from messenger_contacts where member_id_2 = :member_id and updated_at > :anchor_updated_at`

func (repo *ContactRepo) Contacts(ctx context.Context, opts *messenger.FindContactOptions) ([]*messenger.Contact, error) {
	logger := log.FromContext(ctx)

	listContactSQLStmt, err := repo.db.PrepareNamed(listContactSQL)
	if err != nil {
		logger.Errorf("cockroachdb: prepare listContactSQL fail: %v", err)
		return nil, err
	}
	defer listContactSQLStmt.Close()

	contacts := []*messenger.Contact{}
	err = listContactSQLStmt.SelectContext(ctx, &contacts, opts)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Errorf("cockroachdb: get contacts fail: %v", err)
		return nil, err
	}

	return contacts, nil
}

func (repo *ContactRepo) Block(ctx context.Context, target *messenger.Contact) error {
	return nil
}

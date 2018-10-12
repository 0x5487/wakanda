package cockroachdb

import (
	"context"

	"github.com/jasonsoft/log"
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

const insertGroupSQL = `INSERT INTO messenger_groups
(id, "type", "name", description, max_member_count, member_count, state, creator_id)
VALUES(:id, :type, :name, :description, :max_member_count, :member_count, :state, :creator_id);`

func (repo *GroupRepo) InsertTx(ctx context.Context, target *messenger.Group, tx *sqlx.Tx) error {
	logger := log.FromContext(ctx)

	_, err := tx.NamedExec(insertGroupSQL, target)
	if err != nil {
		logger.Errorf("cockroachdb: insert group table failed: %v", err)
		return err
	}

	return nil
}

const listGroupSQL = `SELECT id, "type", "name", description, max_member_count, member_count, state, creator_id, created_at, updated_at 
FROM messenger_groups WHERE 1=1;`

func (repo *GroupRepo) Groups(ctx context.Context, opts *messenger.FindGroupOptions) ([]*messenger.Group, error) {
	return nil, nil
}

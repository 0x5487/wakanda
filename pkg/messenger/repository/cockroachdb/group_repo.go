package cockroachdb

import (
	"context"
	"strings"

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

	_, err := tx.NamedExecContext(ctx, insertGroupSQL, target)
	if err != nil {
		logger.Errorf("cockroachdb: insert group table failed: %v", err)
		return err
	}

	return nil
}

const listGroupSQL = `SELECT id, "type", "name", description, max_member_count, member_count, state, creator_id, created_at, updated_at 
FROM messenger_groups WHERE 1=1`

func (repo *GroupRepo) Groups(ctx context.Context, opts *messenger.FindGroupOptions) ([]*messenger.Group, error) {
	logger := log.FromContext(ctx)

	sqlBuilder := strings.Builder{}
	sqlBuilder.WriteString(listGroupSQL)

	if len(opts.IDs) > 0 {
		sqlBuilder.WriteString(" AND id in (:ids)")
	}

	sqlBuilder.WriteString(" LIMIT :per_page OFFSET :skip ;")

	prepareSQL, args, err := sqlx.Named(sqlBuilder.String(), opts)
	if err != nil {
		return nil, err
	}
	prepareSQL, args, err = sqlx.In(prepareSQL, args...)
	prepareSQL = repo.db.Rebind(prepareSQL)

	groups := []*messenger.Group{}
	repo.db.SelectContext(ctx, &groups, prepareSQL, args...)
	if err != nil {
		logger.Errorf("cockroachdb: get group list fail: %v", err)
		return nil, err
	}
	return groups, nil
}

func (repo *GroupRepo) IsMemberInGroup(ctx context.Context, memberID, groupID string) bool {
	panic("not implemented")
}

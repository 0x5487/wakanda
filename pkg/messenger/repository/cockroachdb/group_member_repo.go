package cockroachdb

import (
	"context"
	"database/sql"
	"strings"

	"github.com/jasonsoft/log"

	"github.com/jasonsoft/wakanda/pkg/messenger"
	"github.com/jmoiron/sqlx"
)

type GroupMemberRepo struct {
	db *sqlx.DB
}

func NewGroupMemberRepo(db *sqlx.DB) *GroupMemberRepo {
	return &GroupMemberRepo{
		db: db,
	}
}

func (repo *GroupMemberRepo) DB() *sqlx.DB {
	return repo.db
}

const listGroupMemberSQL = `SELECT id, group_id, member_id, is_admin, created_at, updated_at
FROM messenger_group_members WHERE 1=1`

func (repo *GroupMemberRepo) GroupMembers(ctx context.Context, opts *messenger.FindGroupMemberOptions) ([]*messenger.GroupMember, error) {
	logger := log.FromContext(ctx)

	sqlBuilder := strings.Builder{}
	sqlBuilder.WriteString(listGroupMemberSQL)

	if len(opts.GroupID) > 0 {
		log.Debugf("cockroachdb: group_id: %s", opts.GroupID)
		sqlBuilder.WriteString(" AND group_id = :group_id")
	}

	if !(opts.PerPage == 0 && opts.Skip == 0) {
		log.Debugf("cockroachdb: per_page: %d, skip: %d", opts.PerPage, opts.Skip)
		sqlBuilder.WriteString(" LIMIT :per_page OFFSET :skip;")
	}

	listGroupMemberStmt, err := repo.db.PrepareNamedContext(ctx, sqlBuilder.String())
	if err != nil {
		logger.Errorf("cockroachdb: prepare listGroupMemberSQL fail: %v", err)
		return nil, err
	}
	defer listGroupMemberStmt.Close()

	groupMembers := []*messenger.GroupMember{}
	err = listGroupMemberStmt.SelectContext(ctx, &groupMembers, opts)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Errorf("cockroachdb: get group members list fail: %v", err)
		return nil, err
	}
	return groupMembers, nil
}

func (repo *GroupMemberRepo) BatchInsertTx(ctx context.Context, members []*messenger.GroupMember, tx *sqlx.Tx) error {
	panic("not implemented")
}

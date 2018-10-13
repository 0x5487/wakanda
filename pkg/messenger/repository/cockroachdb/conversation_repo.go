package cockroachdb

import (
	"context"
	"database/sql"
	"strings"

	"github.com/jasonsoft/log"
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

const listConversationSQL = `SELECT id, group_id, member_id, is_mute, last_ack_message_id, state, created_at, updated_at
FROM messenger_conversations WHERE 1=1`

func (repo *ConversationRepo) Conversations(ctx context.Context, opts *messenger.FindConversionOptions) ([]*messenger.Conversation, error) {
	logger := log.FromContext(ctx)

	sqlBuilder := strings.Builder{}
	sqlBuilder.WriteString(listConversationSQL)

	if len(opts.MemberID) > 0 {
		sqlBuilder.WriteString(" AND member_id = :member_id")
		logger.Debugf("cockraochdb: conversation parameter member_id: %s", opts.MemberID)
	}

	if opts.AnchorUpdatedAt != nil {
		sqlBuilder.WriteString(" AND updated_at >= :anchor_updated_at")
		logger.Debugf("cockraochdb: conversation parameter anchor_updated_at: %s", opts.AnchorUpdatedAt)
	}

	sqlBuilder.WriteString(" ORDER BY updated_at DESC LIMIT :per_page OFFSET :skip;")

	listConversationSQLStmt, err := repo.db.PrepareNamedContext(ctx, sqlBuilder.String())
	if err != nil {
		logger.Errorf("cockroachdb: prepare listContactSQL fail: %v", err)
		return nil, err
	}
	defer listConversationSQLStmt.Close()

	conversations := []*messenger.Conversation{}
	err = listConversationSQLStmt.SelectContext(ctx, &conversations, opts)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Errorf("cockroachdb: get conversation fail: %v", err)
		return nil, err
	}

	return conversations, nil
}

const insertConversationSQL = `INSERT INTO messenger_conversations
(group_id, member_id, is_mute, last_ack_message_id, state)
VALUES(:group_id, :member_id, :is_mute, :last_ack_message_id, :state);`

func (repo *ConversationRepo) InsertTx(ctx context.Context, target *messenger.Conversation, tx *sqlx.Tx) error {
	logger := log.FromContext(ctx)

	target.LastAckMessageID = "00000000-0000-0000-0000-000000000000"
	_, err := tx.NamedExecContext(ctx, insertConversationSQL, target)
	if err != nil {
		logger.Errorf("cockroachdb: insert messenger_conversations table failed: %v", err)
		return err
	}

	return nil
}

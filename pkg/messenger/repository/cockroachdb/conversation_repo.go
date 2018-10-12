package cockroachdb

import (
	"context"

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

const insertConversationSQL = `INSERT INTO messenger_conversations
(group_id, member_id, is_mute, last_ack_message_id)
VALUES(:group_id, :member_id, :is_mute, :last_ack_message_id);`

func (repo *ConversationRepo) InsertTx(ctx context.Context, target *messenger.Conversation, tx *sqlx.Tx) error {
	logger := log.FromContext(ctx)

	target.LastAckMessageID = "00000000-0000-0000-0000-000000000000"
	_, err := tx.NamedExec(insertConversationSQL, target)
	if err != nil {
		logger.Errorf("cockroachdb: insert messenger_conversations table failed: %v", err)
		return err
	}

	return nil
}

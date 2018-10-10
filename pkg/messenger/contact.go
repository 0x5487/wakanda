package messenger

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type ContactState int

const (
	ContactStateNormal ContactState = 1
	ContactStateBlock  ContactState = 2
)

type Contact struct {
	ID        string       `json:"id"`
	MemberID1 string       `json:"member_id_1"`
	MemberID2 string       `json:"member_id_2"`
	State     ContactState `json:"state"`
	CreatedAt *time.Time   `json:"created_at"`
	UpdatedAt *time.Time   `json:"updated_at"`
}

type FindContactOptions struct {
	MemberID        string
	FriendID        string
	AnchorUpdatedAt *time.Time
	Size            int
}

type ContactServicer interface {
	Contacts(ctx context.Context, opts *FindContactOptions) ([]*Contact, error)
	BlockContact(ctx context.Context, memberID, friendID string) error
	AddContact(ctx context.Context, contact *Contact) error
}

type ContactRepository interface {
	DB() *sqlx.DB
	Select(ctx context.Context, opts *FindContactOptions) ([]*Contact, error)
	Insert(ctx context.Context, target *Contact, tx *sqlx.Tx) error
	Block(ctx context.Context, target *Contact) error
}

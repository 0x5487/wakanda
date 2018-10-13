package messenger

import (
	"context"
	"time"

	"github.com/jasonsoft/wakanda/internal/types"

	"github.com/jasonsoft/wakanda/internal/identity"
	"github.com/jmoiron/sqlx"
)

var (
	ErrContactExist = types.AppError{ErrorCode: "contact_exist", Message: "the contact already exists"}
)

type ContactState int

const (
	ContactStateNormal ContactState = 1
	ContactStateBlock  ContactState = 2
)

type Contact struct {
	GroupID   string           `json:"group_id" db:"group_id"`
	MemberID  string           `json:"member_id" db:"member_id"`
	MemberID1 string           `json:"-" db:"member_id_1"`
	MemberID2 string           `json:"-" db:"member_id_2"`
	Member    *identity.Member `json:"member"`
	State     ContactState     `json:"state" db:"state"`
	CreatedAt *time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time       `json:"updated_at" db:"updated_at"`
}

type FindContactOptions struct {
	MemberID        string     `db:"member_id"`
	AnchorUpdatedAt *time.Time `db:"anchor_updated_at"`
	Skip            int        `db:"skip"`
	PerPage         int        `db:"per_page"`
}

type ContactServicer interface {
	Contacts(ctx context.Context, opts *FindContactOptions) ([]*Contact, error)
	BlockContact(ctx context.Context, memberID, friendID string) error
	AddContact(ctx context.Context, memberID, friendID string) error
}

type ContactRepository interface {
	DB() *sqlx.DB
	Contacts(ctx context.Context, opts *FindContactOptions) ([]*Contact, error)
	InsertTx(ctx context.Context, target *Contact, tx *sqlx.Tx) error
	Block(ctx context.Context, target *Contact) error
}

package messenger

import "time"

type ContactState int

const (
	ContactStateNormal = 1
	ContactStateBlock  = 2
)

type Contact struct {
	MemberID   string
	FriendID   string
	FriendName string
	State      ContactState
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}

type FindContactOptions struct {
	MemberID        string
	FriendID        string
	AnchorUpdatedAt *time.Time
	Size            int
}

type ContactServicer interface {
	Contacts(opts FindContactOptions) ([]*Contact, error)
	DeleteContact(memberID, friendID string) error
	AddContact(memberID, friendID string) error
}

type ContactRepository interface {
	Select(opts *FindContactOptions) ([]*Contact, error)
	Insert(target *Contact) error
	Block(target *Contact) error
}

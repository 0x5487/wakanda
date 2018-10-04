package messenger

import "time"

type Contact struct {
	MemberID   string
	FriendID   string
	FriendName string
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}

type FindContactOptions struct {
	MemberID string
	FriendID string
}

type ContactServicer interface {
	Contacts(opts FindContactOptions) ([]*FindContactOptions, error)
	DeleteContact(memberID, friendID string) error
	AddContact(memberID, friendID string) error
}

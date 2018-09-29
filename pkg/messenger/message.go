package core

type Message struct {
	ID int64
}

type Messenger interface {
	MessagesByMemberID(memberID string, from_message_id int64, takes int) ([]*Message, error)
}

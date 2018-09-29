package core

import "time"

type Group struct {
	ID                 string
	Name               string
	Type               int // 1: group 2: channel
	IsMute             bool
	UnseenMessageCount int
	LatestAckMessageID int64
	CreatedAt          *time.Time
	UpdatedAt          *time.Time
}

type GroupServicer interface {
	ListGroups(memberID string) ([]*Group, error)
	AckMessage()
	AddGroup(groupID string) error
}

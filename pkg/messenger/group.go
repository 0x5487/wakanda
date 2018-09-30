package core

import "time"

const (
	GroupType   = 1
	ChannelType = 2
)

type Group struct {
	ID        string
	Name      string
	Type      int
	IconUrl   string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type GroupDetail struct {
	Group
	IsMute             bool
	UnseenMessageCount int
	LatestAckMessageID int64
	LatestMessage      *Message
}

type FindGroupsDetailOptions struct {
	ID       string
	MemberID string
}

type GroupServicer interface {
	GroupsDetails(opts FindGroupsDetailOptions) ([]*GroupDetail, error)
	AddGroup(groupID string, memberID string) error
}

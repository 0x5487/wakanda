package messenger

import "time"

type GroupType int

const (
	GroupTypeGroup   = 1
	GroupTypeChannel = 2
)

type Group struct {
	ID        string
	Name      string
	Type      GroupType
	IconUrl   string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type GroupServicer interface {
	CreateGroup(group Group) error
	LeaveGroup(groupID string) error
	AddGroupMember(groupID string, memberID string) error
	SetAdmin(groupID string, memberID string) error
	RemoveAdmin(groupID string, memberID string) error
	GroupAdmins(groupID string) ([]*Member, error)
	DissolveGroup(groupID string) error // 解散群組
}

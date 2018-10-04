package messenger

import "time"

type GroupType int

const (
	GroupTypePrivateOneByOne = 1
	GroupTypePublic          = 2
)

type Group struct {
	ID        string
	Name      string
	IsSystem  bool
	Type      GroupType
	IconPath  string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type GroupServicer interface {
	CreateGroup(group *Group) error
	DissolveGroup(groupID string) error // 解散群組
	JoinGroup(groupID string, memberID string) error
	LeaveGroup(groupID string) error
	AddGroupMember(groupID string, memberID string) error
	SetAdmin(groupID string, memberID string) error
	RemoveAdmin(groupID string, memberID string) error
	GroupAdmins(groupID string) ([]*Member, error)
}

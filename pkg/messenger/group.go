package messenger

import "time"

type GroupType int

const (
	GroupTypeSystem       = 1
	GroupTypeP2P          = 2
	GroupTypePrivateGroup = 3
)

type GroupState int

const (
	GroupStateNormal    = 1
	GroupStateDissolved = 2
)

type Group struct {
	ID             string
	Name           string
	Description    string
	Type           GroupType
	MaxMemberCount int
	MemberCount    int
	CreatorID      string
	State          GroupState
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
}

type FindGroupOptions struct {
	ID              string
	Type            GroupType
	MemberID        string
	AnchorUpdatedAt *time.Time
	Size            int
}

type FindGroupMemberOptions struct {
	GroupID         string
	AnchorUpdatedAt *time.Time
	Size            int
}

type GroupServicer interface {
	Groups(opts FindGroupOptions) ([]*Group, error)
	CreateGroup(group *Group, memberIDS []string) error
	DissolveGroup(groupID string) error // 解散群組
	JoinGroup(groupID string, memberID string) error
	LeaveGroup(groupID string) error
	AddGroupMember(groupID string, memberID string) error
	GroupMembers(opts FindGroupMemberOptions) ([]*Member, error)
	SetAdmin(groupID string, memberID string) error
	RemoveAdmin(groupID string, memberID string) error
	GroupAdmins(groupID string) ([]*Member, error)
}

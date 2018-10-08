package messenger

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

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
	Groups(ctx context.Context, opts FindGroupOptions) ([]*Group, error)
	CreateGroup(ctx context.Context, group *Group, memberIDS []string) error
	DissolveGroup(ctx context.Context, groupID string) error // 解散群組
	JoinGroup(ctx context.Context, groupID string, memberID string) error
	LeaveGroup(ctx context.Context, groupID string) error
	AddGroupMember(ctx context.Context, groupID string, memberID string) error
	GroupMembers(ctx context.Context, opts FindGroupMemberOptions) ([]*Member, error)
	SetAdmin(ctx context.Context, groupID string, memberID string) error
	RemoveAdmin(ctx context.Context, groupID string, memberID string) error
	GroupAdmins(ctx context.Context, groupID string) ([]*Member, error)
}

type GroupRepository interface {
	DB() *sqlx.DB
	CreateGroup(ctx context.Context, target *Group, memberIDS []string, tx *sqlx.Tx) error
}

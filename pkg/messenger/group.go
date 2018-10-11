package messenger

import (
	"context"
	"time"

	"github.com/jasonsoft/wakanda/internal/identity"
	"github.com/jmoiron/sqlx"
)

type GroupType int

const (
	GroupTypeSystem       GroupType = 1
	GroupTypeP2P          GroupType = 2
	GroupTypePrivateGroup GroupType = 3
)

type GroupState int

const (
	GroupStateNormal    GroupState = 1
	GroupStateDissolved GroupState = 2
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
	Skip            int
	PerPage         int
}

type FindGroupMemberOptions struct {
	GroupID         string
	AnchorUpdatedAt *time.Time
	Skip            int
	PerPage         int
}

type GroupServicer interface {
	Groups(ctx context.Context, opts *FindGroupOptions) ([]*Group, error)
	CreateGroup(ctx context.Context, group *Group, memberIDs []string) error
	DissolveGroup(ctx context.Context, groupID string) error // 解散群組
	JoinGroup(ctx context.Context, groupID string, memberID string) error
	LeaveGroup(ctx context.Context, groupID string) error
	AddGroupMember(ctx context.Context, groupID string, memberID string) error
	GroupMembers(ctx context.Context, opts *FindGroupMemberOptions) ([]*identity.Member, error)
	SetAdmin(ctx context.Context, groupID string, memberID string) error
	RemoveAdmin(ctx context.Context, groupID string, memberID string) error
	GroupAdmins(ctx context.Context, groupID string) ([]*identity.Member, error)
}

type GroupRepository interface {
	DB() *sqlx.DB
	CreateGroup(ctx context.Context, target *Group, memberIDs []string, tx *sqlx.Tx) error
}

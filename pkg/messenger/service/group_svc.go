package service

import (
	"context"

	"github.com/jasonsoft/wakanda/internal/mytime"

	"github.com/jasonsoft/wakanda/pkg/messenger"
)

type GroupService struct {
	groupRepo messenger.GroupRepository
}

func NewGroupService(groupRepo messenger.GroupRepository) *GroupService {
	return &GroupService{}
}

func (svc *GroupService) Groups(ctx context.Context, opts *messenger.FindGroupOptions) ([]*messenger.Group, error) {
	if opts.AnchorUpdatedAt == nil {
		opts.AnchorUpdatedAt = mytime.AnchorUpdateAt()
	}
	return svc.groupRepo.Groups(ctx, opts)
}

func (svc *GroupService) CreateGroup(ctx context.Context, group *messenger.Group, memberIDs []string) error {
	panic("not implemented")
}

func (svc *GroupService) DissolveGroup(ctx context.Context, groupID string) error {
	panic("not implemented")
}

func (svc *GroupService) JoinGroup(ctx context.Context, groupID string, memberID string) error {
	panic("not implemented")
}

func (svc *GroupService) LeaveGroup(ctx context.Context, groupID string) error {
	panic("not implemented")
}

func (svc *GroupService) AddGroupMember(ctx context.Context, groupID string, memberID string) error {
	panic("not implemented")
}

func (svc *GroupService) GroupMembers(ctx context.Context, opts *messenger.FindGroupMemberOptions) ([]*messenger.GroupMember, error) {
	panic("not implemented")
}

func (svc *GroupService) SetAdmin(ctx context.Context, groupID string, memberID string) error {
	panic("not implemented")
}

func (svc *GroupService) RemoveAdmin(ctx context.Context, groupID string, memberID string) error {
	panic("not implemented")
}

func (svc *GroupService) GroupAdmins(ctx context.Context, groupID string) ([]*messenger.GroupMember, error) {
	panic("not implemented")
}

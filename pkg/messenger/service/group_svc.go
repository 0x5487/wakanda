package service

import (
	"context"

	"github.com/jasonsoft/wakanda/pkg/messenger"
)

type GroupService struct {
	groupRepo messenger.GroupRepository
}

func NewGroupService() *GroupService {
	return &GroupService{}
}

func (svc *GroupService) Groups(ctx context.Context, opts *messenger.FindGroupOptions) ([]*messenger.Group, error) {
	return nil, nil
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

func (svc *GroupService) GroupMembers(ctx context.Context, opts messenger.FindGroupMemberOptions) ([]*messenger.Member, error) {
	panic("not implemented")
}

func (svc *GroupService) SetAdmin(ctx context.Context, groupID string, memberID string) error {
	panic("not implemented")
}

func (svc *GroupService) RemoveAdmin(ctx context.Context, groupID string, memberID string) error {
	panic("not implemented")
}

func (svc *GroupService) GroupAdmins(ctx context.Context, groupID string) ([]*messenger.Member, error) {
	panic("not implemented")
}

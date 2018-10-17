package service

import (
	"context"

	"github.com/jasonsoft/wakanda/pkg/messenger"
)

type DeliveryService struct {
	groupMemberRepo messenger.GroupMemberRepository
}

func NewDeliveryService(groupMemberRepo messenger.GroupMemberRepository) *DeliveryService {
	return &DeliveryService{
		groupMemberRepo: groupMemberRepo,
	}
}

func (svc *DeliveryService) DeliveryMessage(ctx context.Context, msg *messenger.Message) error {
	// opts := &messenger.FindGroupMemberOptions{
	// 	GroupID: msg.GroupID,
	// }

	// groupMembers, err := svc.groupMemberRepo.GroupMembers(ctx, opts)
	// if err != nil {
	// 	return err
	// }

	return nil
}

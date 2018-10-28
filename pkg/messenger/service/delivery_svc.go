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

func (svc *DeliveryService) DeliveryMessage(ctx context.Context, msgs *messenger.Message) error {
	// jobRequest := gatewayProto.SendJobRequest{}

	// for _, msg := range msgs {
	// 	opts := &messenger.FindGroupMemberOptions{
	// 		GroupID: msg.GroupID,
	// 	}

	// 	groupMembers, err := svc.groupMemberRepo.GroupMembers(ctx, opts)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	msgBytes, err := json.Marshal(msg)
	// 	if err != nil {
	// 		continue
	// 	}

	// 	command := gateway.Command{
	// 		OP:   "MSG",
	// 		Data: msgBytes,
	// 	}

	// 	for _, member := range groupMembers{
	// 		job := gatewayProto.Job{
	// 			Type: "s",
	// 			TargetID: member.
	// 		}
	// 	}
	// }

	// opts := &messenger.FindGroupMemberOptions{
	// 	GroupID: msg.GroupID,
	// }

	// groupMembers, err := svc.groupMemberRepo.GroupMembers(ctx, opts)
	// if err != nil {
	// 	return err
	// }

	return nil
}

package grpc

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/jasonsoft/log"
	chatroomProto "github.com/jasonsoft/wakanda/pkg/chatroom/proto"
	deliveryProto "github.com/jasonsoft/wakanda/pkg/dispatcher/proto"
)

func (svc *DispatcherServer) handleMSGRM(ctx context.Context, in *deliveryProto.DispatcherCommandRequest) (*deliveryProto.DispatcherCommandReply, error) {
	reqID := ""
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if len(md["req_id"]) > 0 {
			reqID = md["req_id"][0]
		}
	}
	customFields := log.Fields{
		"req_id": reqID,
	}
	logger := log.WithFields(customFields)

	logger.Infof("dispatcher: msg content: %s", string(in.Data))
	ctx = metadata.NewOutgoingContext(ctx, md)

	cmdChatroom := &chatroomProto.ChatroomCommand{
		RoomID:          in.TargetID,
		SenderID:        in.SenderID,
		SenderFirstName: in.SenderFirstName,
		SenderLastName:  in.SenderLastName,
		Data:            in.Data,
	}
	svc.dispatcherPub.PublishToChatMessageChannel(ctx, cmdChatroom)

	return &_emptyDispatcherCommandReply, nil
}

func (svc *DispatcherServer) handleTOKEN(ctx context.Context, in *deliveryProto.DispatcherCommandRequest) (*deliveryProto.DispatcherCommandReply, error) {
	reqID := ""
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if len(md["req_id"]) > 0 {
			reqID = md["req_id"][0]
		}
	}
	customFields := log.Fields{
		"req_id": reqID,
	}
	logger := log.WithFields(customFields)

	logger.Infof("dispatcher: msg content: %s", string(in.Data))
	ctx = metadata.NewOutgoingContext(ctx, md)

	cmdChatroom := &chatroomProto.ChatroomCommand{
		RoomID:          in.TargetID,
		SenderID:        in.SenderID,
		SenderFirstName: in.SenderFirstName,
		SenderLastName:  in.SenderLastName,
		Data:            in.Data,
	}
	svc.dispatcherPub.PublishToChatMessageChannel(ctx, cmdChatroom)

	return &_emptyDispatcherCommandReply, nil
}

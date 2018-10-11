package service

import (
	"context"

	"github.com/jasonsoft/wakanda/pkg/messenger"
)

type ConverstationService struct {
	converstationRepo messenger.ConversationRepository
}

func NewConverstationService(converstationRepo messenger.ConversationRepository) *ConverstationService {
	return &ConverstationService{
		converstationRepo: converstationRepo,
	}
}

func (svc *ConverstationService) Conversations(ctx context.Context, opts *messenger.FindConversionOptions) ([]*messenger.Conversation, error) {
	panic("not implemented")
}

func (svc *ConverstationService) CreateConversation(ctx context.Context, conversation *messenger.Conversation) error {
	panic("not implemented")
}

func (svc *ConverstationService) UnreadMessageCount(ctx context.Context, conversationID string) (int, error) {
	panic("not implemented")
}

func (svc *ConverstationService) MarkAllMessageAsRead(ctx context.Context, conversationID string) error {
	panic("not implemented")
}

func (svc *ConverstationService) GetConversationMessageCount(ctx context.Context, conversationID string) (int, error) {
	panic("not implemented")
}

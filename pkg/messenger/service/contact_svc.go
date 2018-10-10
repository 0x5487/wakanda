package service

import (
	"context"

	"github.com/jasonsoft/cockroach-go/crdb"
	"github.com/jasonsoft/wakanda/pkg/messenger"
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
)

type ContactService struct {
	contactRepo      messenger.ContactRepository
	groupRepo        messenger.GroupRepository
	conversationRepo messenger.ConversationRepository
}

func NewContactService(contactRepo messenger.ContactRepository, groupRepo messenger.GroupRepository, conversationRepo messenger.ConversationRepository) *ContactService {
	return &ContactService{
		contactRepo:      contactRepo,
		groupRepo:        groupRepo,
		conversationRepo: conversationRepo,
	}
}

func (svc *ContactService) Contacts(ctx context.Context, opts *messenger.FindContactOptions) ([]*messenger.Contact, error) {
	return svc.contactRepo.Select(ctx, opts)
}

func (svc *ContactService) AddContact(ctx context.Context, contact *messenger.Contact) error {
	err := crdb.ExecuteTx(ctx, svc.contactRepo.DB(), nil, func(tx *sqlx.Tx) error {

		// add contacts
		originalMemberID := contact.MemberID
		originalFriendID := contact.FriendID
		err := svc.contactRepo.Insert(ctx, contact, tx)
		if err != nil {
			return err
		}
		contact.FriendID = originalMemberID
		contact.MemberID = originalFriendID
		err = svc.contactRepo.Insert(ctx, contact, tx)
		if err != nil {
			return err
		}

		group := &messenger.Group{
			ID:             uuid.NewV4().String(),
			Type:           messenger.GroupTypeP2P,
			CreatorID:      contact.MemberID,
			MaxMemberCount: 2,
			MemberCount:    2,
			State:          messenger.GroupStateNormal,
		}

		memberIDs := []string{contact.MemberID, contact.FriendID}
		err = svc.groupRepo.CreateGroup(ctx, group, memberIDs, tx)
		if err != nil {
			return err
		}

		for _, memberID := range memberIDs {
			conversation := &messenger.Conversation{
				GroupID:  group.ID,
				MemberID: memberID,
			}

			err := svc.conversationRepo.Insert(ctx, conversation, tx)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (svc *ContactService) BlockContact(ctx context.Context, memberID, friendID string) error {
	panic("not implemented")
}

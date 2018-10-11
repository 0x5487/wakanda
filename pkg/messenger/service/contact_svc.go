package service

import (
	"context"

	"github.com/jasonsoft/wakanda/internal/hash"

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
	// TODO: allow max 100 per page

	// get contact
	contacts, err := svc.Contacts(ctx, opts)
	if err != nil {
		return nil, err
	}
	// get members

	return contacts, nil
}

func (svc *ContactService) AddContact(ctx context.Context, memberID, friendID string) error {
	err := crdb.ExecuteTx(ctx, svc.contactRepo.DB(), nil, func(tx *sqlx.Tx) error {
		hashNum1 := hash.FNV32a(memberID)
		hashNum2 := hash.FNV32a(friendID)

		contact := &messenger.Contact{}
		if hashNum1 < hashNum2 {
			contact.MemberID1 = memberID
			contact.MemberID2 = friendID
		} else {
			contact.MemberID1 = friendID
			contact.MemberID2 = memberID
		}

		contact.State = messenger.ContactStateNormal
		err := svc.contactRepo.Insert(ctx, contact, tx)
		if err != nil {
			return err
		}

		group := &messenger.Group{
			ID:             uuid.NewV4().String(),
			Type:           messenger.GroupTypeP2P,
			CreatorID:      memberID,
			MaxMemberCount: 2,
			MemberCount:    2,
			State:          messenger.GroupStateNormal,
		}

		memberIDs := []string{memberID, friendID}
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

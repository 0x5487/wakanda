package http

import (
	"time"

	"github.com/jasonsoft/wakanda/internal/pagination"

	"github.com/jasonsoft/napnap"
	"github.com/jasonsoft/wakanda/internal/identity"
	"github.com/jasonsoft/wakanda/internal/types"
	"github.com/jasonsoft/wakanda/pkg/messenger"
)

func (h *MessengerHandler) contactsMeListEndpoint(c *napnap.Context) {
	ctx := c.StdContext()
	pagination := pagination.FromContext(c)

	claim, found := identity.FromContext(ctx)
	if found == false {
		c.SetStatus(403)
		return
	}

	listContactOpts := &messenger.FindContactOptions{
		MemberID: claim.UserID,
		Skip:     pagination.Skip(),
		PerPage:  pagination.PerPage,
	}

	anchorUpdatedAtStr := c.Query("anchor_updated_at")
	if len(anchorUpdatedAtStr) > 0 {
		anchorUpdatedAt, err := time.Parse(time.RFC3339, anchorUpdatedAtStr)
		if err != nil {
			panic(types.AppError{ErrorCode: "invalid_input", Message: "anchor_updated_at field was invalid"})
		}
		listContactOpts.AnchorUpdatedAt = &anchorUpdatedAt
	}

	contacts, err := h.contactService.Contacts(ctx, listContactOpts)
	if err != nil {
		panic(err)
	}

	c.JSON(200, contacts)
}

type ContactCreatePayload struct {
	FriendID string `json:"friend_id"`
}

func (h *MessengerHandler) contactsMeCreateEndpoint(c *napnap.Context) {
	ctx := c.StdContext()

	claim, found := identity.FromContext(ctx)
	if found == false {
		c.SetStatus(403)
		return
	}

	payload := ContactCreatePayload{}
	err := c.BindJSON(&payload)
	if err != nil {
		panic(err)
	}

	err = h.contactService.AddContact(ctx, claim.UserID, payload.FriendID)
	if err != nil {
		panic(err)
	}

	c.SetStatus(201)
}

func (h *MessengerHandler) contactSearchEndpoint(c *napnap.Context) {

}

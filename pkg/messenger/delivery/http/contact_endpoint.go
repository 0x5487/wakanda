package http

import (
	"time"

	"github.com/jasonsoft/napnap"
	"github.com/jasonsoft/wakanda/internal/identity"
	"github.com/jasonsoft/wakanda/internal/types"
	"github.com/jasonsoft/wakanda/pkg/messenger"
)

func (h *MessengerHandler) contactsMeListEndpoint(c *napnap.Context) {
	ctx := c.StdContext()

	claim, found := identity.FromContext(ctx)
	if found == false {
		c.SetStatus(403)
		return
	}

	listContactOpts := &messenger.FindContactOptions{
		MemberID: claim.UserID,
	}

	anchorUpdatedAtStr := c.Query("anchor_updated_at")
	if len(anchorUpdatedAtStr) > 0 {
		anchorUpdatedAt, err := time.Parse(time.RFC3339, anchorUpdatedAtStr)
		if err != nil {
			panic(types.AppError{ErrorCode: "invalid_input", Message: "anchor_updated_at field was invalid"})
		}
		listContactOpts.AnchorUpdatedAt = &anchorUpdatedAt
	}

	contacts, err := h.ContactService.Contacts(ctx, listContactOpts)
	if err != nil {
		panic(err)
	}

	c.JSON(200, contacts)

}

func (h *MessengerHandler) contactsMeCreateEndpoint(c *napnap.Context) {
	ctx := c.StdContext()

	claim, found := identity.FromContext(ctx)
	if found == false {
		c.SetStatus(403)
		return
	}

	contact := &messenger.Contact{}
	err := c.BindJSON(contact)
	if err != nil {
		panic(err)
	}

	contact.MemberID = claim.UserID

	err = h.ContactService.AddContact(ctx, contact)
	if err != nil {
		panic(err)
	}

	c.SetStatus(201)
}

func (h *MessengerHandler) contactSearchEndpoint(c *napnap.Context) {

}

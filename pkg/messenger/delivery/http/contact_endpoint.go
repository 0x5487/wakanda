package http

import (
	"time"

	"github.com/jasonsoft/napnap"
	"github.com/jasonsoft/wakanda/internal/middleware"
	"github.com/jasonsoft/wakanda/internal/types"
	"github.com/jasonsoft/wakanda/pkg/messenger"
)

func (h *MessengerHandler) contactsListEndpoint(c *napnap.Context) {
	ctx := c.StdContext()

	claim, found := middleware.FromContext(ctx)
	if found == false {
		c.SetStatus(403)
		return
	}

	anchorUpdatedAtStr := c.Query("anchor_updated_at")
	anchorUpdatedAt, err := time.Parse(time.RFC3339, anchorUpdatedAtStr)
	if err != nil {
		panic(types.AppError{ErrorCode: "invalid_input", Message: "anchor_updated_at field was invalid"})
	}

	listContactOpts := messenger.FindContactOptions{
		MemberID:        claim.UserID,
		AnchorUpdatedAt: &anchorUpdatedAt,
	}

	contacts, err := h.ContactService.Contacts(&listContactOpts)
	if err != nil {
		panic(err)
	}

	c.JSON(200, contacts)

}

func (h *MessengerHandler) contactsCreateEndpoint(c *napnap.Context) {
	ctx := c.StdContext()

	claim, found := middleware.FromContext(ctx)
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

	err = h.ContactService.AddContact(contact)
	if err != nil {
		panic(err)
	}

	c.SetStatus(201)
}

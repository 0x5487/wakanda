package http

import (
	"github.com/jasonsoft/napnap"
	"github.com/jasonsoft/wakanda/pkg/identity"
)

func NewIdentityRouter(h *IdentityHttpHandler) *napnap.Router {
	router := napnap.NewRouter()
	router.Get("/tokens/:token_id", h.tokenGetEndpoint)
	return router
}

type IdentityHttpHandler struct {
	accountSvc identity.AccountServicer
}

func NewIdentityHttpHandler(accountSvc identity.AccountServicer) *IdentityHttpHandler {
	return &IdentityHttpHandler{
		accountSvc: accountSvc,
	}
}

func (h *IdentityHttpHandler) tokenGetEndpoint(c *napnap.Context) {

	// // hard-coding for test purpose
	// switch token {
	// case "4d96f463-dc14-44f0-af4f-c284e15c89cc":
	// 	stdctx := c.StdContext()
	// 	claim := types.Claim{
	// 		UserID:    "4d96f463-dc14-44f0-af4f-c284e15c89cc",
	// 		Username:  "angela",
	// 		Firstname: "Angela",
	// 		Lastname:  "Wang",
	// 	}
	// 	ctx := NewContext(stdctx, &claim)
	// 	c.SetStdContext(ctx)
	// case "aa58c0a6-32e3-4621-bb43-f45754f9f3dd":
	// default:
	// 	stdctx := c.StdContext()
	// 	claim := types.Claim{
	// 		UserID:    "aa58c0a6-32e3-4621-bb43-f45754f9f3dd",
	// 		Username:  "jason",
	// 		Firstname: "Jason",
	// 		Lastname:  "Lee",
	// 	}
	// 	ctx := NewContext(stdctx, &claim)
	// 	c.SetStdContext(ctx)
	// }

	account := identity.Account{
		ID:        "aa58c0a6-32e3-4621-bb43-f45754f9f3dd",
		Username:  "jason",
		FirstName: "Jason",
		LastName:  "Lee",
	}

	c.JSON(200, account)

}

package http

import (
	"github.com/jasonsoft/napnap"
	"github.com/jasonsoft/wakanda/pkg/messenger"
)

type MessengerHandler struct {
	ContactService messenger.ContactServicer
}

func NewRouter(h *MessengerHandler) *napnap.Router {
	router := napnap.NewRouter()
	router.Get("/v1/me/contacts", h.contactsListEndpoint)
	router.Post("/v1/me/contacts", h.contactsCreateEndpoint)
	router.Get("/v1/me/messages", h.meMessageListEndpoint)
	router.Get("/v1/me/groups", h.meGroupListEndpoint)
	router.Post("/v1/me/groups/:id/join", h.groupJoinEndpoint)
	return router
}

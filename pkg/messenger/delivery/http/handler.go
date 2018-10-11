package http

import (
	"github.com/jasonsoft/napnap"
	"github.com/jasonsoft/wakanda/pkg/messenger"
)

type MessengerHandler struct {
	contactService      messenger.ContactServicer
	groupService        messenger.GroupServicer
	conversationService messenger.ConversationServicer
}

func NewMessengerHandler(contactService messenger.ContactServicer, groupService messenger.GroupServicer, conversationService messenger.ConversationServicer) *MessengerHandler {
	return &MessengerHandler{
		contactService:      contactService,
		groupService:        groupService,
		conversationService: conversationService,
	}
}

func NewRouter(h *MessengerHandler) *napnap.Router {
	router := napnap.NewRouter()

	// contact
	router.Get("/v1/contacts", h.contactSearchEndpoint)
	router.Get("/v1/me/contacts", h.contactsMeListEndpoint)
	router.Post("/v1/me/contacts", h.contactsMeCreateEndpoint)

	// group
	router.Get("/v1/me/groups", h.meGroupListEndpoint)
	router.Post("/v1/me/groups/:id/join", h.groupJoinEndpoint)

	// conversation
	router.Get("/v1/me/conversations", h.conversationMeListEndpoint)

	// message
	router.Get("/v1/me/messages", h.messageMeListEndpoint)

	return router
}

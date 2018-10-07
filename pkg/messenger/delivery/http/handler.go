package http

import (
	"github.com/jasonsoft/napnap"
)

func NewRouter() *napnap.Router {
	router := napnap.NewRouter()
	router.Post("/v1/me/contacts", contactsCreateEndpoint)
	router.Get("/v1/me/messages", meMessageListEndpoint)
	router.Get("/v1/me/groups", meGroupListEndpoint)
	router.Get("/v1/groups/:id/join", groupJoinEndpoint)
	return router
}

type MessengerHandler struct {
}









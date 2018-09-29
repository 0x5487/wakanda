package http

import (
	"github.com/jasonsoft/napnap"
)

func NewMessengerRouter() *napnap.Router {
	router := napnap.NewRouter()
	router.Get("/v1/me/messages", meMessagesListEndpoint)
	router.Get("/v1/me/group", meGroupListEndpoint)
	router.Get("/v1/groups/:id/join", groupJoinEndpoint)
	return router
}

type MessengerHandler struct {
}

func meMessagesListEndpoint(c *napnap.Context) {

}

func meGroupListEndpoint(c *napnap.Context) {

}

func groupJoinEndpoint(c *napnap.Context) {

}

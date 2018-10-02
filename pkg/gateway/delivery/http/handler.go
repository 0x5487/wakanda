package http

import (
	"net/http"

	"github.com/jasonsoft/log"
	uuid "github.com/satori/go.uuid"

	"github.com/jasonsoft/wakanda/internal/types"
	"github.com/jasonsoft/wakanda/pkg/gateway"

	"github.com/gorilla/websocket"
	"github.com/jasonsoft/napnap"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func NewGatewayRouter() *napnap.Router {
	router := napnap.NewRouter()
	router.Get("/ws", wsEndpoint)
	return router
}

type GatewayHandler struct {
}

func wsEndpoint(c *napnap.Context) {
	defer func() {
		log.Debug("socket end.")
	}()
	member := &types.Member{}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	sessionID := uuid.NewV4().String()
	wsSession := gateway.NewWSSession(sessionID, member, conn)
	wsSession.StartTasks()
}

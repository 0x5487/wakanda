package http

import (
	"net/http"
	"sync/atomic"

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

	seedSessionID uint64
)

func NewGatewayRouter(handler *GatewayHandler) *napnap.Router {
	router := napnap.NewRouter()
	router.Get("/ws", handler.wsEndpoint)
	return router
}

type GatewayHandler struct {
}

func (h *GatewayHandler) wsEndpoint(c *napnap.Context) {

	member := &types.Member{}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	sessionID := atomic.AddUint64(&seedSessionID, 1)
	wsSession := gateway.NewWSSession(sessionID, member, conn)
	wsSession.StarTasks()
}

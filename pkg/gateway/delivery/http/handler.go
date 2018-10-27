package http

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jasonsoft/log"
	"github.com/jasonsoft/napnap"
	"github.com/jasonsoft/wakanda/internal/types"
	"github.com/jasonsoft/wakanda/pkg/dispatcher/proto"
	"github.com/jasonsoft/wakanda/pkg/gateway"
	"github.com/satori/go.uuid"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func NewGatewayRouter(h *GatewayHttpHandler) *napnap.Router {
	router := napnap.NewRouter()
	router.Get("/ws", h.wsEndpoint)
	return router
}

type GatewayHttpHandler struct {
	manager          *gateway.Manager
	dispatcherClient proto.DispatcherClient
}

func NewGatewayHttpHandler(manager *gateway.Manager, dispatcherClient proto.DispatcherClient) *GatewayHttpHandler {
	return &GatewayHttpHandler{
		manager:          manager,
		dispatcherClient: dispatcherClient,
	}
}

func (h *GatewayHttpHandler) wsEndpoint(c *napnap.Context) {
	defer func() {
		log.Debug("gateway: ws socket endpoint end")
	}()
	member := &types.Member{}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	sessionID := uuid.NewV4().String()
	wsSession := gateway.NewWSSession(sessionID, member, conn, h.manager, h.dispatcherClient)
	wsSession.StartTasks()
}

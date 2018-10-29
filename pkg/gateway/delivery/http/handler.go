package http

import (
	"net/http"

	"github.com/jasonsoft/wakanda/internal/identity"

	"github.com/gorilla/websocket"
	"github.com/jasonsoft/log"
	"github.com/jasonsoft/napnap"
	"github.com/jasonsoft/wakanda/pkg/dispatcher/proto"
	"github.com/jasonsoft/wakanda/pkg/gateway"
	routerProto "github.com/jasonsoft/wakanda/pkg/router/proto"
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
	dispatcherClient proto.DispatcherServiceClient
	routerClient     routerProto.RouterServiceClient
}

func NewGatewayHttpHandler(manager *gateway.Manager, dispatcherClient proto.DispatcherServiceClient, routerClient routerProto.RouterServiceClient) *GatewayHttpHandler {
	return &GatewayHttpHandler{
		manager:          manager,
		dispatcherClient: dispatcherClient,
		routerClient:     routerClient,
	}
}

func (h *GatewayHttpHandler) wsEndpoint(c *napnap.Context) {
	ctx := c.StdContext()

	defer func() {
		log.Debug("gateway: ws socket endpoint end")
	}()

	claim, found := identity.FromContext(ctx)
	if found == false {
		c.SetStatus(403)
		return
	}

	member := &identity.Member{
		ID:   claim.UserID,
		Name: claim.Username,
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	sessionID := uuid.NewV4().String()
	wsSession := gateway.NewWSSession(sessionID, member, conn, h.manager, h.dispatcherClient, h.routerClient)
	wsSession.StartTasks()
}

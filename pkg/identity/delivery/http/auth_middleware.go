package http

import (
	"fmt"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/napnap"
	"github.com/jasonsoft/request"
	"github.com/jasonsoft/wakanda/internal/config"
	"github.com/jasonsoft/wakanda/internal/types"
	"github.com/jasonsoft/wakanda/pkg/identity"
)

type AuthMiddleware struct {
	config *config.Configuration
}

func NewAuthMiddleware(config *config.Configuration) *AuthMiddleware {
	return &AuthMiddleware{
		config: config,
	}
}

func (m *AuthMiddleware) Invoke(c *napnap.Context, next napnap.HandlerFunc) {
	//ctx := c.StdContext()
	//logger := log.FromContext(ctx)

	token := c.RequestHeader("Authorization")
	if len(token) == 0 {
		token = c.Query("token")
	}

	if len(token) == 0 {
		c.SetStatus(401)
		return
	}

	// check token
	tokenAPIURL := fmt.Sprintf("%s/v1/tokens/%s", m.config.Identity.AdvertiseAddr, token)
	resp, err := request.GET(tokenAPIURL).End()

	if err != nil {
		log.Errorf("identity: get token network fail: %v, apiURL: %s", err, tokenAPIURL)
		panic(types.AppError{
			ErrorCode: "network_error",
			Message:   err.Error(),
		})
	}

	if !resp.OK {
		// unpredictable error
		log.Errorf("identity: get token resp fail: %v, apiURL: %s, resp: %s", err, tokenAPIURL, resp.String())
		appErr := types.AppError{
			ErrorCode: "request_error",
			Message:   resp.String(),
		}
		panic(appErr)
	}

	t := identity.Token{}
	resp.JSON(&t)

	stdctx := c.StdContext()
	ctx := identity.NewContext(stdctx, &t.Claims)
	c.SetStdContext(ctx)

	next(c)
}

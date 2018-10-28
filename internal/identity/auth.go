package identity

import (
	"context"

	"github.com/jasonsoft/napnap"
	"github.com/jasonsoft/wakanda/internal/types"
)

func NewContext(ctx context.Context, claim *types.Claim) context.Context {
	return context.WithValue(ctx, "identity_userkey", claim)
}

func FromContext(ctx context.Context) (*types.Claim, bool) {
	val, ok := ctx.Value("identity_userkey").(*types.Claim)
	if !ok {
		return nil, false
	}
	return val, true
}

type IdentityMiddleware struct {
}

func NewMiddleware() *IdentityMiddleware {
	return &IdentityMiddleware{}
}

func (l *IdentityMiddleware) Invoke(c *napnap.Context, next napnap.HandlerFunc) {
	//ctx := c.StdContext()
	//logger := log.FromContext(ctx)

	token := c.RequestHeader("Authorization")

	// hard-coding for test purpose
	switch token {
	case "4d96f463-dc14-44f0-af4f-c284e15c89cc":
		stdctx := c.StdContext()
		claim := types.Claim{
			UserID:   "4d96f463-dc14-44f0-af4f-c284e15c89cc",
			Username: "angela",
		}
		ctx := NewContext(stdctx, &claim)
		c.SetStdContext(ctx)
	case "aa58c0a6-32e3-4621-bb43-f45754f9f3dd":
	default:
		stdctx := c.StdContext()
		claim := types.Claim{
			UserID:   "aa58c0a6-32e3-4621-bb43-f45754f9f3dd",
			Username: "jason",
		}
		ctx := NewContext(stdctx, &claim)
		c.SetStdContext(ctx)
	}

	next(c)
}

package identity

import (
	"context"
)

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	Claim       Claim  `json:"claim,omitempty" db:"claim"`
}

type Claim struct {
	AccountID  string   `json:"account_id"`
	Username   string   `json:"username"`
	Firstname  string   `json:"first_name"`
	Lastname   string   `json:"last_name"`
	ConsumerID string   `json:"consumer_id"`
	Roles      []string `json:"roles"`
	Modules    []string `json:"modules"`
}

func NewContext(ctx context.Context, claim *Claim) context.Context {
	return context.WithValue(ctx, "identity_userkey", claim)
}

func FromContext(ctx context.Context) (*Claim, bool) {
	val, ok := ctx.Value("identity_userkey").(*Claim)
	if !ok {
		return nil, false
	}
	return val, true
}

type TokenServicer interface {
	Token(ctx context.Context, app, accessToken string) (*Token, error)
	DeleteToken(ctx context.Context, app, accessToken string) error
	CreateToken(ctx context.Context, login *LoginInfo) (*Token, error, int)
}

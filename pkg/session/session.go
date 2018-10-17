package session

import (
	"context"
	"time"
)

type SessionInfo struct {
	ID        string
	MemberID  string
	Server    string
	UpdatedAt *time.Time
}

type FindSessionInfo struct {
	MemberIDs []string
}

type SessionServicer interface {
	SessionInfo(sctx context.Context, opts *FindSessionInfo) ([]*SessionInfo, error)
	UpdateSessionInfo(ctx context.Context, target *SessionInfo) error
	DeleteSessionInfo(ctx context.Context, sessionID string) error
}

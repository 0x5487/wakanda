package identity

import (
	"context"
	"time"
)

type Role struct {
	ID        string     `json:"id,omitempty" db:"id"`
	App       string     `json:"app,omitempty" db:"app"`
	Name      string     `json:"name,omitempty" db:"name"`
	CreatedAt *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

type RoleServicer interface {
	Roles(ctx context.Context) ([]*Role, error)
	CreateRole(ctx context.Context, role *Role) error
	UpdateUserRole(ctx context.Context, app, accountID string, roles []string) error
}

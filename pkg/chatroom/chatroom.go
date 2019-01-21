package chatroom

import (
	"time"
)

type Room struct {
	ID         string
	MerchantID string
	Name       string
	CreatedAt  *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at" db:"updated_at"`
}

type ChatroomServicer interface {
}

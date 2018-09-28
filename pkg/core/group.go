package core

import "time"

type Group struct {
	ID        string
	Name      string
	Type      int
	IsMute    bool
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type GroupServicer struct {
}

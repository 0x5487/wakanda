package gateway

import "sync"

type Room struct {
	rwMutex  sync.RWMutex
	id       string
	sessions sync.Map
}

func (r *Room) Room(roomID string) *Room {

}

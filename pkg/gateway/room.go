package gateway

import (
	"sync"

	"github.com/jasonsoft/log"
)

type Room struct {
	id       string
	sessions sync.Map
}

func NewRoom(roomID string) *Room {
	return &Room{
		id: roomID,
	}
}

func (r *Room) join(session *WSSession) {
	r.sessions.Store(session.ID, session)
	session.rooms.Store(r.id, true)
	log.Infof("gateway: session id %d was joined to room id %s", session.ID, r.id)
}

func (r *Room) leave(session *WSSession) {
	r.sessions.Delete(session.ID)
	session.rooms.Delete(r.id)
	log.Infof("gateway: session id %d leaved the room id %s", session.ID, r.id)
}

func (r *Room) count() int {
	length := 0
	r.sessions.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	return length
}

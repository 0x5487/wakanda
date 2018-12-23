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
	log.Infof("gateway: session id %s was joined to room id %s", session.ID, r.id)
}

func (r *Room) leave(session *WSSession) {
	r.sessions.Delete(session.ID)
	session.rooms.Delete(r.id)
	log.Infof("gateway: session id %s leaved room id %s", session.ID, r.id)
}

func (r *Room) count() int {
	length := 0
	r.sessions.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	return length
}

func (r *Room) push(command *Command) {
	var (
		session *WSSession
		ok      bool
	)
	r.sessions.Range(func(key, value interface{}) bool {
		session, ok = value.(*WSSession)
		if ok {
			// msg, err := command.ToWSMessage()
			// if err != nil {
			// 	return true
			// }
			session.SendCommand(command)
			//session.SendMessage(msg)
		}
		return true
	})
}

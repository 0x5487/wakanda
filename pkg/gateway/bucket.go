package gateway

import (
	"sync"

	"github.com/jasonsoft/log"
)

type Bucket struct {
	id       int
	rooms    sync.Map
	sessions sync.Map
	jobChan  chan Job
}

func NewBucket(id, workerCount int) *Bucket {
	b := &Bucket{
		id:      id,
		jobChan: make(chan Job, 1000),
	}
	// create workers
	for i := 0; i < workerCount; i++ {
		go b.doJob()
	}
	return b
}

func (b *Bucket) addSession(session *WSSession) {
	b.sessions.Store(session.ID, session)
	log.Infof("gateway: session id %s was added to bucket id %d", session.ID, b.id)
}

func (b *Bucket) deleteSession(session *WSSession) {
	b.sessions.Delete(session.ID)
	log.Infof("gateway: session id %s was deleted to bucket id %d", session.ID, b.id)
}

func (b *Bucket) pushAll(command *Command) {
	var (
		session *WSSession
		ok      bool
	)

	b.sessions.Range(func(key, value interface{}) bool {
		session, ok = value.(*WSSession)
		if ok {
			msg, err := command.ToWSMessage()
			if err != nil {
				return true
			}
			session.SendMessage(msg)
		}
		return true
	})
}

func (b *Bucket) room(roomID string) *Room {
	room, found := b.rooms.Load(roomID)
	if found {
		room, ok := room.(*Room)
		if ok {
			return room
		}
	}
	return nil
}

func (b *Bucket) joinRoom(roomID string, session *WSSession) error {
	room := b.room(roomID)
	if room == nil {
		room = NewRoom(roomID)
		log.Infof("gateway: room id %s was created", room.id)
	}

	room.join(session)
	b.rooms.Store(roomID, room)
	return nil
}

func (b *Bucket) leaveRoom(roomID string, session *WSSession) {
	room := b.room(roomID)
	if room == nil {
		return
	}

	room.leave(session)
}

func (b *Bucket) pushRoom(roomID string, command *Command) {
	room := b.room(roomID)
	if room == nil {
		return
	}
	room.push(command)
}

func (b *Bucket) session(sessionID string) *WSSession {
	session, found := b.sessions.Load(sessionID)
	if found {
		session, ok := session.(*WSSession)
		if ok {
			return session
		}
	}
	return nil
}

func (b *Bucket) push(sessionID string, command *Command) {
	session := b.session(sessionID)
	msg, err := command.ToWSMessage()
	if err != nil {
		log.Errorf("gateway: command to message fail: %v", err)
		return
	}
	session.SendMessage(msg)
}

func (b *Bucket) count() int {
	length := 0
	b.sessions.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	return length
}

// doJob function which will be executed by workers.
func (b *Bucket) doJob() {
	var job Job
	for {
		select {
		case job = <-b.jobChan:
			switch job.OP {
			case OP_PUSH:
				b.push(job.SessionID, job.Command)
			case OP_PUSH_ALL:
				b.pushAll(job.Command)
			case OP_PUSH_ROOM:
				b.pushRoom(job.RoomID, job.Command)
			}
		}
	}
}

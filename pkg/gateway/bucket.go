package gateway

import (
	"sync"
)

type Bucket struct {
	rwMutex  sync.RWMutex
	id       int
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
}

func (b *Bucket) pushAll(message *WSMessage) {
	b.rwMutex.RLock()
	defer b.rwMutex.RUnlock()

	var (
		session *WSSession
		ok      bool
	)

	b.sessions.Range(func(key, value interface{}) bool {
		session, ok = value.(*WSSession)
		if ok {
			session.SendMessage(message)
		}
		return true
	})
}

func (b *Bucket) room(roomID string) *Room {
}

func (b *Bucket) pushRoom(roomID string, message *WSMessage) {
	b.rwMutex.RLock()
	defer b.rwMutex.RUnlock()

	var (
		session *WSSession
		ok      bool
	)

	b.sessions.Range(func(key, value interface{}) bool {
		session, ok = value.(*WSSession)
		if ok {
			session.SendMessage(message)
		}
		return true
	})
}

// doJob function which will be executed by workers.
func (b *Bucket) doJob() {
	var job Job
	for {
		select {
		case job = <-b.jobChan:
			switch job.OP {
			case PUSH_ALL:
				b.pushAll(job.WSMessage)
			case PUSH_ROOM:
				b.pushRoom(job.RoomID, job.WSMessage)
			}
		}
	}
}

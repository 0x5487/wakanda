package gateway

type Manager struct {
	buckets []*Bucket
}

func NewManager() *Manager {
	m := &Manager{
		buckets: make([]*Bucket, 1024),
	}
	for id, bucket := range m.buckets {
		bucket = NewBucket(id, 32)
	}
	return m
}

func (m *Manager) Bucket(sessionID uint64) *Bucket {
	return m.buckets[sessionID%uint64(len(m.buckets))]
}

func (m *Manager) AddSession(session *WSSession) {
	bucket := m.Bucket(session.ID)
	bucket.addSession(session)
}

func (m *Manager) PushAll(command *Command) {
	job := Job{
		OP:      PUSH_ALL,
		Command: command,
	}
	for _, bucket := range m.buckets {
		bucket.jobChan <- job
	}
}

func (m *Manager) PushRoom(roomID string, command *Command) {
	job := Job{
		OP:      PUSH_ROOM,
		Command: command,
	}
	for _, bucket := range m.buckets {
		bucket.jobChan <- job
	}
}

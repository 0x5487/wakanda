package gateway

import "sync/atomic"

type Status struct {
	OnlinePeople int64
}

func (s *Status) IncreaseOnlinePeople() {
	atomic.AddInt64(&s.OnlinePeople, 1)
}

func (s *Status) DecreaseOnlinePeople() {
	atomic.AddInt64(&s.OnlinePeople, -1)
}

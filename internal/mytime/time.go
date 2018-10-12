package mytime

import "time"

const (
	utc = "2006-01-02T15:04:05.000Z"
)

var (
	_anchorUpdateAt time.Time
)

func init() {
	str := "2018-01-01T11:45:26.371Z"
	t, _ := time.Parse(utc, str)
	_anchorUpdateAt = t
}

func AnchorUpdateAt() *time.Time {
	return &_anchorUpdateAt
}

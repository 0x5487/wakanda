package friendship

import "time"

type Friendship struct {
	MemberID  string
	FirstName string
	LastName  string
	LastSeen  *time.Time
}

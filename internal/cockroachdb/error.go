package cockroachdb

import (
	"github.com/lib/pq"
)

func IsErrDBDuplicate(err error) bool {
	if pgerr, ok := err.(*pq.Error); ok {
		if pgerr.Code == "23505" {
			return true
		}
	}

	return false
}

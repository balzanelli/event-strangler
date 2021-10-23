package eventstrangler

import "time"

type Record struct {
	ID        string
	Data      map[string]interface{}
	ExpiresAt time.Time
}

func (r *Record) isExpired() bool {
	return r.ExpiresAt.After(time.Now().UTC())
}

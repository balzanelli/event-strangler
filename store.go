package eventstrangler

import (
	"fmt"
	"time"
)

type StoreEventNotFoundError struct {
	HashKey string
}

func (e StoreEventNotFoundError) Error() string {
	return fmt.Sprintf("Hash Key '%s' Not Found", e.HashKey)
}

type Store interface {
	Exists(hashKey string) (bool, error)
	Get(hashKey string) (*Record, error)
	Put(hashKey string, record *Record, timeToLive int) error
	Delete(hashKey string) error
	Close() error
}

func getTimeToLive(timeToLive int) time.Duration {
	if timeToLive != 0 {
		return time.Second * time.Duration(timeToLive)
	}
	return time.Hour * 1
}

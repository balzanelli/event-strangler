package eventstrangler

import "fmt"

type StoreEventNotFoundError struct {
	HashKey string
}

func (e StoreEventNotFoundError) Error() string {
	return fmt.Sprintf("Hash Key '%s' Not Found", e.HashKey)
}

type Store interface {
	Exists(hashKey string) (bool, error)
	Get(hashKey string) (*Record, error)
	Put(hashKey string, record *Record) error
	Delete(hashKey string) error
	Close() error
}

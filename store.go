package eventstrangler

import (
	"encoding/json"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
)

type StoreEventNotFoundError struct {
	HashKey string
}

func (e StoreEventNotFoundError) Error() string {
	return fmt.Sprintf("Hash Key '%s' Not Found", e.HashKey)
}

type Store struct {
	db *leveldb.DB
}

func (s *Store) Exists(hashKey string) (bool, error) {
	return s.db.Has([]byte(hashKey), nil)
}

func (s *Store) Get(hashKey string) (*Record, error) {
	value, err := s.db.Get([]byte(hashKey), nil)
	if err != nil {
		return nil, StoreEventNotFoundError{
			HashKey: hashKey,
		}
	}

	var record Record
	if err := json.Unmarshal(value, &record); err != nil {
		return nil, err
	}
	return &record, nil
}

func (s *Store) Put(hashKey string, record *Record) error {
	serialized, err := json.Marshal(record)
	if err != nil {
		return err
	}
	return s.db.Put([]byte(hashKey), serialized, nil)
}

func (s *Store) Delete(hashKey string) error {
	return s.db.Delete([]byte(hashKey), nil)
}

func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func NewStore() (*Store, error) {
	db, err := leveldb.OpenFile("./tmp/strangler.db", nil)
	if err != nil {
		return nil, err
	}
	return &Store{
		db,
	}, nil
}

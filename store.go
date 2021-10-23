package eventstrangler

import (
	"encoding/json"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
)

type StoreEventError struct {
	HashKey string
}

func (e StoreEventError) Error() string {
	return fmt.Sprintf("Hash Key '%s' Not Found", e.HashKey)
}

type StoreEventNotFoundError struct {
	StoreEventError
}

type StoreEventAlreadyExistsError struct {
	StoreEventError
}

type Store struct {
	db *leveldb.DB
}

func (s *Store) Exists(hashKey string) (bool, error) {
	iter := s.db.NewIterator(nil, nil)
	defer iter.Release()
	return iter.Seek([]byte(hashKey)), nil
}

func (s *Store) GetEvent(hashKey string) (*Record, error) {
	iter := s.db.NewIterator(nil, nil)
	defer iter.Release()
	if !iter.Seek([]byte(hashKey)) {
		return nil, StoreEventNotFoundError{
			StoreEventError{HashKey: hashKey},
		}
	}

	var record Record
	if err := json.Unmarshal(iter.Value(), &record); err != nil {
		return nil, err
	}
	return &record, nil
}

func (s *Store) PutEvent(hashKey string, record *Record) error {
	exists, err := s.Exists(hashKey)
	if err != nil {
		return err
	}
	if exists {
		return StoreEventAlreadyExistsError{
			StoreEventError{HashKey: hashKey},
		}
	}

	serialized, err := json.Marshal(record)
	if err != nil {
		return err
	}
	return s.db.Put([]byte(hashKey), serialized, nil)
}

func (s *Store) DeleteEvent(hashKey string) error {
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

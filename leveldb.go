package eventstrangler

import (
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDBStore struct {
	db *leveldb.DB
}

func (s *LevelDBStore) Exists(hashKey string) (bool, error) {
	return s.db.Has([]byte(hashKey), nil)
}

func (s *LevelDBStore) Get(hashKey string) (*Record, error) {
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

func (s *LevelDBStore) Put(hashKey string, record *Record) error {
	serialized, err := json.Marshal(record)
	if err != nil {
		return err
	}
	return s.db.Put([]byte(hashKey), serialized, nil)
}

func (s *LevelDBStore) Delete(hashKey string) error {
	return s.db.Delete([]byte(hashKey), nil)
}

func (s *LevelDBStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func NewLevelDBStore() (*LevelDBStore, error) {
	db, err := leveldb.OpenFile("./tmp/strangler.db", nil)
	if err != nil {
		return nil, err
	}
	return &LevelDBStore{
		db: db,
	}, nil
}

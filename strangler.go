package eventstrangler

import (
	"time"
)

type Config struct {
	HashKey    *HashKeyOptions
	TimeToLive *time.Duration
}

type Strangler struct {
	Config *Config
	Store  Store
}

func (s *Strangler) Once(object map[string]interface{}) (string, error) {
	hashKey, err := getHashKey(object, s.Config.HashKey)
	if err != nil {
		return "", err
	}

	exists, err := s.Store.Exists(hashKey)
	if err != nil {
		return hashKey, err
	}
	if exists {
		record, err := s.Store.Get(hashKey)
		if err != nil {
			return hashKey, err
		}
		if record.isProcessing() {
			return hashKey, EventAlreadyProcessingError{
				HashKey:   hashKey,
				StartedAt: record.CreatedAt,
			}
		}
		return hashKey, EventAlreadyProcessedError{HashKey: hashKey}
	}

	record := newRecord(hashKey, s.Config.TimeToLive)
	if err = s.Store.Put(hashKey, record); err != nil {
		return hashKey, err
	}
	return hashKey, nil
}

func (s *Strangler) Complete(hashKey string) error {
	record, err := s.Store.Get(hashKey)
	if err != nil {
		return err
	}

	record.complete()
	if err = s.Store.Put(hashKey, record); err != nil {
		return err
	}
	return nil
}

func (s *Strangler) Fail(hashKey string) error {
	return s.Store.Delete(hashKey)
}

func NewStrangler(store Store, config *Config) (*Strangler, error) {
	if store == nil {
		return nil, StoreNotFoundError{}
	}
	return &Strangler{
		Config: config,
		Store:  store,
	}, nil
}

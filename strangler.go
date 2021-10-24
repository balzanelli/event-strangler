package eventstrangler

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"github.com/jmespath/go-jmespath"
	"time"
)

type Config struct {
	InstanceName      string
	HashKeyExpression string
	TimeToLive        *time.Duration
}

type Strangler struct {
	Config *Config
	Store  *Store
}

func (s *Strangler) Once(dict map[string]interface{}) (string, error) {
	identifierKey, err := s.getIdentifierKey(dict)
	if err != nil {
		return "", err
	}
	hashKey, err := s.getHashKey(identifierKey)
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

func (s *Strangler) getIdentifierKey(dict map[string]interface{}) ([]byte, error) {
	var identifierKey interface{}
	if len(s.Config.HashKeyExpression) != 0 {
		result, err := jmespath.Search(s.Config.HashKeyExpression, dict)
		if err != nil {
			return nil, err
		}
		identifierKey = result
	} else {
		identifierKey = dict
	}

	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(identifierKey); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (s *Strangler) getHashKey(identifierKey []byte) (string, error) {
	hash := sha256.New()
	if len(s.Config.InstanceName) != 0 {
		if _, err := hash.Write([]byte(s.Config.InstanceName)); err != nil {
			return "", err
		}
	}
	if _, err := hash.Write(identifierKey); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
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

package eventstrangler

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

type Config struct {
	InstanceName string
}

type Strangler struct {
	config *Config
	store  *Store
}

func (s *Strangler) Once(dict map[string]interface{}) (bool, error) {
	identifierKey, err := s.getIdentifierKey(dict)
	if err != nil {
		return false, err
	}

	hashKey, err := s.getHashKey(identifierKey)
	if err != nil {
		return false, err
	}
	log.Println(hashKey)
	return false, nil
}

func (s *Strangler) Complete() error {
	return nil
}

type IdentifierKeyNotFound struct {
	IdentifierKey string
}

func (e IdentifierKeyNotFound) Error() string {
	return fmt.Sprintf("Identifier Key '%s' Not Found", e.IdentifierKey)
}

const IdentifierKey = "transaction_id"

func (s *Strangler) getIdentifierKey(dict map[string]interface{}) ([]byte, error) {
	if transactionID, exists := dict[IdentifierKey]; exists {
		var buffer bytes.Buffer
		if err := gob.NewEncoder(&buffer).Encode(transactionID); err != nil {
			return nil, err
		}
		return buffer.Bytes(), nil
	}
	return nil, IdentifierKeyNotFound{IdentifierKey}
}

func (s *Strangler) getHashKey(identifierKey []byte) (string, error) {
	parts := [][]byte{identifierKey}
	if len(s.config.InstanceName) != 0 {
		parts = append(parts, []byte(s.config.InstanceName))
	}

	hash := sha256.New()
	if _, err := hash.Write(bytes.Join(parts, []byte("-"))); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func NewStrangler(config *Config) (*Strangler, error) {
	store, err := NewStore()
	if err != nil {
		return nil, err
	}
	return &Strangler{
		config: config,
		store:  store,
	}, nil
}

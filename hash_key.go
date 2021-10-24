package eventstrangler

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"github.com/jmespath/go-jmespath"
)

type HashKeyOptions struct {
	Name       string
	Expression string
}

func getHashKeyName(opt *HashKeyOptions) (bool, string) {
	if opt != nil && len(opt.Name) != 0 {
		return true, opt.Name
	}
	return false, ""
}

func getIdempotencyKey(object map[string]interface{}, opt *HashKeyOptions) ([]byte, error) {
	var idempotencyKey interface{}
	if opt != nil && len(opt.Expression) != 0 {
		result, err := jmespath.Search(opt.Expression, object)
		if err != nil {
			return nil, err
		}
		idempotencyKey = result
	} else {
		idempotencyKey = object
	}

	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(idempotencyKey); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func getHashKey(object map[string]interface{}, opt *HashKeyOptions) (string, error) {
	hash := sha256.New()
	if exists, name := getHashKeyName(opt); exists {
		if _, err := hash.Write([]byte(name)); err != nil {
			return "", err
		}
	}
	idempotencyKey, err := getIdempotencyKey(object, opt)
	if err != nil {
		return "", nil
	}
	if _, err := hash.Write(idempotencyKey); err != nil {
		return "", nil
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

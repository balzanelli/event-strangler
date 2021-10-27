package eventstrangler

import (
	"github.com/karlseguin/ccache/v2"
)

type LRUCacheStore struct {
	cache *ccache.Cache
}

func (s *LRUCacheStore) Exists(hashKey string) (bool, error) {
	_, err := s.Get(hashKey)
	if err != nil {
		if _, ok := err.(StoreEventNotFoundError); ok {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *LRUCacheStore) Get(hashKey string) (*Record, error) {
	item := s.cache.Get(hashKey)
	if item == nil {
		return nil, StoreEventNotFoundError{
			HashKey: hashKey,
		}
	}
	return item.Value().(*Record), nil
}

func (s *LRUCacheStore) Put(hashKey string, record *Record, timeToLive int) error {
	s.cache.Set(hashKey, record, getTimeToLive(timeToLive))
	return nil
}

func (s *LRUCacheStore) Delete(hashKey string) error {
	if !s.cache.Delete(hashKey) {
		return StoreEventNotFoundError{
			HashKey: hashKey,
		}
	}
	return nil
}

func (s *LRUCacheStore) Close() error {
	if s.cache != nil {
		s.cache.Clear()
	}
	return nil
}

func NewLRUCacheStore() *LRUCacheStore {
	return &LRUCacheStore{
		cache: ccache.New(ccache.Configure()),
	}
}

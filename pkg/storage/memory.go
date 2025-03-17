package storage

import (
	"sync"
	"time"
)

type KeyValueStore map[string]string

type InMemoryStorage struct {
	data       KeyValueStore
	expiration map[string]time.Time
	mu         sync.RWMutex
}

func NewInMemoryStorage() *InMemoryStorage {
	s := &InMemoryStorage{
		data:       make(map[string]string),
		expiration: make(map[string]time.Time),
	}

	go func() {
		for {
			time.Sleep(1 * time.Second)

			s.mu.Lock()
			for key, expiry := range s.expiration {
				if time.Now().After(expiry) {
					delete(s.data, key)
					delete(s.expiration, key)
				}
			}
			s.mu.Unlock()
		}
	}()

	return s
}

func (s *InMemoryStorage) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	expiry, exists := s.expiration[key]
	if exists && time.Now().After(expiry) {
		return "", false
	}

	value, ok := s.data[key]
	return value, ok
}
func (s *InMemoryStorage) Set(key string, value string, ttl int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value

	if ttl > 0 {
		s.expiration[key] = time.Now().Add(time.Duration(ttl) * time.Second)
	}

}

func (s *InMemoryStorage) Delete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.data[key]
	if ok {
		delete(s.data, key)
		delete(s.expiration, key)
	}
	return ok
}

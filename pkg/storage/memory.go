package storage

import (
	"context"
	"sync"
	"time"
)

type KeyValueStore map[string]string

type InMemoryStorage struct {
	data       KeyValueStore
	expiration map[string]time.Time
	mu         sync.RWMutex
	cancel     context.CancelFunc
}

func NewInMemoryStorage() *InMemoryStorage {
	s := &InMemoryStorage{
		data:       make(map[string]string),
		expiration: make(map[string]time.Time),
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.cleanupExpiredKeys()
			}
		}
	}()

	s.cancel = cancel

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

func (s *InMemoryStorage) cleanupExpiredKeys() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for key, expiry := range s.expiration {
		if time.Now().After(expiry) {
			delete(s.data, key)
			delete(s.expiration, key)
		}
	}
}

func (s *InMemoryStorage) Shutdown() {
	if s.cancel != nil {
		s.cancel()
	}
}

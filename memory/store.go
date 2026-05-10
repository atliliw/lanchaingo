package memory

import "sync"

// Store provides long-term memory across sessions.
type Store interface {
	Get(userID, key string) (any, bool)
	Put(userID, key string, value any)
	Delete(userID, key string)
	Keys(userID string) []string
}

// InMemoryStore implements Store in memory.
type InMemoryStore struct {
	mu   sync.RWMutex
	data map[string]map[string]any
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{data: make(map[string]map[string]any)}
}

func (s *InMemoryStore) Get(userID, key string) (any, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if user, ok := s.data[userID]; ok {
		v, ok := user[key]
		return v, ok
	}
	return nil, false
}

func (s *InMemoryStore) Put(userID, key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[userID]; !ok {
		s.data[userID] = make(map[string]any)
	}
	s.data[userID][key] = value
}

func (s *InMemoryStore) Delete(userID, key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if user, ok := s.data[userID]; ok {
		delete(user, key)
	}
}

func (s *InMemoryStore) Keys(userID string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if user, ok := s.data[userID]; ok {
		keys := make([]string, 0, len(user))
		for k := range user {
			keys = append(keys, k)
		}
		return keys
	}
	return nil
}

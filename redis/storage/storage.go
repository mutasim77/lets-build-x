package storage

import (
	"fmt"
	"strconv"
	"sync"
)

// Storage represents the in-memory key-value store
type Storage struct {
	data map[string]string // Internal map to store key-value pairs
	mu   sync.RWMutex      // Read-Write mutex for thread-safe operations
}

// NewStorage creates and returns a new Storage instance
func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

// Set stores a key-value pair in the storage
func (s *Storage) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// Get retrieves the value associated with the given key
// Returns the value and a boolean indicating if the key exists
func (s *Storage) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.data[key]
	return value, ok
}

// Del removes the specified key from the storage
// Returns true if the key was present and removed, false otherwise
func (s *Storage) Del(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.data[key]
	if ok {
		delete(s.data, key)
	}
	return ok
}

// Exists checks if the specified keys exist in the storage
// Returns the count of existing keys
func (s *Storage) Exists(keys ...string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	count := 0
	for _, key := range keys {
		if _, ok := s.data[key]; ok {
			count++
		}
	}
	return count
}

// IncrBy increments the value of the key by the given amount
// If the key doesn't exist, it's set to 0 before performing the operation
// Returns the new value and any error that occurred
func (s *Storage) IncrBy(key string, amount int) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.data[key]
	if !ok {
		s.data[key] = "0"
		value = "0"
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("value is not an integer")
	}

	newValue := intValue + amount
	s.data[key] = strconv.Itoa(newValue)

	return newValue, nil
}

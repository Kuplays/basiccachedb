package storage

import (
	"errors"
	"sync"
)

var (
	ErrorKeyNotFound = errors.New("Key not found")
	ErrorEmptyKey    = errors.New("Key cannot be empty")
)

type Storage struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

func (s *Storage) Set(key string, value string) error {
	if key == "" {
		return ErrorEmptyKey
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value

	return nil
}

func (s *Storage) Get(key string) (string, error) {
	if key == "" {
		return "", ErrorEmptyKey
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.data[key]
	if !ok {
		return "", ErrorKeyNotFound
	}

	return value, nil
}

func (s *Storage) Delete(key string) (bool, error) {
	if key == "" {
		return false, ErrorEmptyKey
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[key]; !ok {
		return false, nil
	}

	delete(s.data, key)

	return true, nil
}

func (s *Storage) Exists(key string) (bool, error) {
	if key == "" {
		return false, ErrorEmptyKey
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data[key]

	return ok, nil
}

func (s *Storage) Length() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data)
}

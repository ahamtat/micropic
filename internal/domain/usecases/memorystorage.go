package usecases

import (
	"sync"

	"github.com/AcroManiac/micropic/internal/domain/interfaces"
	"github.com/pkg/errors"
)

// MemoryStorage implementation
type MemoryStorage struct {
	mx   sync.RWMutex
	data map[string][]byte
}

// NewMemoryStorage constructor
func NewMemoryStorage() interfaces.Storage {
	return &MemoryStorage{
		mx:   sync.RWMutex{},
		data: make(map[string][]byte),
	}
}

// Save data in memory by hash key
func (s *MemoryStorage) Save(hash string, data []byte) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.data[hash] = data
	return nil
}

// Get stored data by hash key
func (s *MemoryStorage) Get(hash string) ([]byte, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	data, ok := s.data[hash]
	if !ok {
		return nil, errors.New("error finding value")
	}
	return data, nil
}

// Remove data from storage
func (s *MemoryStorage) Remove(hash string) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	delete(s.data, hash)
	return nil
}

// Clean storage
func (s *MemoryStorage) Clean() error {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.data = make(map[string][]byte)
	return nil
}

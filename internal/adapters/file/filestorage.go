package file

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/uuid"

	"github.com/AcroManiac/micropic/internal/domain/interfaces"
)

// Storage implementation
type Storage struct {
	dirName string
}

// NewFileStorage constructor
func NewFileStorage(dirName string) interfaces.Storage {
	dirName += "/" + uuid.New().String()
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		_ = os.MkdirAll(dirName, os.ModeDir|os.ModePerm)
	}
	return &Storage{dirName: dirName}
}

func (s Storage) createPath(hash uint64) string {
	return fmt.Sprintf("%s/%d", s.dirName, hash)
}

// Save data with hash key
func (s *Storage) Save(hash uint64, data []byte) error {
	return ioutil.WriteFile(s.createPath(hash), data, os.ModePerm)
}

// Get stored data for hash key
func (s *Storage) Get(hash uint64) ([]byte, error) {
	return ioutil.ReadFile(s.createPath(hash))
}

// Remove stored data by hash key
func (s *Storage) Remove(hash uint64) error {
	return os.Remove(s.createPath(hash))
}

// Remove all data from storage
func (s *Storage) Clean() error {
	return os.RemoveAll(s.dirName)
}

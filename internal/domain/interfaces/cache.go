package interfaces

import "github.com/AcroManiac/micropic/internal/domain/entities"

// Cache interface
type Cache interface {
	// HasPreview searches preview in cache
	HasPreview(params *entities.PreviewParams) bool

	// Save preview to cache
	Save(preview *entities.Preview) error

	// Get preview from cache
	Get(params *entities.PreviewParams) (*entities.Preview, error)

	// Evict cache item
	Evict() error

	// Clean cache totally
	Clean() error
}

// Storage interface
type Storage interface {
	// Save data with hash key
	Save(hash uint64, data []byte) error

	// Get stored data for hash key
	Get(hash uint64) ([]byte, error)

	// Remove stored data by hash key
	Remove(hash uint64) error

	// Remove all data from storage
	Clean() error
}

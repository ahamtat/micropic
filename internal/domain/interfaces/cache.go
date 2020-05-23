package interfaces

import "github.com/ahamtat/micropic/internal/domain/entities"

// CacheClient interface
type CacheClient interface {
	// Save preview to cache
	Save(preview *entities.Preview) error

	// Get preview from cache
	Get(params *entities.PreviewParams) (*entities.Preview, error)
}

// Cache interface
type Cache interface {
	// CacheClient interface included
	CacheClient

	// HasPreview searches preview in cache
	HasPreview(params *entities.PreviewParams) bool

	// Clean cache totally
	Clean() error
}

// Storage interface
type Storage interface {
	// Save data with hash key
	Save(hash string, data []byte) error

	// Get stored data for hash key
	Get(hash string) ([]byte, error)

	// Remove stored data by hash key
	Remove(hash string) error

	// Remove all data from storage
	Clean() error
}

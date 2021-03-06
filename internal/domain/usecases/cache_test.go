package usecases_test

import (
	"sync"
	"testing"

	"github.com/ahamtat/micropic/internal/domain/usecases"

	"github.com/ahamtat/micropic/internal/domain/interfaces"
	"github.com/pkg/errors"

	"github.com/ahamtat/micropic/internal/domain/entities"
	"github.com/stretchr/testify/require"
)

var testData = []struct {
	preview *entities.Preview
}{
	{
		preview: &entities.Preview{
			Params: &entities.PreviewParams{
				Width:  300,
				Height: 200,
				URL:    "www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk.jpg",
			},
			Image: []byte("Some Base64 code ;)"),
		},
	},
	{
		preview: &entities.Preview{
			Params: &entities.PreviewParams{
				Width:  400,
				Height: 300,
				URL:    "www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk.jpg",
			},
			Image: []byte("Some Base64 code ;)"),
		},
	},
	{
		preview: &entities.Preview{
			Params: &entities.PreviewParams{
				Width:  1024,
				Height: 768,
				URL:    "www.audubon.org/sites/default/files/apa_2016_atlantic-puffin_bonnie_block_kk.jpg",
			},
			Image: []byte("Some Base64 code ;)"),
		},
	},
	{
		preview: &entities.Preview{
			Params: &entities.PreviewParams{
				Width:  1280,
				Height: 1024,
				URL:    "www.audubon.org/sites/default/files/apa-2012-eastern-bluebird_24931-193267_kk_photo-danny-brown.jpg",
			},
			Image: []byte("Some Base64 code ;)"),
		},
	},
}

func createAndFillCache(t *testing.T, size int) interfaces.Cache {
	cache := usecases.NewLRUCache(size, NewMemoryStorage())

	// Save previews to cache
	for i := 0; i < size; i++ {
		err := cache.Save(testData[i].preview)
		require.Nil(t, err, "error should be nil")
	}

	return cache
}

func TestLRUCache_Save(t *testing.T) {
	_ = createAndFillCache(t, 4)
}

func TestLRUCache_Get(t *testing.T) {
	cache := createAndFillCache(t, 3)

	// Get stored preview from cache
	preview, err := cache.Get(testData[0].preview.Params)
	require.Nil(t, err, "error should be nil")
	require.Equal(t, testData[0].preview, preview, "cached preview should be equal to test preview")

	// Get missed preview from cache
	_, err = cache.Get(testData[3].preview.Params)
	require.NotNil(t, err, "method should return error")
}

func TestLRUCache_HasPreview(t *testing.T) {
	size := 4
	cache := createAndFillCache(t, size)

	// Check previews an cache
	for i := 0; i < size; i++ {
		inCache := cache.HasPreview(testData[i].preview.Params)
		require.True(t, inCache, "preview should be in cache")
	}
}

func TestLRUCache_Evict(t *testing.T) {
	cache := createAndFillCache(t, 3)

	// Save one more preview
	err := cache.Save(testData[3].preview)
	require.Nil(t, err, "error should be nil")

	// Check if first element was evicted from cache
	inCache := cache.HasPreview(testData[0].preview.Params)
	require.False(t, inCache, "preview should not be in cache")
}

func TestLRUCache_Clean(t *testing.T) {
	cache := createAndFillCache(t, 4)

	// Clean cache from data
	err := cache.Clean()
	require.Nil(t, err, "error should be nil")
}

// MemoryStorage implementation
type MemoryStorage struct {
	mx   sync.RWMutex
	data map[string][]byte
}

// NewMemoryStorage constructor.
func NewMemoryStorage() interfaces.Storage {
	return &MemoryStorage{
		mx:   sync.RWMutex{},
		data: make(map[string][]byte),
	}
}

// Save data in memory by hash key.
func (s *MemoryStorage) Save(hash string, data []byte) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.data[hash] = data
	return nil
}

// Get stored data by hash key.
func (s *MemoryStorage) Get(hash string) ([]byte, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	data, ok := s.data[hash]
	if !ok {
		return nil, errors.New("error finding value")
	}
	return data, nil
}

// Remove data from storage.
func (s *MemoryStorage) Remove(hash string) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	delete(s.data, hash)
	return nil
}

// Clean storage.
func (s *MemoryStorage) Clean() error {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.data = make(map[string][]byte)
	return nil
}

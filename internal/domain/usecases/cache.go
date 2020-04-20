package usecases

import (
	"fmt"
	"hash/maphash"
	"sync"

	"github.com/pkg/errors"

	"github.com/AcroManiac/micropic/internal/domain/interfaces"

	"github.com/AcroManiac/micropic/internal/domain/entities"
)

// LRUCache implements Least Recently Used cache
type LRUCache struct {
	size        int
	storage     interfaces.Storage
	mx          sync.RWMutex
	hashItemMap map[uint64]*dllItem
	paramsList  doublyLinkedList
}

// NewLRUCache constructor
func NewLRUCache(size int, storage interfaces.Storage) interfaces.Cache {
	return &LRUCache{
		size:        size,
		storage:     storage,
		mx:          sync.RWMutex{},
		hashItemMap: make(map[uint64]*dllItem),
		paramsList:  doublyLinkedList{},
	}
}

// HasPreview searches preview in cache
func (c *LRUCache) HasPreview(params *entities.PreviewParams) bool {
	// Evaluate hash key from preview params
	hash := createHash(params)

	c.mx.RLock()
	defer c.mx.RUnlock()
	_, ok := c.hashItemMap[hash]
	return ok
}

// Save preview to cache
func (c *LRUCache) Save(preview *entities.Preview) error {
	// Push preview params to doubly linked list
	item := c.paramsList.PushHead(preview.Params)
	if c.paramsList.GetLength() > c.size {
		// Remove preview from cache
		_ = c.Evict()
	}

	// Evaluate hash key from preview params
	hash := createHash(preview.Params)

	// Add list item to hash map
	c.mx.Lock()
	c.hashItemMap[hash] = item
	c.mx.Unlock()

	// Store preview image
	return c.storage.Save(hash, preview.Image)
}

// Get preview from cache
func (c *LRUCache) Get(params *entities.PreviewParams) (*entities.Preview, error) {
	// Evaluate hash key from preview params
	hash := createHash(params)

	// Search item in hash map
	c.mx.RLock()
	item, ok := c.hashItemMap[hash]
	if !ok {
		c.mx.RUnlock()
		return nil, errors.New("no preview in cache")
	}
	c.mx.RUnlock()

	// Update item position in list
	c.paramsList.MoveHead(item)

	// Get preview image from storage
	image, err := c.storage.Get(hash)
	if err != nil {
		return nil, errors.Wrap(err, "error loading preview image from storage")
	}

	return &entities.Preview{
		Params: params,
		Image:  image,
	}, nil
}

// Evict cache item
func (c *LRUCache) Evict() error {
	// Get LRU item
	evictedItem := c.paramsList.PopTail()
	if evictedItem == nil {
		return errors.New("failed removing item from doubly linked list tail")
	}
	params, ok := evictedItem.(*entities.PreviewParams)
	if !ok {
		return errors.New("failed type asserting empty interface to preview params")
	}

	// Remove item from hash map
	hash := createHash(params)
	c.mx.Lock()
	delete(c.hashItemMap, hash)
	c.mx.Unlock()

	// Remove preview image from storage
	return c.storage.Remove(hash)
}

// Clean cache totally
func (c *LRUCache) Clean() error {
	c.hashItemMap = make(map[uint64]*dllItem)
	c.paramsList.Clean()
	return c.storage.Clean()
}

// createHash creates hash key from preview params
func createHash(params *entities.PreviewParams) uint64 {
	// The zero Hash value is valid and ready to use; setting an
	// initial seed is not necessary.
	var h maphash.Hash

	// Add a string to the hash, and return the current hash value.
	_, _ = h.WriteString(fmt.Sprintf("%d/%d/%s", params.Width, params.Height, params.URL))
	return h.Sum64()
}

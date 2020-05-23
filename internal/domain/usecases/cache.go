package usecases

import (
	"container/list"
	"crypto/md5" // nolint:gosec
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/ahamtat/micropic/internal/domain/interfaces"

	"github.com/ahamtat/micropic/internal/domain/entities"
)

// LRUCache implements Least Recently Used cache
type LRUCache struct {
	size        int
	storage     interfaces.Storage
	mx          sync.RWMutex
	hashItemMap map[string]*list.Element
	params      *list.List
}

// NewLRUCache constructor.
func NewLRUCache(size int, storage interfaces.Storage) interfaces.Cache {
	return &LRUCache{
		size:        size,
		storage:     storage,
		mx:          sync.RWMutex{},
		hashItemMap: make(map[string]*list.Element),
		params:      list.New(),
	}
}

// HasPreview searches preview in cache.
func (c *LRUCache) HasPreview(params *entities.PreviewParams) bool {
	// Evaluate hash key from preview params
	hash := createHash(params)

	c.mx.RLock()
	defer c.mx.RUnlock()
	_, ok := c.hashItemMap[hash]
	return ok
}

// Save preview to cache.
func (c *LRUCache) Save(preview *entities.Preview) error {
	// Check preview in cache
	if c.HasPreview(preview.Params) {
		// Preview is in cache already. No need to save it
		return nil
	}

	// Push preview params to doubly linked list
	item := c.params.PushFront(preview.Params)
	if c.params.Len() > c.size {
		// Remove preview from cache
		_ = c.evict()
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

// Get preview from cache.
func (c *LRUCache) Get(params *entities.PreviewParams) (*entities.Preview, error) {
	// Evaluate hash key from preview params
	hash := createHash(params)

	// Search item in hash map
	c.mx.RLock()
	item, ok := c.hashItemMap[hash]
	if !ok || item == nil {
		c.mx.RUnlock()
		return nil, errors.New("no preview in cache")
	}
	c.mx.RUnlock()

	// Update item position in list
	c.params.MoveToFront(item)

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

// Evict cache item.
func (c *LRUCache) evict() error {
	// Get LRU item
	item := c.params.Back()
	if item == nil {
		return errors.New("no list back")
	}
	evictedItem := c.params.Remove(item)
	if evictedItem == nil {
		return errors.New("failed removing item from list back")
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

// Clean cache totally.
func (c *LRUCache) Clean() error {
	c.hashItemMap = make(map[string]*list.Element)
	c.params.Init()
	return c.storage.Clean()
}

// createHash creates hash key from preview params.
func createHash(params *entities.PreviewParams) string {
	hash := md5.Sum([]byte(fmt.Sprintf("%d/%d/%s", params.Width, params.Height, params.URL))) // nolint:gosec
	return hex.EncodeToString(hash[:])
}

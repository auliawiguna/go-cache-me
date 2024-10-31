package helpers

import (
	"go-cache-me/models"
	"sync"
	"time"
)

type Cache struct {
	mu    sync.RWMutex
	items map[string]models.CacheItem
}

func NewCache() *Cache {
	c := &Cache{
		items: make(map[string]models.CacheItem),
	}

	go c.CleanupExpired()

	return c
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = models.CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

func (c *Cache) GetAll() map[string]models.CacheItem {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Create a copy of the cache items
	itemsCopy := make(map[string]models.CacheItem, len(c.items))
	for k, v := range c.items {
		itemsCopy[k] = v
	}

	return itemsCopy
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found || item.ExpiresAt.Before(time.Now()) {
		return nil, false
	}

	return item.Value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

func (c *Cache) CleanupExpired() {
	for {
		time.Sleep(time.Minute)
		c.mu.Lock()
		for key, item := range c.items {
			if item.ExpiresAt.Before(time.Now()) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

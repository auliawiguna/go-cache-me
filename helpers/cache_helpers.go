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

	go c.cleanupExpired()

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

func (c *Cache) cleanupExpired() {
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
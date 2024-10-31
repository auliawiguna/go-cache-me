package helpers

import (
	"database/sql"
	"go-cache-me/models"
	"log"
	"sync"
	"time"
)

type Cache struct {
	mu    sync.RWMutex
	items map[string]models.CacheItem
}

var DbInstance *sql.DB
var CacheInstance *Cache

func NewCache() *Cache {
	CacheInstance = &Cache{
		items: make(map[string]models.CacheItem),
	}

	go CacheInstance.CleanupExpired()

	return CacheInstance
}

func InitDb(db *sql.DB) {
	DbInstance = db
}

func (c *Cache) DirectCacheSet(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = models.CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.DirectCacheSet(key, value, ttl)

	DbInstance.Exec("INSERT INTO cache (key, value, expires_at) VALUES (?, ?, ?)", key, value, time.Now().Add(ttl))
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
	if found && item.ExpiresAt.After(time.Now()) {
		return item.Value, true
	}

	// Find from DB
	log.Println("Cache miss for key", key)
	row := DbInstance.QueryRow("SELECT value, expires_at FROM cache WHERE key = ?", key)
	var value string
	var expiresAt time.Time

	log.Println("Querying from DB for", key)
	err := row.Scan(&value, &expiresAt)
	log.Println("Scanning from DB for", key)
	if err == sql.ErrNoRows || expiresAt.Before(time.Now()) {
		return nil, false
	} else if err != nil {
		return nil, false
	}

	c.items[key] = models.CacheItem{
		Value:     value,
		ExpiresAt: expiresAt,
	}

	return value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)

	DbInstance.Exec("DELETE FROM cache WHERE key = ?", key)
}

func (c *Cache) CleanupExpired() {
	for {
		time.Sleep(time.Minute)
		c.mu.Lock()
		for key, item := range c.items {
			if item.ExpiresAt.Before(time.Now()) {
				delete(c.items, key)
				DbInstance.Exec("DELETE FROM cache WHERE key = ?", key)
			}
		}
		c.mu.Unlock()
	}
}

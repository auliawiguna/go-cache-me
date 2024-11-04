package helpers

import (
	"database/sql"
	"go-cache-me/configs"
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
var once sync.Once

func NewCache() *Cache {
	once.Do(func() {
		CacheInstance = &Cache{
			items: make(map[string]models.CacheItem),
		}

		go CleanupExpiredCache()
	})

	return CacheInstance
}

func InitDb(db *sql.DB) {
	DbInstance = db
}

func DirectCacheSet(key string, value interface{}, ttl time.Duration) {
	CacheInstance.mu.Lock()
	defer CacheInstance.mu.Unlock()

	CacheInstance.items[key] = models.CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

func updateCache(db *sql.DB, key string, value interface{}, ttl time.Duration) {
	done := make(chan error, 1)

	go func() {
		tx, err := db.Begin()
		if err != nil {
			done <- err
			return
		}

		_, err = tx.Exec("REPLACE INTO cache (key, value, expires_at) VALUES (?, ?, ?)", key, value, time.Now().Add(ttl))
		if err != nil {
			tx.Rollback()
			done <- err
			return
		}

		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			done <- err
			return
		}

		done <- nil
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Printf("Failed to update cache: %v", err)
		} else {
			log.Println("Cache updated successfully")
		}
	case <-time.After(5 * time.Second): // Adjust timeout as needed
		log.Println("Timeout while updating cache")
	}
}

func SetCookie(key string, value interface{}, ttl time.Duration) {
	DirectCacheSet(key, value, ttl)
	if configs.GetEnv("CACHE_SYSTEM") == "database" {
		updateCache(DbInstance, key, value, ttl)
	}
}

func GetAllCache() map[string]models.CacheItem {
	CacheInstance.mu.RLock()
	defer CacheInstance.mu.RUnlock()

	// Create a copy of the cache items
	itemsCopy := make(map[string]models.CacheItem, len(CacheInstance.items))
	for k, v := range CacheInstance.items {
		itemsCopy[k] = v
	}

	return itemsCopy
}

func GetCache(key string) (interface{}, bool) {

	CacheInstance.mu.RLock()
	defer CacheInstance.mu.RUnlock()

	item, found := CacheInstance.items[key]
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
	if err == sql.ErrNoRows {
		return nil, false
	} else if expiresAt.Before(time.Now()) {
		DbInstance.Exec("DELETE FROM cache WHERE key = ?", key)
		return nil, false
	} else if err != nil {
		return nil, false
	}

	CacheInstance.items[key] = models.CacheItem{
		Value:     value,
		ExpiresAt: expiresAt,
	}

	return value, true
}

func DeleteCache(key string) {
	CacheInstance.mu.Lock()
	defer CacheInstance.mu.Unlock()

	delete(CacheInstance.items, key)

	DbInstance.Exec("DELETE FROM cache WHERE key = ?", key)
}

func CleanupExpiredCache() {
	for {
		time.Sleep(time.Minute)
		CacheInstance.mu.Lock()
		for key, item := range CacheInstance.items {
			if item.ExpiresAt.Before(time.Now()) {
				delete(CacheInstance.items, key)
				DbInstance.Exec("DELETE FROM cache WHERE key = ?", key)
				log.Println("Deleted expired cache with key", key)
			} else {
				log.Println("Updating cache in db with key", key)
				updateCache(DbInstance, key, item.Value, time.Until(item.ExpiresAt))
			}
		}
		CacheInstance.mu.Unlock()
	}
}

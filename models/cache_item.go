package models

import "time"

// CacheItem represents an individual cache entry
type CacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
}

package helpers

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func PreparingDbCache(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS cache (key TEXT PRIMARY KEY, value TEXT, expires_at DATETIME)")

	if err != nil {
		return err
	}

	return nil
}

func SaveCacheToDatabase(db *sql.DB, cache *Cache) error {
	err := PreparingDbCache(db)

	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM cache")
	if err != nil {
		return err
	}

	var items = cache.GetAll()
	// Log the retrieved items
	log.Println("Cache items:", items)

	for key, item := range items {
		_, err = db.Exec("INSERT INTO cache (key, value, expires_at) VALUES (?, ?, ?)", key, item.Value, item.ExpiresAt)
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

func LoadCacheFromDatabase(db *sql.DB, cache *Cache) error {
	err := PreparingDbCache(db)

	if err != nil {
		return err
	}

	rows, err := db.Query("SELECT key, value, expires_at FROM cache")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		var value string
		var expiresAt time.Time

		if err := rows.Scan(&key, &value, &expiresAt); err != nil {
			return err
		}

		cache.Set(key, value, expiresAt.Sub(time.Now()))
	}

	return nil
}

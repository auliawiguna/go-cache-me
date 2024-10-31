package jobs

import (
	"go-cache-me/helpers"
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

func StartCacheCleanupJob(cache *helpers.Cache) {
	log.Println("Scheduling cache cleanup job")
	s := gocron.NewScheduler(time.Local)

	s.Every(60).Seconds().Do(func() {
		log.Println("Running cache cleanup job")
		cache.CleanupExpired()
	})

	s.StartAsync()
}

package main

import (
	"database/sql"
	"go-cache-me/configs"
	"go-cache-me/helpers"
	"go-cache-me/jobs"
	"go-cache-me/middlewares"
	"go-cache-me/routes"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"

	_ "go-cache-me/docs" // replace with the actual path to your docs folder

	swagger "github.com/gofiber/swagger"
)

func main() {
	// Open SQLite
	db, err := sql.Open("sqlite3", "./db/cache.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := fiber.New()
	cache := helpers.NewCache()

	routes.RegisterCacheRoutes(app, cache)

	middlewares.DefaultMiddleware(app)

	jobs.StartCacheCleanupJob(cache)

	// Load existing cache from database
	if err := helpers.LoadCacheFromDatabase(db, cache); err != nil {
		log.Fatal("Error loading cache from DB", err)
	}

	// Save cache to database on shutdown
	go func() {
		// Retrieve all cache items
		items := cache.GetAll()

		// Log the retrieved items
		log.Println("Cache items before saving to DB:", items)

		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c

		log.Println("Saving cache to database")
		if err := helpers.SaveCacheToDatabase(db, cache); err != nil {
			log.Fatal("Error saving cache to DB", err)
		}

		log.Println("Shutting down")
		os.Exit(0)
	}()

	if configs.GetEnv("ENV") == "dev" {
		app.Get("/api/v0/docs/*", swagger.HandlerDefault) // default
	}

	if configs.GetEnv("ENV") == "dev" {
		log.Println("Starting server on development mode")
		helpers.StartServer(app)
	} else {
		log.Println("Starting server on production mode")
		helpers.StartServerWithGracefulShutdown(app)

	}
}

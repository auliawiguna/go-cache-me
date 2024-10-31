package main

import (
	"database/sql"
	"go-cache-me/configs"
	"go-cache-me/helpers"
	"go-cache-me/jobs"
	"go-cache-me/middlewares"
	"go-cache-me/routes"
	"log"

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
	helpers.InitDb(db)

	middlewares.DefaultMiddleware(app)

	routes.RegisterCacheRoutes(app, cache)

	jobs.StartCacheCleanupJob(cache)

	if configs.GetEnv("ENV") == "dev" {
		app.Get("/api/v0/docs/*", swagger.HandlerDefault) // default
	}

	// Load existing cache from database
	if err := helpers.LoadCacheFromDatabase(db, cache); err != nil {
		log.Fatal("Error loading cache from DB", err)
	}

	if configs.GetEnv("ENV") == "dev" {
		log.Println("Starting server on development mode")
		helpers.StartServer(app)
	} else {
		log.Println("Starting server on production mode")
		helpers.StartServerWithGracefulShutdown(app)

	}
}

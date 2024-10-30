package main

import (
	"go-cache-me/configs"
	"go-cache-me/helpers"
	"go-cache-me/middlewares"
	"go-cache-me/routes"

	"github.com/gofiber/fiber/v2"

	_ "go-cache-me/docs" // replace with the actual path to your docs folder

	swagger "github.com/gofiber/swagger"
)

func main() {
	app := fiber.New()
	cache := helpers.NewCache()

	routes.RegisterCacheRoutes(app, cache)

	middlewares.DefaultMiddleware(app)

	if configs.GetEnv("ENV") == "dev" {
		helpers.StartServer(app)
		app.Get("/swagger/*", swagger.HandlerDefault) // default
	} else {
		helpers.StartServerWithGracefulShutdown(app)

	}
}

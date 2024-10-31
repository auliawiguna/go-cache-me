package routes

import (
	"go-cache-me/controllers"
	"go-cache-me/helpers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func RegisterCacheRoutes(app *fiber.App, cache *helpers.Cache) {
	cacheRoute := app.Group("/api/caches", logger.New())

	cacheRoute.Post("/", controllers.SetCache)
	cacheRoute.Get("/", controllers.GetAllCache)
	cacheRoute.Get("/key/:key", controllers.GetCache)
	cacheRoute.Delete("/:key", controllers.DeleteCache)
	cacheRoute.Post("/get-or-set", controllers.GetOrSetCache)
}

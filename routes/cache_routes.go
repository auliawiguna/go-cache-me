package routes

import (
	"go-cache-me/helpers"
	"time"

	"github.com/goccy/go-json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func RegisterCacheRoutes(app *fiber.App, cache *helpers.Cache) {
	cacheRoute := app.Group("/cache", logger.New())

	cacheRoute.Post("/", func(c *fiber.Ctx) error {
		key := c.Query("key")
		ttl, err := time.ParseDuration(c.Query("ttl"))

		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid TTL")
		}

		var value interface{}

		if err := json.Unmarshal(c.Body(), &value); err != nil {
			// If not json, back to string
			value = string(c.Body())
		}

		cache.Set(key, value, ttl)

		return c.SendStatus(fiber.StatusCreated)
	})

	cacheRoute.Get("/:key", func(c *fiber.Ctx) error {
		key := c.Params("key")

		value, found := cache.Get(key)

		if !found {
			return fiber.NewError(fiber.StatusNotFound, "Key not found")
		}

		return c.JSON(value)
	})

	cacheRoute.Post("/get-or-set", func(c *fiber.Ctx) error {
		key := c.Query("key")
		ttl, err := time.ParseDuration(c.Query("ttl"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid TTL")
		}

		var value interface{}
		if err := json.Unmarshal(c.Body(), &value); err != nil {
			// If not json, back to string
			value = string(c.Body())
		}

		cache.Set(key, value, ttl)

		return c.JSON(value)
	})

	cacheRoute.Delete("/:key", func(c *fiber.Ctx) error {
		key := c.Params("key")

		cache.Delete(key)

		return c.SendStatus(fiber.StatusNoContent)
	})
}

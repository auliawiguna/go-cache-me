package controllers

import (
	"log"
	"time"

	"go-cache-me/helpers"

	"github.com/gofiber/fiber/v2"
)

type CacheRequest struct {
	Key   string      `json:"key"`
	TTL   string      `json:"ttl"`
	Value interface{} `json:"value"`
}

var cache = helpers.NewCache()

// @Summary Set a key-value pair
// @Description Set a key-value pair in the cache
// @Tags cache
// @Accept application/json
// @Produce json
// @Param cacheRequest body controllers.CacheRequest true "Cache Request"
// @Success 201 "Created"
// @Failure 400 "Invalid TTL"
// @Router /cache [post]
func SetCache(c *fiber.Ctx) error {
	var req CacheRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	ttl, err := time.ParseDuration(req.TTL)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid TTL")
	}

	cache.Set(req.Key, req.Value, ttl)

	return c.SendStatus(fiber.StatusCreated)
}

// @Summary Get a key-value pair
// @Description Get a key-value pair from the cache
// @Tags cache
// @Accept json
// @Produce json
// @Param key path string true "Key"
// @Success 200 "OK"
// @Failure 404 "Key not found"
// @Router /cache/{key} [get]
func GetCache(c *fiber.Ctx) error {
	key := c.Params("key")

	value, found := cache.Get(key)

	if !found {
		return fiber.NewError(fiber.StatusNotFound, "Key not found")
	}

	return c.JSON(value)
}

// @Summary Delete a key-value pair
// @Description Delete a key-value pair from the cache
// @Tags cache
// @Accept json
// @Produce json
// @Param key path string true "Key"
// @Success 204 "No Content"
// @Router /cache/{key} [delete]
func DeleteCache(c *fiber.Ctx) error {
	key := c.Params("key")

	cache.Delete(key)

	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary Get or set a key-value pair
// @Description Get or set a key-value pair in the cache
// @Tags cache
// @Accept json
// @Produce json
// @Param cacheRequest body controllers.CacheRequest true "Cache Request"
// @Success 200 "OK"
// @Failure 400 "Invalid TTL"
// @Router /cache/get-or-set [post]
func GetOrSetCache(c *fiber.Ctx) error {
	var req CacheRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	ttl, err := time.ParseDuration(req.TTL)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid TTL")
	}

	value, found := cache.Get(req.Key)
	log.Printf("Value: %v, Found: %v", value, found)

	if !found {
		cache.Set(req.Key, req.Value, ttl)
		return c.JSON(req.Value)
	}

	return c.JSON(value)
}

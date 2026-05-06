package handler

import (
	"github.com/gofiber/fiber/v3"
)

// HealthHandler health check do servidor
func HealthHandler(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": "aulaflash",
	})
}



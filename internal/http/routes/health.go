package routes

import (
	"time"

	"backend-boiler/internal/shared/response"
	"github.com/gofiber/fiber/v3"
)

func registerHealthRoutes(router fiber.Router) {
	router.Get("/health", func(c fiber.Ctx) error {
		return response.OK(c, fiber.Map{
			"status":    "ok",
			"timestamp": time.Now().UTC(),
		})
	})

	router.Get("/ready", func(c fiber.Ctx) error {
		return response.OK(c, fiber.Map{
			"status": "ready",
		})
	})
}

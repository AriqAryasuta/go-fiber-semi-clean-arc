package routes

import (
	"backend-boiler/internal/shared/response"
	"github.com/gofiber/fiber/v3"
)

func registerSwaggerRoutes(router fiber.Router) {
	router.Get("/swagger/*", func(c fiber.Ctx) error {
		return response.OK(c, fiber.Map{
			"message": "swagger is available in docs/swagger. run `make swagger` to regenerate.",
		})
	})
}

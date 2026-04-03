package user

import (
	"backend-boiler/internal/shared/validator"
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(router fiber.Router) {
	// Placeholder wiring. In production this should receive dependencies
	// from app container to avoid opening DB here.
	service := NewService(NewRepository(nil))
	controller := NewController(service, validator.New())

	group := router.Group("/users")
	group.Get("/", controller.List)
	group.Post("/", controller.Create)
}

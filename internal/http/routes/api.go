package routes

import (
	"backend-boiler/internal/modules/user"
	"github.com/gofiber/fiber/v3"
)

func Register(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	registerHealthRoutes(v1)
	registerSwaggerRoutes(app)
	user.RegisterRoutes(v1)
}

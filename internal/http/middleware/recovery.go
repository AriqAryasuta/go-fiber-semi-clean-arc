package middleware

import "github.com/gofiber/fiber/v3"

func Recovery() fiber.Handler {
	return func(c fiber.Ctx) error {
		defer func() {
			if recover() != nil {
				_ = c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"success": false,
					"message": "internal server error",
				})
			}
		}()
		return c.Next()
	}
}

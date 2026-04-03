package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
)

func Logging(logger *slog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		logger.Info("http request",
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"duration_ms", time.Since(start).Milliseconds(),
			"request_id", c.Locals("request_id"),
		)

		return err
	}
}

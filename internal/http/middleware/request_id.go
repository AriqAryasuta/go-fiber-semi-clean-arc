package middleware

import (
	"backend-boiler/internal/shared/utils"
	"github.com/gofiber/fiber/v3"
)

const requestIDHeader = "X-Request-ID"

func RequestID() fiber.Handler {
	return func(c fiber.Ctx) error {
		requestID := c.Get(requestIDHeader)
		if requestID == "" {
			requestID = utils.NewID()
		}
		c.Set(requestIDHeader, requestID)
		c.Locals("request_id", requestID)
		return c.Next()
	}
}

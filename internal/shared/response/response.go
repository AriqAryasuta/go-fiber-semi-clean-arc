package response

import "github.com/gofiber/fiber/v3"

type Payload struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func OK(c fiber.Ctx, data any) error {
	return c.Status(fiber.StatusOK).JSON(Payload{
		Success: true,
		Data:    data,
	})
}

func Created(c fiber.Ctx, data any) error {
	return c.Status(fiber.StatusCreated).JSON(Payload{
		Success: true,
		Data:    data,
	})
}

func Fail(c fiber.Ctx, status int, message string, err any) error {
	return c.Status(status).JSON(Payload{
		Success: false,
		Message: message,
		Error:   err,
	})
}

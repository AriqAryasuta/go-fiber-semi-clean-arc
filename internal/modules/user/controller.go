package user

import (
	"backend-boiler/internal/shared/response"
	sharedValidator "backend-boiler/internal/shared/validator"
	"github.com/gofiber/fiber/v3"
)

type Controller struct {
	service   *Service
	validator *sharedValidator.Validator
}

func NewController(service *Service, validator *sharedValidator.Validator) *Controller {
	return &Controller{
		service:   service,
		validator: validator,
	}
}

func (c *Controller) Create(ctx fiber.Ctx) error {
	payload := CreateUserRequest{}
	if err := ctx.Bind().Body(&payload); err != nil {
		return response.Fail(ctx, fiber.StatusBadRequest, "invalid request body", err.Error())
	}
	if err := c.validator.ValidateStruct(payload); err != nil {
		return response.Fail(ctx, fiber.StatusBadRequest, "validation failed", err.Error())
	}
	result, err := c.service.Create(payload)
	if err != nil {
		return response.Fail(ctx, fiber.StatusInternalServerError, "failed to create user", err.Error())
	}
	return response.Created(ctx, result)
}

func (c *Controller) List(ctx fiber.Ctx) error {
	result, err := c.service.List()
	if err != nil {
		return response.Fail(ctx, fiber.StatusInternalServerError, "failed to fetch users", err.Error())
	}
	return response.OK(ctx, result)
}

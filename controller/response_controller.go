package controller

import (
	"net/http"
	"solution-challange/response"

	"github.com/gofiber/fiber/v2"
)

func APIResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(http.StatusOK).JSON(response.UserResponse{
		Status:  status,
		Message: message,
		Data:    &fiber.Map{"data": data},
	})
}

package helper

import (
	"github.com/gofiber/fiber/v2"
)

type (
	response struct {
		Status  bool        `json:"status"`
		Message string      `json:"message"`
		Errors  []string    `json:"errors"`
		Data    interface{} `json:"data"`
	}
)

func ResponseSuccessJson(c *fiber.Ctx, data interface{}) error {
	message := "Succeed to " + string(c.Request().Header.Method()) + " data"
	res := response{
		Status:  true,
		Message: message,
		Errors:  nil,
		Data:    data,
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func ResponseErrorJson(c *fiber.Ctx, code int, err error) error {
	message := "Failed to " + string(c.Request().Header.Method()) + " data"
	res := response{
		Status:  false,
		Message: message,
		Errors:  []string{err.Error()},
		Data:    nil,
	}
	return c.Status(code).JSON(res)
}

package helper

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func JwtMiddlewareErrorHandler(c *fiber.Ctx, err error) error {
	return ResponseErrorJson(c, fiber.StatusUnauthorized, err)
}

func CheckAdminTokenHandler(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	isAdmin := claims["isAdmin"].(bool)
	if isAdmin {
		return c.Next()
	} else {
		return ResponseErrorJson(c, fiber.StatusUnauthorized, errors.New("unauthorized"))
	}
}

func GetUserIdFromToken(c *fiber.Ctx) (int, error) {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	idString := claims["idString"].(string)
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		return -1, errors.New("invalid or malformed JWT")
	}
	return idInt, nil
}

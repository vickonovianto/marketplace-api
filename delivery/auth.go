package delivery

import (
	"marketplace-api/helper"
	"marketplace-api/model"

	"github.com/gofiber/fiber/v2"
)

type authDelivery struct {
	userUsecase model.UserUsecase
}

type AuthDelivery interface {
	MountUnprotectedRoutes(group fiber.Router)
}

func NewAuthDelivery(userUsecase model.UserUsecase) AuthDelivery {
	return &authDelivery{userUsecase: userUsecase}
}

func (a *authDelivery) MountUnprotectedRoutes(group fiber.Router) {
	group.Post("/register", a.RegisterUserHandler)
	group.Post("/login", a.LoginUserHandler)
}

func (a *authDelivery) RegisterUserHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	var req model.UserRegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	req.Trim()
	if err := req.Validate(); err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	userRegisterResponse, err := a.userUsecase.RegisterUser(ctx, &req)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	return helper.ResponseSuccessJson(c, userRegisterResponse)
}

func (a *authDelivery) LoginUserHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	var req model.UserLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	req.Trim()
	if err := req.Validate(); err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	userLoginResponse, err := a.userUsecase.LoginUser(ctx, &req)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusUnauthorized, err)
	}
	return helper.ResponseSuccessJson(c, userLoginResponse)
}

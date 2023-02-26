package delivery

import (
	"marketplace-api/helper"
	"marketplace-api/model"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type userDelivery struct {
	userUsecase model.UserUsecase
}

type UserDelivery interface {
	MountProtectedRoutes(jwtMiddleware func(*fiber.Ctx) error, group fiber.Router)
}

func NewUserDelivery(userUsecase model.UserUsecase) UserDelivery {
	return &userDelivery{userUsecase: userUsecase}
}

func (p *userDelivery) MountProtectedRoutes(jwtMiddleware func(*fiber.Ctx) error, group fiber.Router) {
	group.Get("", jwtMiddleware, p.GetCurrentUserHandler)
	group.Put("", jwtMiddleware, p.EditCurrentUserHandler)
}

func (p *userDelivery) GetCurrentUserHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	userId, err := helper.GetUserIdFromToken(c)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	userResponse, err := p.userUsecase.GetCurrentUser(ctx, userId)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	return helper.ResponseSuccessJson(c, userResponse)
}

func (p *userDelivery) EditCurrentUserHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	var req model.UserUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	req.Trim()
	if err := req.Validate(); err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	userId, err := helper.GetUserIdFromToken(c)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	userResponse, err := p.userUsecase.EditCurrentUser(ctx, userId, &req)
	if err != nil {
		return helper.ResponseErrorJson(c, http.StatusBadRequest, err)
	}
	return helper.ResponseSuccessJson(c, userResponse)
}

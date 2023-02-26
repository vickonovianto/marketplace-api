package delivery

import (
	"errors"
	"marketplace-api/helper"
	"marketplace-api/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type alamatDelivery struct {
	alamatUsecase model.AlamatUsecase
}

type AlamatDelivery interface {
	MountProtectedRoutes(jwtMiddleware func(*fiber.Ctx) error, group fiber.Router)
}

func NewAlamatDelivery(alamatUsecase model.AlamatUsecase) AlamatDelivery {
	return &alamatDelivery{alamatUsecase: alamatUsecase}
}

func (p *alamatDelivery) MountProtectedRoutes(jwtMiddleware func(*fiber.Ctx) error, group fiber.Router) {
	group.Post("", jwtMiddleware, p.StoreAlamatHandler)
	group.Get("", jwtMiddleware, p.FetchAndFilterAlamatHandler)
	group.Get("/:id", jwtMiddleware, p.DetailAlamatHandler)
	group.Put("/:id", jwtMiddleware, p.EditAlamatHandler)
	group.Delete("/:id", jwtMiddleware, p.DeleteAlamatHandler)
}

func (p *alamatDelivery) StoreAlamatHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	var req model.AlamatRequest
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
	req.IdUser = userId
	alamatResponse, err := p.alamatUsecase.StoreAlamat(ctx, &req)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusInternalServerError, err)
	}
	return helper.ResponseSuccessJson(c, alamatResponse)
}

func (p *alamatDelivery) FetchAndFilterAlamatHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	judulAlamat := strings.TrimSpace(c.Query("judul_alamat"))
	if len(judulAlamat) > 255 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("judul alamat cannot exceed 255 characters"))
	}
	userId, err := helper.GetUserIdFromToken(c)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	alamatResponses, err := p.alamatUsecase.FetchAndFilterAlamat(ctx, userId, judulAlamat)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusInternalServerError, err)
	}
	return helper.ResponseSuccessJson(c, alamatResponses)
}

func (p *alamatDelivery) DetailAlamatHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	idString := c.Params("id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid id"))
	}
	userId, err := helper.GetUserIdFromToken(c)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	alamatResponse, err := p.alamatUsecase.GetAlamatByID(ctx, idInt, userId)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	return helper.ResponseSuccessJson(c, alamatResponse)
}

func (p *alamatDelivery) EditAlamatHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	var req model.AlamatRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	req.Trim()
	if err := req.Validate(); err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	idString := c.Params("id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid id"))
	}
	userId, err := helper.GetUserIdFromToken(c)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	req.IdUser = userId
	alamatResponse, err := p.alamatUsecase.EditAlamatByID(ctx, idInt, &req)
	if err != nil {
		return helper.ResponseErrorJson(c, http.StatusBadRequest, err)
	}
	return helper.ResponseSuccessJson(c, alamatResponse)
}

func (p *alamatDelivery) DeleteAlamatHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	idString := c.Params("id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid id"))
	}
	userId, err := helper.GetUserIdFromToken(c)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	err = p.alamatUsecase.DestroyAlamat(ctx, idInt, userId)
	if err != nil {
		return helper.ResponseErrorJson(c, http.StatusBadRequest, err)
	}
	return helper.ResponseSuccessJson(c, "")
}

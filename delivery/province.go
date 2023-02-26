package delivery

import (
	"errors"
	"marketplace-api/helper"
	"marketplace-api/model"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type provinceDelivery struct {
	provinceUsecase model.ProvinceUsecase
}

type ProvinceDelivery interface {
	MountUnprotectedRoutes(group fiber.Router)
}

func NewProvinceDelivery(provinceUsecase model.ProvinceUsecase) ProvinceDelivery {
	return &provinceDelivery{provinceUsecase: provinceUsecase}
}

func (p *provinceDelivery) MountUnprotectedRoutes(group fiber.Router) {
	group.Get("/listprovincies", p.FetchProvinceHandler)
	group.Get("/detailprovince/:prov_id", p.DetailProvinceHandler)
}

func (p *provinceDelivery) FetchProvinceHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	provinces, err := p.provinceUsecase.FetchAllProvince(ctx)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusInternalServerError, err)
	}
	return helper.ResponseSuccessJson(c, provinces)
}

func (p *provinceDelivery) DetailProvinceHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	provIdString := c.Params("prov_id")
	if len(provIdString) != 2 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid prov_id"))
	}
	_, err := strconv.Atoi(provIdString)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid prov_id"))
	}
	province, err := p.provinceUsecase.GetProvinceByID(ctx, provIdString)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	return helper.ResponseSuccessJson(c, province)
}

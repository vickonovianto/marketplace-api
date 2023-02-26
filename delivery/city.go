package delivery

import (
	"errors"
	"marketplace-api/helper"
	"marketplace-api/model"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type cityDelivery struct {
	cityUsecase model.CityUsecase
}

type CityDelivery interface {
	MountUnprotectedRoutes(group fiber.Router)
}

func NewCityDelivery(cityUsecase model.CityUsecase) CityDelivery {
	return &cityDelivery{cityUsecase: cityUsecase}
}

func (p *cityDelivery) MountUnprotectedRoutes(group fiber.Router) {
	group.Get("/listcities/:prov_id", p.FetchCityHandler)
	group.Get("/detailcity/:city_id", p.DetailCityHandler)
}

func (p *cityDelivery) FetchCityHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	provIdString := c.Params("prov_id")
	if len(provIdString) != 2 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid prov_id"))
	}
	_, err := strconv.Atoi(provIdString)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid prov_id"))
	}
	cities, err := p.cityUsecase.FetchAllCity(ctx, provIdString)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	return helper.ResponseSuccessJson(c, cities)
}

func (p *cityDelivery) DetailCityHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	cityIdString := c.Params("city_id")
	if len(cityIdString) != 4 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid city_id"))
	}
	_, err := strconv.Atoi(cityIdString)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid city_id"))
	}
	city, err := p.cityUsecase.GetCityByID(ctx, cityIdString)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	return helper.ResponseSuccessJson(c, city)
}

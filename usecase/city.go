package usecase

import (
	"context"
	"marketplace-api/model"
)

type cityUsecase struct {
	cityRepository model.CityRepository
}

func NewCityUsecase(cityRepository model.CityRepository) model.CityUsecase {
	return &cityUsecase{cityRepository: cityRepository}
}

func (c *cityUsecase) FetchAllCity(ctx context.Context, provinceId string) ([]*model.City, error) {
	cities, err := c.cityRepository.FetchAll(ctx, provinceId)
	if err != nil {
		return nil, err
	}
	return cities, nil
}

func (c *cityUsecase) GetCityByID(ctx context.Context, cityId string) (*model.City, error) {
	provinceId := cityId[0:2]
	city, err := c.cityRepository.FindByID(ctx, provinceId, cityId)
	if err != nil {
		return nil, err
	}
	return city, nil
}

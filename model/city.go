package model

import (
	"context"
)

type (
	City struct {
		ID         string `json:"id"`
		ProvinceID string `json:"province_id"`
		Name       string `json:"name"`
	}

	CityRepository interface {
		FetchAll(ctx context.Context, provinceId string) ([]*City, error)
		FindByID(ctx context.Context, provinceId string, cityId string) (*City, error)
	}

	CityUsecase interface {
		FetchAllCity(ctx context.Context, provinceId string) ([]*City, error)
		GetCityByID(ctx context.Context, cityId string) (*City, error)
	}
)

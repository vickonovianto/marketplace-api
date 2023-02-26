package model

import (
	"context"
)

type (
	Province struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	ProvinceRepository interface {
		FetchAll(ctx context.Context) ([]*Province, error)
		FindByID(ctx context.Context, id string) (*Province, error)
	}

	ProvinceUsecase interface {
		FetchAllProvince(ctx context.Context) ([]*Province, error)
		GetProvinceByID(ctx context.Context, id string) (*Province, error)
	}
)

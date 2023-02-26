package usecase

import (
	"context"
	"marketplace-api/model"
)

type provinceUsecase struct {
	provinceRepository model.ProvinceRepository
}

func NewProvinceUsecase(provinceRepository model.ProvinceRepository) model.ProvinceUsecase {
	return &provinceUsecase{provinceRepository: provinceRepository}
}

func (p *provinceUsecase) FetchAllProvince(ctx context.Context) ([]*model.Province, error) {
	provinces, err := p.provinceRepository.FetchAll(ctx)
	if err != nil {
		return nil, err
	}
	return provinces, nil
}

func (p *provinceUsecase) GetProvinceByID(ctx context.Context, id string) (*model.Province, error) {
	province, err := p.provinceRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return province, nil
}

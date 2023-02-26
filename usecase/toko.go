package usecase

import (
	"context"
	"errors"
	"marketplace-api/model"

	"github.com/jinzhu/copier"
)

type tokoUsecase struct {
	tokoRepository model.TokoRepository
}

func NewTokoUsecase(tokoRepository model.TokoRepository) model.TokoUsecase {
	return &tokoUsecase{tokoRepository: tokoRepository}
}

func (t *tokoUsecase) FetchAndPaginateToko(ctx context.Context, req *model.TokoFetchPaginateRequest) (*model.TokoFetchPaginateResponse, error) {
	listToko, err := t.tokoRepository.FetchAndPaginate(ctx, req)
	if err != nil {
		return nil, err
	}
	tokoFetchPaginateResponse := new(model.TokoFetchPaginateResponse)
	copier.Copy(tokoFetchPaginateResponse, req)
	data := []*model.TokoGetByIDResponse{}
	copier.Copy(&data, listToko)
	tokoFetchPaginateResponse.Data = data
	return tokoFetchPaginateResponse, nil
}

func (t *tokoUsecase) GetTokoByID(ctx context.Context, tokoId int) (*model.TokoGetByIDResponse, error) {
	toko, err := t.tokoRepository.FindByTokoID(ctx, tokoId)
	if err != nil {
		return nil, err
	}
	tokoGetByIDResponse := new(model.TokoGetByIDResponse)
	copier.Copy(tokoGetByIDResponse, toko)
	return tokoGetByIDResponse, nil
}

func (t *tokoUsecase) GetMyToko(ctx context.Context, userId int) (*model.GetMyTokoResponse, error) {
	toko, err := t.tokoRepository.FindByUserID(ctx, userId)
	if err != nil {
		return nil, err
	}
	getMyTokoResponse := new(model.GetMyTokoResponse)
	copier.Copy(getMyTokoResponse, toko)
	return getMyTokoResponse, nil
}

func (t *tokoUsecase) EditToko(ctx context.Context, req *model.TokoUpdateRequest) (*model.TokoUpdateResponse, error) {
	myToko, err := t.tokoRepository.FindByUserID(ctx, req.IdUser)
	if err != nil {
		return nil, err
	}
	if myToko.ID != req.ID {
		return nil, errors.New("unauthorized")
	}
	toko := new(model.Toko)
	copier.Copy(toko, req)
	toko, err = t.tokoRepository.UpdateByTokoID(ctx, toko.ID, toko)
	if err != nil {
		return nil, err
	}
	tokoUpdateResponse := new(model.TokoUpdateResponse)
	copier.Copy(tokoUpdateResponse, toko)
	return tokoUpdateResponse, nil
}

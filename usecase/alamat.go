package usecase

import (
	"context"
	"errors"
	"marketplace-api/model"

	"github.com/jinzhu/copier"
)

type alamatUsecase struct {
	alamatRepository model.AlamatRepository
}

func NewAlamatUsecase(alamatRepository model.AlamatRepository) model.AlamatUsecase {
	return &alamatUsecase{alamatRepository: alamatRepository}
}

func (a *alamatUsecase) StoreAlamat(ctx context.Context, req *model.AlamatRequest) (*model.AlamatResponse, error) {
	alamat := new(model.Alamat)
	copier.Copy(alamat, req)
	alamat, err := a.alamatRepository.Create(ctx, alamat)
	if err != nil {
		return nil, err
	}
	alamatResponse := new(model.AlamatResponse)
	copier.Copy(alamatResponse, alamat)
	return alamatResponse, nil
}

func (a *alamatUsecase) FetchAndFilterAlamat(ctx context.Context, userId int, judulAlamat string) ([]*model.AlamatResponse, error) {
	alamatList, err := a.alamatRepository.FetchAndFilter(ctx, userId, judulAlamat)
	if err != nil {
		return nil, err
	}
	alamatResponses := []*model.AlamatResponse{}
	copier.Copy(&alamatResponses, &alamatList)
	return alamatResponses, nil
}

func (a *alamatUsecase) GetAlamatByID(ctx context.Context, alamatId int, userId int) (*model.AlamatResponse, error) {
	alamat, err := a.alamatRepository.FindByID(ctx, alamatId)
	if err != nil {
		return nil, err
	}
	if alamat.IdUser != userId {
		return nil, errors.New("unauthorized")
	}
	alamatResponse := new(model.AlamatResponse)
	copier.Copy(alamatResponse, alamat)
	return alamatResponse, nil
}

func (a *alamatUsecase) EditAlamatByID(ctx context.Context, alamatId int, req *model.AlamatRequest) (*model.AlamatResponse, error) {
	alamat, err := a.alamatRepository.FindByID(ctx, alamatId)
	if err != nil {
		return nil, err
	}
	if alamat.IdUser != req.IdUser {
		return nil, errors.New("unauthorized")
	}
	copier.Copy(alamat, req)
	alamat, err = a.alamatRepository.UpdateByID(ctx, alamatId, alamat)
	if err != nil {
		return nil, err
	}
	alamatResponse := new(model.AlamatResponse)
	copier.Copy(alamatResponse, alamat)
	return alamatResponse, nil
}

func (a *alamatUsecase) DestroyAlamat(ctx context.Context, alamatId int, userId int) error {
	alamat, err := a.alamatRepository.FindByID(ctx, alamatId)
	if err != nil {
		return err
	}
	if alamat.IdUser != userId {
		return errors.New("unauthorized")
	}
	err = a.alamatRepository.Delete(ctx, alamatId)
	if err != nil {
		return err
	}
	return nil
}

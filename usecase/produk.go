package usecase

import (
	"context"
	"errors"
	"marketplace-api/model"

	"github.com/gosimple/slug"
	"github.com/jinzhu/copier"
)

type produkUsecase struct {
	produkRepository     model.ProdukRepository
	fotoProdukRepository model.FotoProdukRepository
	tokoRepository       model.TokoRepository
	categoryRepository   model.CategoryRepository
}

func NewProdukUsecase(
	produkRepository model.ProdukRepository,
	fotoProdukRepository model.FotoProdukRepository,
	tokoRepository model.TokoRepository,
	categoryRepository model.CategoryRepository,
) model.ProdukUsecase {
	return &produkUsecase{
		produkRepository:     produkRepository,
		fotoProdukRepository: fotoProdukRepository,
		tokoRepository:       tokoRepository,
		categoryRepository:   categoryRepository,
	}
}

func (p *produkUsecase) StoreProduk(ctx context.Context, req *model.ProdukRequest, userId int) (*model.ProdukResponse, error) {
	produkSlug := slug.Make(req.NamaProduk)
	req.Slug = produkSlug

	toko, err := p.tokoRepository.FindByUserID(ctx, userId)
	if err != nil {
		return nil, err
	}
	req.IdToko = toko.ID

	category, err := p.categoryRepository.FindByID(ctx, req.IdCategory)
	if err != nil {
		return nil, err
	}

	produk := new(model.Produk)
	copier.Copy(produk, req)
	produk, fotoProdukList, err := p.produkRepository.CreateProdukAndFotoProduk(ctx, produk, req.PhotoUrls)
	if err != nil {
		return nil, err
	}
	produkResponse := new(model.ProdukResponse)
	copier.Copy(produkResponse, produk)

	fotoProdukResponses := []*model.FotoProdukResponse{}
	copier.Copy(&fotoProdukResponses, &fotoProdukList)
	produkResponse.Photos = fotoProdukResponses

	tokoGetByIDResponse := new(model.TokoGetByIDResponse)
	copier.Copy(tokoGetByIDResponse, toko)
	produkResponse.Toko = tokoGetByIDResponse

	categoryResponse := new(model.CategoryResponse)
	copier.Copy(categoryResponse, category)
	produkResponse.Category = categoryResponse

	return produkResponse, nil
}

func (p *produkUsecase) FetchProduk(ctx context.Context, req *model.ProdukFetchRequest) ([]*model.ProdukResponse, error) {
	produkList, err := p.produkRepository.Fetch(ctx, req)
	if err != nil {
		return nil, err
	}
	produkResponses := []*model.ProdukResponse{}

	for _, produk := range produkList {
		produkResponse := new(model.ProdukResponse)
		copier.Copy(produkResponse, produk)

		toko, err := p.tokoRepository.FindByTokoID(ctx, produk.IdToko)
		if err != nil {
			return nil, err
		}
		tokoGetByIDResponse := new(model.TokoGetByIDResponse)
		copier.Copy(tokoGetByIDResponse, toko)
		produkResponse.Toko = tokoGetByIDResponse

		category, err := p.categoryRepository.FindByID(ctx, produk.IdCategory)
		if err != nil {
			return nil, err
		}
		categoryResponse := new(model.CategoryResponse)
		copier.Copy(categoryResponse, category)
		produkResponse.Category = categoryResponse

		fotoProdukResponses := []*model.FotoProdukResponse{}
		fotoProdukList, err := p.fotoProdukRepository.FetchByProdukId(ctx, produk.ID)
		if err != nil {
			return nil, err
		}
		copier.Copy(&fotoProdukResponses, &fotoProdukList)
		produkResponse.Photos = fotoProdukResponses

		produkResponses = append(produkResponses, produkResponse)
	}

	return produkResponses, nil
}

func (p *produkUsecase) GetProdukByID(ctx context.Context, produkId int) (*model.ProdukResponse, error) {
	produk, err := p.produkRepository.FindByID(ctx, produkId)
	if err != nil {
		return nil, err
	}
	produkResponse := new(model.ProdukResponse)
	copier.Copy(produkResponse, produk)

	toko, err := p.tokoRepository.FindByTokoID(ctx, produk.IdToko)
	if err != nil {
		return nil, err
	}
	tokoGetByIDResponse := new(model.TokoGetByIDResponse)
	copier.Copy(tokoGetByIDResponse, toko)
	produkResponse.Toko = tokoGetByIDResponse

	category, err := p.categoryRepository.FindByID(ctx, produk.IdCategory)
	if err != nil {
		return nil, err
	}
	categoryResponse := new(model.CategoryResponse)
	copier.Copy(categoryResponse, category)
	produkResponse.Category = categoryResponse

	fotoProdukResponses := []*model.FotoProdukResponse{}
	fotoProdukList, err := p.fotoProdukRepository.FetchByProdukId(ctx, produk.ID)
	if err != nil {
		return nil, err
	}
	copier.Copy(&fotoProdukResponses, &fotoProdukList)
	produkResponse.Photos = fotoProdukResponses

	return produkResponse, nil
}

func (p *produkUsecase) EditProdukByID(ctx context.Context, produkId int, userId int, req *model.ProdukRequest) (*model.ProdukResponse, error) {
	oldProduk, err := p.produkRepository.FindByID(ctx, produkId)
	if err != nil {
		return nil, err
	}
	toko, err := p.tokoRepository.FindByTokoID(ctx, oldProduk.IdToko)
	if err != nil {
		return nil, err
	}
	if toko.IdUser != userId {
		return nil, errors.New("unauthorized")
	}
	req.IdToko = toko.ID

	produkSlug := slug.Make(req.NamaProduk)
	req.Slug = produkSlug

	category, err := p.categoryRepository.FindByID(ctx, req.IdCategory)
	if err != nil {
		return nil, err
	}

	produk := new(model.Produk)
	copier.Copy(produk, req)

	fotoProdukList := []*model.FotoProduk{}
	for _, photoUrl := range req.PhotoUrls {
		fotoProduk := new(model.FotoProduk)
		fotoProduk.IdProduk = produkId
		fotoProduk.Url = photoUrl
		fotoProdukList = append(fotoProdukList, fotoProduk)
	}

	produk, fotoProdukList, err = p.produkRepository.UpdateProdukAndFotoProduk(ctx, produkId, produk, fotoProdukList)
	if err != nil {
		return nil, err
	}
	produkResponse := new(model.ProdukResponse)
	copier.Copy(produkResponse, produk)

	fotoProdukResponses := []*model.FotoProdukResponse{}
	copier.Copy(&fotoProdukResponses, &fotoProdukList)
	produkResponse.Photos = fotoProdukResponses

	tokoGetByIDResponse := new(model.TokoGetByIDResponse)
	copier.Copy(tokoGetByIDResponse, toko)
	produkResponse.Toko = tokoGetByIDResponse

	categoryResponse := new(model.CategoryResponse)
	copier.Copy(categoryResponse, category)
	produkResponse.Category = categoryResponse

	return produkResponse, nil
}

func (p *produkUsecase) DestroyProduk(ctx context.Context, produkId int, userId int) error {
	produk, err := p.produkRepository.FindByID(ctx, produkId)
	if err != nil {
		return err
	}
	toko, err := p.tokoRepository.FindByTokoID(ctx, produk.IdToko)
	if err != nil {
		return err
	}
	if toko.IdUser != userId {
		return errors.New("unauthorized")
	}

	err = p.produkRepository.DeleteProdukAndFotoProduk(ctx, produkId)
	if err != nil {
		return err
	}

	return nil
}

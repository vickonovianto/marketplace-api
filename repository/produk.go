package repository

import (
	"context"
	"errors"
	"marketplace-api/config"
	"marketplace-api/model"

	"gorm.io/gorm"
)

type produkRepository struct {
	Cfg config.Config
}

func NewProdukRepository(cfg config.Config) model.ProdukRepository {
	return &produkRepository{Cfg: cfg}
}

func (p *produkRepository) CreateProdukAndFotoProduk(
	ctx context.Context,
	produk *model.Produk,
	photoUrls []string,
) (*model.Produk, []*model.FotoProduk, error) {

	transaction := p.Cfg.Database().WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			transaction.Rollback()
		}
	}()

	if err := transaction.Error; err != nil {
		return nil, nil, err
	}

	if err := transaction.Create(&produk).Error; err != nil {
		transaction.Rollback()
		return nil, nil, err
	}

	fotoProdukList := []*model.FotoProduk{}
	for _, photoUrl := range photoUrls {
		fotoProduk := new(model.FotoProduk)
		fotoProduk.IdProduk = produk.ID
		fotoProduk.Url = photoUrl
		fotoProdukList = append(fotoProdukList, fotoProduk)
	}

	if err := transaction.Create(&fotoProdukList).Error; err != nil {
		transaction.Rollback()
		return nil, nil, err
	}

	return produk, fotoProdukList, transaction.Commit().Error
}

func (p *produkRepository) Fetch(ctx context.Context, req *model.ProdukFetchRequest) ([]*model.Produk, error) {
	var data []*model.Produk

	offset := (req.Page - 1) * req.Limit
	query := p.Cfg.Database().WithContext(ctx).Where("nama_produk LIKE ?", "%"+req.NamaProduk+"%")
	if req.CategoryId != -1 {
		query = query.Where("id_category = ?", req.CategoryId)
	}
	if req.TokoId != -1 {
		query = query.Where("id_toko = ?", req.TokoId)
	}
	if req.MinHarga != -1 {
		query = query.Where("harga_konsumen >= ?", req.MinHarga)
	}
	if req.MaxHarga != -1 {
		query = query.Where("harga_konsumen <= ?", req.MaxHarga)
	}
	if err := query.Limit(req.Limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (p *produkRepository) FindByID(ctx context.Context, produkId int) (*model.Produk, error) {
	produk := new(model.Produk)

	if err := p.Cfg.Database().
		WithContext(ctx).
		First(produk, produkId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("produk not found")
		}
		return nil, err
	}
	return produk, nil
}

func (p *produkRepository) UpdateProdukAndFotoProduk(
	ctx context.Context,
	produkId int,
	produk *model.Produk,
	fotoProdukList []*model.FotoProduk,
) (*model.Produk, []*model.FotoProduk, error) {

	transaction := p.Cfg.Database().WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			transaction.Rollback()
		}
	}()

	if err := transaction.Error; err != nil {
		return nil, nil, err
	}

	if err := transaction.
		Model(&model.Produk{ID: produkId}).Updates(produk).Find(produk).Error; err != nil {
		transaction.Rollback()
		return nil, nil, err
	}

	res := transaction.Delete(&model.FotoProduk{}, "id_produk = ?", produkId)
	if res.Error != nil {
		transaction.Rollback()
		return nil, nil, res.Error
	}

	if err := transaction.Create(&fotoProdukList).Error; err != nil {
		transaction.Rollback()
		return nil, nil, err
	}

	return produk, fotoProdukList, transaction.Commit().Error
}

func (p *produkRepository) DeleteProdukAndFotoProduk(ctx context.Context, produkId int) error {

	transaction := p.Cfg.Database().WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			transaction.Rollback()
		}
	}()

	if err := transaction.Error; err != nil {
		return err
	}

	res := transaction.Delete(&model.FotoProduk{}, "id_produk = ?", produkId)
	if res.Error != nil {
		transaction.Rollback()
		return res.Error
	}

	res = transaction.Delete(&model.Produk{}, produkId)
	if res.Error != nil {
		transaction.Rollback()
		return res.Error
	}

	return transaction.Commit().Error
}

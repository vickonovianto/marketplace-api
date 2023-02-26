package repository

import (
	"context"
	"marketplace-api/config"
	"marketplace-api/model"
)

type fotoProdukRepository struct {
	Cfg config.Config
}

func NewFotoProdukRepository(cfg config.Config) model.FotoProdukRepository {
	return &fotoProdukRepository{Cfg: cfg}
}

func (f *fotoProdukRepository) FetchByProdukId(ctx context.Context, produkId int) ([]*model.FotoProduk, error) {
	var data []*model.FotoProduk

	if err := f.Cfg.Database().WithContext(ctx).
		Where("id_produk = ?", produkId).
		Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

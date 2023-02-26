package repository

import (
	"context"
	"marketplace-api/config"
	"marketplace-api/model"
)

type detailTrxRepository struct {
	Cfg config.Config
}

func NewDetailTrxRepository(cfg config.Config) model.DetailTrxRepository {
	return &detailTrxRepository{Cfg: cfg}
}

func (d *detailTrxRepository) FindByTrxID(ctx context.Context, trxId int) ([]*model.DetailTrx, error) {
	var data []*model.DetailTrx

	if err := d.Cfg.Database().WithContext(ctx).
		Where("id_trx = ?", trxId).
		Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

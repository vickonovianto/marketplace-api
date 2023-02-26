package repository

import (
	"context"
	"errors"
	"marketplace-api/config"
	"marketplace-api/model"

	"gorm.io/gorm"
)

type logProdukRepository struct {
	Cfg config.Config
}

func NewLogProdukRepository(cfg config.Config) model.LogProdukRepository {
	return &logProdukRepository{Cfg: cfg}
}

func (l *logProdukRepository) FindByID(ctx context.Context, logProdukId int) (*model.LogProduk, error) {
	logProduk := new(model.LogProduk)

	if err := l.Cfg.Database().
		WithContext(ctx).
		First(logProduk, logProdukId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("log produk not found")
		}
		return nil, err
	}
	return logProduk, nil
}

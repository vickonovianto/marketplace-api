package repository

import (
	"context"
	"errors"
	"marketplace-api/config"
	"marketplace-api/model"

	"gorm.io/gorm"
)

type tokoRepository struct {
	Cfg config.Config
}

func NewTokoRepository(cfg config.Config) model.TokoRepository {
	return &tokoRepository{Cfg: cfg}
}

func (t *tokoRepository) Create(ctx context.Context, toko *model.Toko) (*model.Toko, error) {
	if err := t.Cfg.Database().WithContext(ctx).Create(&toko).Error; err != nil {
		return nil, err
	}
	return toko, nil
}

func (t *tokoRepository) FetchAndPaginate(ctx context.Context, req *model.TokoFetchPaginateRequest) ([]*model.Toko, error) {
	var data []*model.Toko

	offset := (req.Page - 1) * req.Limit
	if err := t.Cfg.Database().WithContext(ctx).
		Where("nama_toko LIKE ?", "%"+req.Nama+"%").
		Limit(req.Limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (t *tokoRepository) FindByTokoID(ctx context.Context, tokoId int) (*model.Toko, error) {
	toko := new(model.Toko)

	if err := t.Cfg.Database().
		WithContext(ctx).
		First(toko, tokoId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("toko not found")
		}
		return nil, err
	}
	return toko, nil
}

func (t *tokoRepository) FindByUserID(ctx context.Context, userId int) (*model.Toko, error) {
	toko := new(model.Toko)

	if err := t.Cfg.Database().
		WithContext(ctx).
		Where("id_user = ?", userId).
		First(toko).Error; err != nil {
		return nil, err
	}
	return toko, nil
}

func (t *tokoRepository) UpdateByTokoID(ctx context.Context, tokoId int, toko *model.Toko) (*model.Toko, error) {
	_, err := t.FindByTokoID(ctx, tokoId)
	if err != nil {
		return nil, err
	}

	if err := t.Cfg.Database().WithContext(ctx).
		Model(&model.Toko{ID: tokoId}).Updates(toko).Find(toko).Error; err != nil {
		return nil, err
	}
	return toko, nil
}

package repository

import (
	"context"
	"errors"
	"marketplace-api/config"
	"marketplace-api/model"

	"gorm.io/gorm"
)

type alamatRepository struct {
	Cfg config.Config
}

func NewAlamatRepository(cfg config.Config) model.AlamatRepository {
	return &alamatRepository{Cfg: cfg}
}

func (a *alamatRepository) Create(ctx context.Context, alamat *model.Alamat) (*model.Alamat, error) {
	if err := a.Cfg.Database().WithContext(ctx).Create(&alamat).Error; err != nil {
		return nil, err
	}
	return alamat, nil
}

func (a *alamatRepository) FetchAndFilter(ctx context.Context, userId int, judulAlamat string) ([]*model.Alamat, error) {
	var data []*model.Alamat

	if err := a.Cfg.Database().WithContext(ctx).
		Where("id_user = ? AND judul_alamat LIKE ?", userId, "%"+judulAlamat+"%").
		Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (a *alamatRepository) FindByID(ctx context.Context, alamatId int) (*model.Alamat, error) {
	alamat := new(model.Alamat)

	if err := a.Cfg.Database().
		WithContext(ctx).
		First(alamat, alamatId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("alamat not found")
		}
		return nil, err
	}
	return alamat, nil
}

func (a *alamatRepository) UpdateByID(ctx context.Context, alamatId int, alamat *model.Alamat) (*model.Alamat, error) {
	_, err := a.FindByID(ctx, alamatId)
	if err != nil {
		return nil, err
	}

	if err := a.Cfg.Database().WithContext(ctx).
		Model(&model.Alamat{ID: alamatId}).Updates(alamat).Find(alamat).Error; err != nil {
		return nil, err
	}
	return alamat, nil
}

func (a *alamatRepository) Delete(ctx context.Context, alamatId int) error {
	_, err := a.FindByID(ctx, alamatId)
	if err != nil {
		return err
	}

	res := a.Cfg.Database().WithContext(ctx).
		Delete(&model.Alamat{}, alamatId)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

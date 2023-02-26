package repository

import (
	"context"
	"errors"
	"marketplace-api/config"
	"marketplace-api/model"

	"gorm.io/gorm"
)

type userRepository struct {
	Cfg config.Config
}

func NewUserRepository(cfg config.Config) model.UserRepository {
	return &userRepository{Cfg: cfg}
}

func (u *userRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	if err := u.Cfg.Database().WithContext(ctx).Create(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepository) FindByNoTelp(ctx context.Context, noTelp string) (*model.User, error) {
	user := new(model.User)
	if err := u.Cfg.Database().
		WithContext(ctx).
		Where("no_telp = ?", noTelp).
		First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no telp atau kata sandi salah")
		}
		return nil, err
	}
	return user, nil
}

func (u *userRepository) FindByID(ctx context.Context, id int) (*model.User, error) {
	user := new(model.User)

	if err := u.Cfg.Database().
		WithContext(ctx).
		First(user, id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepository) UpdateByID(ctx context.Context, id int, user *model.User) (*model.User, error) {
	_, err := u.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := u.Cfg.Database().WithContext(ctx).
		Model(&model.User{ID: id}).Updates(user).Find(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

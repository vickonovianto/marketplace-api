package repository

import (
	"context"
	"errors"
	"marketplace-api/config"
	"marketplace-api/model"

	"gorm.io/gorm"
)

type categoryRepository struct {
	Cfg config.Config
}

func NewCategoryRepository(cfg config.Config) model.CategoryRepository {
	return &categoryRepository{Cfg: cfg}
}

func (c *categoryRepository) Create(ctx context.Context, category *model.Category) (*model.Category, error) {
	if err := c.Cfg.Database().WithContext(ctx).Create(&category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (c *categoryRepository) FetchAll(ctx context.Context) ([]*model.Category, error) {
	var data []*model.Category

	if err := c.Cfg.Database().WithContext(ctx).
		Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (c *categoryRepository) FindByID(ctx context.Context, id int) (*model.Category, error) {
	category := new(model.Category)

	if err := c.Cfg.Database().
		WithContext(ctx).
		First(category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return category, nil
}

func (c *categoryRepository) UpdateByID(ctx context.Context, id int, category *model.Category) (*model.Category, error) {
	_, err := c.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := c.Cfg.Database().WithContext(ctx).
		Model(&model.Category{ID: id}).Updates(category).Find(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (c *categoryRepository) Delete(ctx context.Context, id int) error {
	_, err := c.FindByID(ctx, id)
	if err != nil {
		return err
	}

	res := c.Cfg.Database().WithContext(ctx).
		Delete(&model.Category{}, id)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

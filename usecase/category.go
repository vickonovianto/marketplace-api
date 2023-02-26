package usecase

import (
	"context"
	"marketplace-api/model"

	"github.com/jinzhu/copier"
)

type categoryUsecase struct {
	categoryRepository model.CategoryRepository
}

func NewCategoryUsecase(categoryRepository model.CategoryRepository) model.CategoryUsecase {
	return &categoryUsecase{categoryRepository: categoryRepository}
}

func (c *categoryUsecase) StoreCategory(ctx context.Context, req *model.CategoryRequest) (*model.CategoryResponse, error) {
	category := new(model.Category)
	copier.Copy(category, req)
	category, err := c.categoryRepository.Create(ctx, category)
	if err != nil {
		return nil, err
	}
	categoryResponse := new(model.CategoryResponse)
	copier.Copy(categoryResponse, category)
	return categoryResponse, nil
}

func (c *categoryUsecase) FetchAllCategory(ctx context.Context) ([]*model.CategoryResponse, error) {
	categories, err := c.categoryRepository.FetchAll(ctx)
	if err != nil {
		return nil, err
	}
	categoryResponses := []*model.CategoryResponse{}
	copier.Copy(&categoryResponses, &categories)
	return categoryResponses, nil
}

func (c *categoryUsecase) GetCategoryByID(ctx context.Context, id int) (*model.CategoryResponse, error) {
	category, err := c.categoryRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	categoryResponse := new(model.CategoryResponse)
	copier.Copy(categoryResponse, category)
	return categoryResponse, nil
}

func (c *categoryUsecase) EditCategory(ctx context.Context, id int, req *model.CategoryRequest) (*model.CategoryResponse, error) {
	category := new(model.Category)
	copier.Copy(category, req)
	category, err := c.categoryRepository.UpdateByID(ctx, id, category)
	if err != nil {
		return nil, err
	}
	categoryResponse := new(model.CategoryResponse)
	copier.Copy(categoryResponse, category)
	return categoryResponse, nil
}

func (c *categoryUsecase) DestroyCategory(ctx context.Context, id int) error {
	err := c.categoryRepository.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

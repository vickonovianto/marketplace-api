package model

import (
	"context"
	"time"
)

type (
	Category struct {
		ID           int       `gorm:"column:id"`
		NamaCategory string    `gorm:"column:nama_category;size:255;not null"`
		CreatedAt    time.Time `gorm:"column:created_at"`
		UpdatedAt    time.Time `gorm:"column:updated_at"`
	}

	CategoryRepository interface {
		Create(ctx context.Context, category *Category) (*Category, error)
		FetchAll(ctx context.Context) ([]*Category, error)
		FindByID(ctx context.Context, id int) (*Category, error)
		UpdateByID(ctx context.Context, id int, category *Category) (*Category, error)
		Delete(ctx context.Context, id int) error
	}

	CategoryUsecase interface {
		StoreCategory(ctx context.Context, req *CategoryRequest) (*CategoryResponse, error)
		FetchAllCategory(ctx context.Context) ([]*CategoryResponse, error)
		GetCategoryByID(ctx context.Context, id int) (*CategoryResponse, error)
		EditCategory(ctx context.Context, id int, req *CategoryRequest) (*CategoryResponse, error)
		DestroyCategory(ctx context.Context, id int) error
	}

	CategoryRequest struct {
		NamaCategory string `json:"nama_category"`
	}

	CategoryResponse struct {
		ID           int    `json:"id"`
		NamaCategory string `json:"nama_category"`
	}
)

// override gorm table name
func (Category) TableName() string {
	return "category"
}

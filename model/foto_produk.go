package model

import (
	"context"
	"time"
)

type (
	FotoProduk struct {
		ID        int       `gorm:"column:id"`
		IdProduk  int       `gorm:"column:id_produk"`
		Produk    *Produk   `gorm:"foreignKey:IdProduk"`
		Url       string    `gorm:"column:url;size:255;not null"`
		CreatedAt time.Time `gorm:"column:created_at"`
		UpdatedAt time.Time `gorm:"column:updated_at"`
	}

	FotoProdukRepository interface {
		FetchByProdukId(ctx context.Context, produkId int) ([]*FotoProduk, error)
	}

	FotoProdukResponse struct {
		ID       int    `json:"id"`
		IdProduk int    `json:"product_id"`
		Url      string `json:"url"`
	}
)

// override gorm table name
func (FotoProduk) TableName() string {
	return "foto_produk"
}

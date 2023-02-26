package model

import (
	"context"
	"time"
)

type (
	// Trx means Transaction
	LogProduk struct {
		ID            int       `gorm:"column:id"`
		IdProduk      int       `gorm:"column:id_produk;not null"`
		Produk        *Produk   `gorm:"foreignKey:IdProduk"`
		NamaProduk    string    `gorm:"column:nama_produk;size:255;not null"`
		Slug          string    `gorm:"column:slug;size:255;not null"`
		HargaReseller string    `gorm:"column:harga_reseller;size:255;not null"`
		HargaKonsumen string    `gorm:"column:harga_konsumen;size:255;not null"`
		Deskripsi     string    `gorm:"column:deskripsi;not null"`
		CreatedAt     time.Time `gorm:"column:created_at"`
		UpdatedAt     time.Time `gorm:"column:updated_at"`
		IdToko        int       `gorm:"column:id_toko;not null"`
		Toko          *Toko     `gorm:"foreignKey:IdToko"`
		IdCategory    int       `gorm:"column:id_category;not null"`
		Category      *Category `gorm:"foreignKey:IdCategory"`
	}

	LogProdukRepository interface {
		FindByID(ctx context.Context, logProdukId int) (*LogProduk, error)
	}

	LogProdukRequest struct {
		IdProduk      int
		NamaProduk    string
		Slug          string
		HargaReseller string
		HargaKonsumen string
		Deskripsi     string
		IdToko        int
		IdCategory    int
	}

	LogProdukResponse struct {
		IdProduk      int                    `json:"id"`
		NamaProduk    string                 `json:"nama_produk"`
		Slug          string                 `json:"slug"`
		HargaReseller string                 `json:"harga_reseller"`
		HargaKonsumen string                 `json:"harga_konsumen"`
		Deskripsi     string                 `json:"deskripsi"`
		Toko          *TokoLogProdukResponse `json:"toko"`
		Category      *CategoryResponse      `json:"category"`
		Photos        []*FotoProdukResponse  `json:"photos"`
	}
)

// override gorm table name
func (LogProduk) TableName() string {
	return "log_produk"
}

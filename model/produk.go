package model

import (
	"context"
	"time"
)

type (
	Produk struct {
		ID            int       `gorm:"column:id"`
		NamaProduk    string    `gorm:"column:nama_produk;size:255;not null"`
		Slug          string    `gorm:"column:slug;size:255;not null"`
		HargaReseller string    `gorm:"column:harga_reseller;size:255;not null"`
		HargaKonsumen string    `gorm:"column:harga_konsumen;size:255;not null"`
		Stok          int       `gorm:"column:stok;not null"`
		Deskripsi     string    `gorm:"column:deskripsi;not null"`
		IdToko        int       `gorm:"column:id_toko"`
		Toko          *Toko     `gorm:"foreignKey:IdToko"`
		IdCategory    int       `gorm:"column:id_category"`
		Category      *Category `gorm:"foreignKey:IdCategory"`
		CreatedAt     time.Time `gorm:"column:created_at"`
		UpdatedAt     time.Time `gorm:"column:updated_at"`
	}

	ProdukRepository interface {
		CreateProdukAndFotoProduk(
			ctx context.Context,
			produk *Produk,
			photoUrls []string,
		) (*Produk, []*FotoProduk, error)
		Fetch(ctx context.Context, req *ProdukFetchRequest) ([]*Produk, error)
		FindByID(ctx context.Context, produkId int) (*Produk, error)
		UpdateProdukAndFotoProduk(
			ctx context.Context,
			produkId int,
			produk *Produk,
			fotoProdukList []*FotoProduk,
		) (*Produk, []*FotoProduk, error)
		DeleteProdukAndFotoProduk(ctx context.Context, produkId int) error
	}

	ProdukUsecase interface {
		StoreProduk(ctx context.Context, req *ProdukRequest, userId int) (*ProdukResponse, error)
		FetchProduk(ctx context.Context, req *ProdukFetchRequest) ([]*ProdukResponse, error)
		GetProdukByID(ctx context.Context, produkId int) (*ProdukResponse, error)
		EditProdukByID(ctx context.Context, produkId int, userId int, req *ProdukRequest) (*ProdukResponse, error)
		DestroyProduk(ctx context.Context, produkId int, userId int) error
	}

	ProdukFetchRequest struct {
		NamaProduk string
		Limit      int
		Page       int
		CategoryId int
		TokoId     int
		MaxHarga   int
		MinHarga   int
	}

	ProdukRequest struct {
		NamaProduk    string `json:"nama_produk"`
		Slug          string
		HargaReseller string `json:"harga_reseller"`
		HargaKonsumen string `json:"harga_konsumen"`
		Stok          int    `json:"stok"`
		Deskripsi     string `json:"deskripsi"`
		IdToko        int
		IdCategory    int `json:"id_category"`
		PhotoUrls     []string
	}

	ProdukResponse struct {
		ID            int                   `json:"id"`
		NamaProduk    string                `json:"nama_produk"`
		Slug          string                `json:"slug"`
		HargaReseller string                `json:"harga_reseller"`
		HargaKonsumen string                `json:"harga_konsumen"`
		Stok          int                   `json:"stok"`
		Deskripsi     string                `json:"deskripsi"`
		Toko          *TokoGetByIDResponse  `json:"toko"`
		Category      *CategoryResponse     `json:"category"`
		Photos        []*FotoProdukResponse `json:"photos"`
	}
)

// override gorm table name
func (Produk) TableName() string {
	return "produk"
}

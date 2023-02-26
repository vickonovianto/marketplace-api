package model

import (
	"context"
	"time"
)

type (
	// Trx means Transaction
	DetailTrx struct {
		ID          int        `gorm:"column:id"`
		IdTrx       int        `gorm:"column:id_trx;not null"`
		Trx         *Trx       `gorm:"foreignKey:IdTrx"`
		IdLogProduk int        `gorm:"column:id_log_produk;not null"`
		LogProduk   *LogProduk `gorm:"foreignKey:IdLogProduk"`
		IdToko      int        `gorm:"column:id_toko;not null"`
		Toko        *Toko      `gorm:"foreignKey:IdToko"`
		Kuantitas   int        `gorm:"column:kuantitas;not null"`
		HargaTotal  int        `gorm:"column:harga_total;not null"`
		CreatedAt   time.Time  `gorm:"column:created_at"`
		UpdatedAt   time.Time  `gorm:"column:updated_at"`
	}

	DetailTrxRepository interface {
		FindByTrxID(ctx context.Context, trxId int) ([]*DetailTrx, error)
	}

	DetailTrxWithLogProduk struct {
		DetailTrx *DetailTrx
		LogProduk *LogProduk
	}

	DetailTrxRequest struct {
		ProductId int `json:"product_id"`
		Kuantitas int `json:"kuantitas"`
	}

	DetailTrxResponse struct {
		LogProduk  *LogProdukResponse   `json:"product"`
		Toko       *TokoGetByIDResponse `json:"toko"`
		Kuantitas  int                  `json:"kuantitas"`
		HargaTotal int                  `json:"harga_total"`
	}
)

// override gorm table name
func (DetailTrx) TableName() string {
	return "detail_trx"
}

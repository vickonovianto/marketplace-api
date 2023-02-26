package model

import (
	"context"
	"time"
)

const KODE_INVOICE_PREFIX = "INV-"

type (
	// Trx means Transaction
	Trx struct {
		ID               int       `gorm:"column:id"`
		IdUser           int       `gorm:"column:id_user;not null"`
		User             *User     `gorm:"foreignKey:IdUser"`
		AlamatPengiriman int       `gorm:"column:alamat_pengiriman;not null"`
		Alamat           *Alamat   `gorm:"foreignKey:AlamatPengiriman"`
		HargaTotal       int       `gorm:"column:harga_total;not null"`
		KodeInvoice      string    `gorm:"column:kode_invoice;size:255;not null"`
		MethodBayar      string    `gorm:"column:method_bayar;size:255;not null"`
		CreatedAt        time.Time `gorm:"column:created_at"`
		UpdatedAt        time.Time `gorm:"column:updated_at"`
	}

	TrxRepository interface {
		CreateTrx(
			ctx context.Context,
			trx *Trx,
			detailTrxWithLogProdukList []*DetailTrxWithLogProduk,
		) (*Trx, error)
		Fetch(ctx context.Context, req *TrxFetchRequest, userId int) ([]*Trx, error)
		FindByID(ctx context.Context, trxId int) (*Trx, error)
	}

	TrxUsecase interface {
		StoreTrx(ctx context.Context, req *TrxStoreRequest, userId int) (*TrxGetByIDResponse, error)
		FetchTrx(ctx context.Context, req *TrxFetchRequest, userId int) (*TrxFetchResponse, error)
		GetTrxByID(ctx context.Context, trxId int, userId int) (*TrxGetByIDResponse, error)
	}

	TrxStoreRequest struct {
		MethodBayar       string              `json:"method_bayar"`
		AlamatPengiriman  int                 `json:"alamat_kirim"`
		DetailTrxRequests []*DetailTrxRequest `json:"detail_trx"`
	}

	TrxFetchRequest struct {
		Search string
		Limit  int
		Page   int
	}

	TrxFetchResponse struct {
		Limit int                   `json:"limit"`
		Page  int                   `json:"page"`
		Data  []*TrxGetByIDResponse `json:"data"`
	}

	TrxGetByIDResponse struct {
		ID                 int                  `json:"id"`
		HargaTotal         int                  `json:"harga_total"`
		KodeInvoice        string               `json:"kode_invoice"`
		MethodBayar        string               `json:"method_bayar"`
		AlamatPengiriman   *AlamatResponse      `json:"alamat_kirim"`
		DetailTrxResponses []*DetailTrxResponse `json:"detail_trx"`
	}
)

// override gorm table name
func (Trx) TableName() string {
	return "trx"
}

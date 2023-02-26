package model

import (
	"context"
	"time"
)

type (
	Toko struct {
		ID        int       `gorm:"column:id"`
		IdUser    int       `gorm:"column:id_user"`
		User      *User     `gorm:"foreignKey:IdUser"`
		NamaToko  string    `gorm:"column:nama_toko;size:255;not null"`
		UrlFoto   string    `gorm:"column:url_foto;size:255;not null"`
		CreatedAt time.Time `gorm:"column:created_at"`
		UpdatedAt time.Time `gorm:"column:updated_at"`
	}

	TokoRepository interface {
		Create(ctx context.Context, toko *Toko) (*Toko, error)
		FetchAndPaginate(ctx context.Context, req *TokoFetchPaginateRequest) ([]*Toko, error)
		FindByTokoID(ctx context.Context, tokoId int) (*Toko, error)
		FindByUserID(ctx context.Context, userId int) (*Toko, error)
		UpdateByTokoID(ctx context.Context, tokoId int, toko *Toko) (*Toko, error)
	}

	TokoUsecase interface {
		FetchAndPaginateToko(ctx context.Context, req *TokoFetchPaginateRequest) (*TokoFetchPaginateResponse, error)
		GetTokoByID(ctx context.Context, tokoId int) (*TokoGetByIDResponse, error)
		GetMyToko(ctx context.Context, userId int) (*GetMyTokoResponse, error)
		EditToko(ctx context.Context, req *TokoUpdateRequest) (*TokoUpdateResponse, error)
	}

	TokoCreateRequest struct {
		IdUser   int    `json:"id_user"`
		NamaToko string `json:"nama_toko"`
		UrlFoto  string `json:"url_foto"`
	}

	TokoFetchPaginateRequest struct {
		Limit int
		Page  int
		Nama  string
	}

	TokoUpdateRequest struct {
		ID       int
		NamaToko string `json:"nama_toko"`
		UrlFoto  string `json:"url_foto"`
		IdUser   int
	}

	TokoCreateResponse struct {
		ID       int    `json:"id"`
		NamaToko string `json:"nama_toko"`
		UrlFoto  string `json:"url_foto"`
		IdUser   int    `json:"user_id"`
	}

	TokoFetchPaginateResponse struct {
		Limit int                    `json:"limit"`
		Page  int                    `json:"page"`
		Data  []*TokoGetByIDResponse `json:"data"`
	}

	GetMyTokoResponse struct {
		ID       int    `json:"id"`
		NamaToko string `json:"nama_toko"`
		UrlFoto  string `json:"url_foto"`
		IdUser   int    `json:"user_id"`
	}

	TokoGetByIDResponse struct {
		ID       int    `json:"id"`
		NamaToko string `json:"nama_toko"`
		UrlFoto  string `json:"url_foto"`
	}

	TokoLogProdukResponse struct {
		NamaToko string `json:"nama_toko"`
		UrlFoto  string `json:"url_foto"`
	}

	TokoUpdateResponse struct {
		ID       int    `json:"id"`
		NamaToko string `json:"nama_toko"`
		UrlFoto  string `json:"url_foto"`
		IdUser   int    `json:"user_id"`
	}
)

// override gorm table name
func (Toko) TableName() string {
	return "toko"
}

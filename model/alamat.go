package model

import (
	"context"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type (
	Alamat struct {
		ID           int       `gorm:"column:id"`
		IdUser       int       `gorm:"column:id_user"`
		User         *User     `gorm:"foreignKey:IdUser"`
		JudulAlamat  string    `gorm:"column:judul_alamat;size:255;not null"`
		NamaPenerima string    `gorm:"column:nama_penerima;size:255;not null"`
		NoTelp       string    `gorm:"column:no_telp;size:255;not null"`
		DetailAlamat string    `gorm:"column:detail_alamat;size:255;not null"`
		CreatedAt    time.Time `gorm:"column:created_at"`
		UpdatedAt    time.Time `gorm:"column:updated_at"`
	}

	AlamatRepository interface {
		Create(ctx context.Context, alamat *Alamat) (*Alamat, error)
		FetchAndFilter(ctx context.Context, userId int, judulAlamat string) ([]*Alamat, error)
		FindByID(ctx context.Context, alamatId int) (*Alamat, error)
		UpdateByID(ctx context.Context, alamatId int, alamat *Alamat) (*Alamat, error)
		Delete(ctx context.Context, alamatId int) error
	}

	AlamatUsecase interface {
		StoreAlamat(ctx context.Context, req *AlamatRequest) (*AlamatResponse, error)
		FetchAndFilterAlamat(ctx context.Context, userId int, judulAlamat string) ([]*AlamatResponse, error)
		GetAlamatByID(ctx context.Context, alamatId int, userId int) (*AlamatResponse, error)
		EditAlamatByID(ctx context.Context, alamatId int, req *AlamatRequest) (*AlamatResponse, error)
		DestroyAlamat(ctx context.Context, alamatId int, userId int) error
	}

	AlamatRequest struct {
		IdUser       int
		JudulAlamat  string `json:"judul_alamat"`
		NamaPenerima string `json:"nama_penerima"`
		NoTelp       string `json:"no_telp"`
		DetailAlamat string `json:"detail_alamat"`
	}

	AlamatResponse struct {
		ID           int    `json:"id"`
		JudulAlamat  string `json:"judul_alamat"`
		NamaPenerima string `json:"nama_penerima"`
		NoTelp       string `json:"no_telp"`
		DetailAlamat string `json:"detail_alamat"`
	}
)

// override gorm table name
func (Alamat) TableName() string {
	return "alamat"
}

func (req AlamatRequest) Validate() error {
	return validation.ValidateStruct(
		&req,
		validation.Field(&req.JudulAlamat, validation.Required, validation.Length(1, 255)),
		validation.Field(&req.NamaPenerima, validation.Required, validation.Length(1, 255)),
		validation.Field(&req.NoTelp, validation.Required, is.Digit, validation.Length(10, 13)),
		validation.Field(&req.DetailAlamat, validation.Required, validation.Length(1, 255)),
	)
}

func (req *AlamatRequest) Trim() {
	req.JudulAlamat = strings.TrimSpace(req.JudulAlamat)
	req.NamaPenerima = strings.TrimSpace(req.NamaPenerima)
	req.NoTelp = strings.TrimSpace(req.NoTelp)
	req.DetailAlamat = strings.TrimSpace(req.DetailAlamat)
}

package model

import (
	"context"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

const TANGGAL_LAHIR_DATE_FORMAT = "02/01/2006"

type (
	User struct {
		ID           int       `gorm:"column:id"`
		Nama         string    `gorm:"column:nama;size:255;not null"`
		KataSandi    string    `gorm:"column:kata_sandi;size:255;not null"`
		NoTelp       string    `gorm:"column:no_telp;size:255;not null;unique"`
		TanggalLahir time.Time `gorm:"column:tanggal_lahir;type:DATE NOT NULL"`
		JenisKelamin string    `gorm:"column:jenis_kelamin;size:255;not null"`
		Tentang      string    `gorm:"column:tentang;not null"`
		Pekerjaan    string    `gorm:"column:pekerjaan;size:255;not null"`
		Email        string    `gorm:"column:email;size:255;not null;unique"`
		IdProvinsi   string    `gorm:"column:id_provinsi;size:255;not null"`
		IdKota       string    `gorm:"column:id_kota;size:255;not null"`
		IsAdmin      bool      `gorm:"column:is_admin;not null;default:0"`
		CreatedAt    time.Time `gorm:"column:created_at"`
		UpdatedAt    time.Time `gorm:"column:updated_at"`
	}

	UserRepository interface {
		Create(ctx context.Context, user *User) (*User, error)
		FindByNoTelp(ctx context.Context, noTelp string) (*User, error)
		FindByID(ctx context.Context, id int) (*User, error)
		UpdateByID(ctx context.Context, id int, user *User) (*User, error)
	}

	UserUsecase interface {
		RegisterUser(ctx context.Context, req *UserRegisterRequest) (*UserRegisterResponse, error)
		LoginUser(ctx context.Context, req *UserLoginRequest) (*UserLoginResponse, error)
		GetCurrentUser(ctx context.Context, userId int) (*UserResponse, error)
		EditCurrentUser(ctx context.Context, userId int, req *UserUpdateRequest) (*UserResponse, error)
	}

	UserRegisterRequest struct {
		Nama         string `json:"nama"`
		KataSandi    string `json:"kata_sandi"`
		NoTelp       string `json:"no_telp"`
		TanggalLahir string `json:"tanggal_lahir"`
		JenisKelamin string `json:"jenis_kelamin"`
		Tentang      string `json:"tentang"`
		Pekerjaan    string `json:"pekerjaan"`
		Email        string `json:"email"`
		IdProvinsi   string `json:"id_provinsi"`
		IdKota       string `json:"id_kota"`
	}

	UserLoginRequest struct {
		NoTelp    string `json:"no_telp"`
		KataSandi string `json:"kata_sandi"`
	}

	UserUpdateRequest struct {
		Nama         string `json:"nama"`
		KataSandi    string `json:"kata_sandi"`
		NoTelp       string `json:"no_telp"`
		TanggalLahir string `json:"tanggal_lahir"`
		JenisKelamin string `json:"jenis_kelamin"`
		Tentang      string `json:"tentang"`
		Pekerjaan    string `json:"pekerjaan"`
		Email        string `json:"email"`
		IdProvinsi   string `json:"id_provinsi"`
		IdKota       string `json:"id_kota"`
	}

	UserRegisterResponse struct {
		Nama         string              `json:"nama"`
		NoTelp       string              `json:"no_telp"`
		TanggalLahir string              `json:"tanggal_lahir"`
		JenisKelamin string              `json:"jenis_kelamin"`
		Tentang      string              `json:"tentang"`
		Pekerjaan    string              `json:"pekerjaan"`
		Email        string              `json:"email"`
		IdProvinsi   *Province           `json:"id_provinsi"`
		IdKota       *City               `json:"id_kota"`
		Toko         *TokoCreateResponse `json:"toko"`
	}

	UserResponse struct {
		Nama         string    `json:"nama"`
		NoTelp       string    `json:"no_telp"`
		TanggalLahir string    `json:"tanggal_lahir"`
		JenisKelamin string    `json:"jenis_kelamin"`
		Tentang      string    `json:"tentang"`
		Pekerjaan    string    `json:"pekerjaan"`
		Email        string    `json:"email"`
		IdProvinsi   *Province `json:"id_provinsi"`
		IdKota       *City     `json:"id_kota"`
	}

	UserLoginResponse struct {
		Nama         string    `json:"nama"`
		NoTelp       string    `json:"no_telp"`
		TanggalLahir string    `json:"tanggal_lahir"`
		JenisKelamin string    `json:"jenis_kelamin"`
		Tentang      string    `json:"tentang"`
		Pekerjaan    string    `json:"pekerjaan"`
		Email        string    `json:"email"`
		IdProvinsi   *Province `json:"id_provinsi"`
		IdKota       *City     `json:"id_kota"`
		Token        string    `json:"token"`
	}
)

// override gorm table name
func (User) TableName() string {
	return "user"
}

func (req UserRegisterRequest) Validate() error {
	return validation.ValidateStruct(
		&req,
		validation.Field(&req.Nama, validation.Required, validation.Length(1, 255)),
		validation.Field(&req.KataSandi, validation.Required, validation.Length(6, 255)),
		validation.Field(&req.NoTelp, validation.Required, is.Digit, validation.Length(10, 13)),
		validation.Field(&req.TanggalLahir, validation.Required,
			validation.Date(TANGGAL_LAHIR_DATE_FORMAT).Max(time.Now()).Error("invalid or incorrect format, must be in format: dd/mm/yyyy")),
		validation.Field(&req.JenisKelamin, validation.Required, validation.Length(1, 255)),
		validation.Field(&req.Pekerjaan, validation.Required, validation.Length(1, 255)),
		validation.Field(&req.Email, validation.Required, is.Email, validation.Length(3, 255)),
		validation.Field(&req.IdProvinsi, validation.Required, is.Digit, validation.Length(2, 2)),
		validation.Field(&req.IdKota, validation.Required, is.Digit, validation.Length(4, 4)),
	)
}

func (req UserLoginRequest) Validate() error {
	return validation.ValidateStruct(
		&req,
		validation.Field(&req.NoTelp, validation.Required, is.Digit, validation.Length(10, 13)),
		validation.Field(&req.KataSandi, validation.Required, validation.Length(6, 255)),
	)
}

func (req UserUpdateRequest) Validate() error {
	return validation.ValidateStruct(
		&req,
		validation.Field(&req.Nama, validation.Required, validation.Length(1, 255)),
		validation.Field(&req.KataSandi, validation.Required, validation.Length(6, 255)),
		validation.Field(&req.NoTelp, validation.Required, is.Digit, validation.Length(10, 13)),
		validation.Field(&req.TanggalLahir, validation.Required,
			validation.Date(TANGGAL_LAHIR_DATE_FORMAT).Max(time.Now()).Error("invalid or incorrect format, must be in format: dd/mm/yyyy")),
		validation.Field(&req.JenisKelamin, validation.Required, validation.Length(1, 255)),
		validation.Field(&req.Pekerjaan, validation.Required, validation.Length(1, 255)),
		validation.Field(&req.Email, validation.Required, is.Email, validation.Length(3, 255)),
		validation.Field(&req.IdProvinsi, validation.Required, is.Digit, validation.Length(2, 2)),
		validation.Field(&req.IdKota, validation.Required, is.Digit, validation.Length(4, 4)),
	)
}

func (req *UserRegisterRequest) Trim() {
	req.Nama = strings.TrimSpace(req.Nama)
	req.KataSandi = strings.TrimSpace(req.KataSandi)
	req.NoTelp = strings.TrimSpace(req.NoTelp)
	req.TanggalLahir = strings.TrimSpace(req.TanggalLahir)
	req.JenisKelamin = strings.TrimSpace(req.JenisKelamin)
	req.Tentang = strings.TrimSpace(req.Tentang)
	req.Pekerjaan = strings.TrimSpace(req.Pekerjaan)
	req.Email = strings.TrimSpace(req.Email)
	req.IdProvinsi = strings.TrimSpace(req.IdProvinsi)
	req.IdKota = strings.TrimSpace(req.IdKota)
}

func (req *UserLoginRequest) Trim() {
	req.NoTelp = strings.TrimSpace(req.NoTelp)
	req.KataSandi = strings.TrimSpace(req.KataSandi)
}

func (req *UserUpdateRequest) Trim() {
	req.Nama = strings.TrimSpace(req.Nama)
	req.KataSandi = strings.TrimSpace(req.KataSandi)
	req.NoTelp = strings.TrimSpace(req.NoTelp)
	req.TanggalLahir = strings.TrimSpace(req.TanggalLahir)
	req.JenisKelamin = strings.TrimSpace(req.JenisKelamin)
	req.Tentang = strings.TrimSpace(req.Tentang)
	req.Pekerjaan = strings.TrimSpace(req.Pekerjaan)
	req.Email = strings.TrimSpace(req.Email)
	req.IdProvinsi = strings.TrimSpace(req.IdProvinsi)
	req.IdKota = strings.TrimSpace(req.IdKota)
}

package usecase

import (
	"context"
	"errors"
	"marketplace-api/model"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepository     model.UserRepository
	tokoRepository     model.TokoRepository
	provinceRepository model.ProvinceRepository
	cityRepository     model.CityRepository
}

func NewUserUsecase(
	userRepository model.UserRepository,
	tokoRepository model.TokoRepository,
	provinceRepository model.ProvinceRepository,
	cityRepository model.CityRepository,
) model.UserUsecase {
	return &userUsecase{
		userRepository:     userRepository,
		tokoRepository:     tokoRepository,
		provinceRepository: provinceRepository,
		cityRepository:     cityRepository,
	}
}

func (u *userUsecase) RegisterUser(ctx context.Context, req *model.UserRegisterRequest) (*model.UserRegisterResponse, error) {
	user := new(model.User)

	tanggalLahir, err := time.Parse(model.TANGGAL_LAHIR_DATE_FORMAT, req.TanggalLahir)
	if err != nil {
		return nil, err
	}
	user.TanggalLahir = tanggalLahir

	_, err = u.provinceRepository.FindByID(ctx, req.IdProvinsi)
	if err != nil {
		return nil, err
	}

	if req.IdKota[0:2] != req.IdProvinsi {
		return nil, errors.New("2 angka awal id kota harus sama dengan id provinsi")
	}
	_, err = u.cityRepository.FindByID(ctx, req.IdKota[0:2], req.IdKota)
	if err != nil {
		return nil, err
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(req.KataSandi), 14)
	if err != nil {
		return nil, err
	}
	hashedPassword := string(bytes)
	req.KataSandi = hashedPassword

	copier.Copy(user, req)
	user, err = u.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	userRegisterResponse := new(model.UserRegisterResponse)
	copier.Copy(userRegisterResponse, user)

	userRegisterResponse.TanggalLahir = user.TanggalLahir.Format(model.TANGGAL_LAHIR_DATE_FORMAT)

	province, err := u.provinceRepository.FindByID(ctx, user.IdProvinsi)
	if err != nil {
		return nil, err
	}
	userRegisterResponse.IdProvinsi = province

	city, err := u.cityRepository.FindByID(ctx, user.IdKota[0:2], user.IdKota)
	if err != nil {
		return nil, err
	}
	userRegisterResponse.IdKota = city

	toko := new(model.Toko)
	toko.IdUser = user.ID
	toko, err = u.tokoRepository.Create(ctx, toko)
	if err != nil {
		return nil, err
	}
	tokoCreateResponse := new(model.TokoCreateResponse)
	copier.Copy(tokoCreateResponse, toko)
	userRegisterResponse.Toko = tokoCreateResponse

	return userRegisterResponse, nil
}

func (u *userUsecase) LoginUser(ctx context.Context, req *model.UserLoginRequest) (*model.UserLoginResponse, error) {
	user, err := u.userRepository.FindByNoTelp(ctx, req.NoTelp)
	if err != nil {
		return nil, err
	}
	hashedPassword := user.KataSandi
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.KataSandi))
	if err != nil {
		return nil, errors.New("no telp atau kata sandi salah")
	}
	userLoginResponse := new(model.UserLoginResponse)
	copier.Copy(userLoginResponse, user)

	// Create the claims
	idString := strconv.Itoa(user.ID)
	claims := jwt.MapClaims{
		"idString": idString,
		"isAdmin":  user.IsAdmin,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return nil, err
	}
	userLoginResponse.Token = signedToken

	userLoginResponse.TanggalLahir = user.TanggalLahir.Format(model.TANGGAL_LAHIR_DATE_FORMAT)

	province, err := u.provinceRepository.FindByID(ctx, user.IdProvinsi)
	if err != nil {
		return nil, err
	}
	userLoginResponse.IdProvinsi = province

	city, err := u.cityRepository.FindByID(ctx, user.IdKota[0:2], user.IdKota)
	if err != nil {
		return nil, err
	}
	userLoginResponse.IdKota = city

	return userLoginResponse, nil
}

func (u *userUsecase) GetCurrentUser(ctx context.Context, userId int) (*model.UserResponse, error) {
	user, err := u.userRepository.FindByID(ctx, userId)
	if err != nil {
		return nil, err
	}
	userResponse := new(model.UserResponse)
	copier.Copy(userResponse, user)

	userResponse.TanggalLahir = user.TanggalLahir.Format(model.TANGGAL_LAHIR_DATE_FORMAT)

	province, err := u.provinceRepository.FindByID(ctx, user.IdProvinsi)
	if err != nil {
		return nil, err
	}
	userResponse.IdProvinsi = province

	city, err := u.cityRepository.FindByID(ctx, user.IdKota[0:2], user.IdKota)
	if err != nil {
		return nil, err
	}
	userResponse.IdKota = city

	return userResponse, nil
}

func (u *userUsecase) EditCurrentUser(ctx context.Context, userId int, req *model.UserUpdateRequest) (*model.UserResponse, error) {
	user := new(model.User)

	tanggalLahir, err := time.Parse(model.TANGGAL_LAHIR_DATE_FORMAT, req.TanggalLahir)
	if err != nil {
		return nil, err
	}
	user.TanggalLahir = tanggalLahir

	_, err = u.provinceRepository.FindByID(ctx, req.IdProvinsi)
	if err != nil {
		return nil, err
	}

	if req.IdKota[0:2] != req.IdProvinsi {
		return nil, errors.New("2 angka awal id kota harus sama dengan id provinsi")
	}
	_, err = u.cityRepository.FindByID(ctx, req.IdKota[0:2], req.IdKota)
	if err != nil {
		return nil, err
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(req.KataSandi), 14)
	if err != nil {
		return nil, err
	}
	hashedPassword := string(bytes)
	req.KataSandi = hashedPassword

	copier.Copy(user, req)
	user, err = u.userRepository.UpdateByID(ctx, userId, user)
	if err != nil {
		return nil, err
	}
	userResponse := new(model.UserResponse)
	copier.Copy(userResponse, user)

	userResponse.TanggalLahir = user.TanggalLahir.Format(model.TANGGAL_LAHIR_DATE_FORMAT)

	province, err := u.provinceRepository.FindByID(ctx, user.IdProvinsi)
	if err != nil {
		return nil, err
	}
	userResponse.IdProvinsi = province

	city, err := u.cityRepository.FindByID(ctx, user.IdKota[0:2], user.IdKota)
	if err != nil {
		return nil, err
	}
	userResponse.IdKota = city

	return userResponse, nil
}

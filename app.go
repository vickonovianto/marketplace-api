package main

import (
	"errors"
	"fmt"
	"log"
	"marketplace-api/config"
	"marketplace-api/delivery"
	"marketplace-api/helper"
	"marketplace-api/repository"
	"marketplace-api/usecase"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v3"
)

type (
	server struct {
		httpServer *fiber.App
		cfg        config.Config
	}

	Server interface {
		Run()
	}
)

func InitServer(cfg config.Config) Server {
	app := fiber.New()

	// Middleware
	app.Use(logger.New())

	// Serving uploads folder
	app.Static("/uploads", "./uploads")

	return &server{
		httpServer: app,
		cfg:        cfg,
	}
}

func (s *server) Run() {
	api := s.httpServer.Group(os.Getenv("API_PREFIX"))

	rootFolderPath, err := filepath.Abs("./")
	if err != nil {
		log.Fatal(err)
	}

	uploadFolderPath := filepath.Join(rootFolderPath, "uploads")
	if _, err := os.Stat(uploadFolderPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(uploadFolderPath, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	tokoFolderPath := filepath.Join(uploadFolderPath, "toko")
	if _, err := os.Stat(tokoFolderPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(tokoFolderPath, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	produkFolderPath := filepath.Join(uploadFolderPath, "produk")
	if _, err := os.Stat(produkFolderPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(produkFolderPath, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	jwtMiddleware := jwtware.New(
		jwtware.Config{
			SigningKey:   []byte(os.Getenv("SECRET")),
			ErrorHandler: helper.JwtMiddlewareErrorHandler,
		},
	)

	provinceRepository := repository.NewProvinceRepository(s.cfg)
	provinceUsecase := usecase.NewProvinceUsecase(provinceRepository)
	provinceDelivery := delivery.NewProvinceDelivery(provinceUsecase)

	cityRepository := repository.NewCityRepository(s.cfg)
	cityUsecase := usecase.NewCityUsecase(cityRepository)
	cityDelivery := delivery.NewCityDelivery(cityUsecase)

	provinceCityGroup := api.Group("/provcity")
	provinceDelivery.MountUnprotectedRoutes(provinceCityGroup)
	cityDelivery.MountUnprotectedRoutes(provinceCityGroup)

	tokoRepository := repository.NewTokoRepository(s.cfg)
	tokoUsecase := usecase.NewTokoUsecase(tokoRepository)
	tokoDelivery := delivery.NewTokoDelivery(tokoUsecase)
	tokoGroup := api.Group("/toko")
	tokoDelivery.MountProtectedRoutes(jwtMiddleware, tokoGroup)

	userRepository := repository.NewUserRepository(s.cfg)
	userUsecase := usecase.NewUserUsecase(userRepository, tokoRepository, provinceRepository, cityRepository)
	userDelivery := delivery.NewUserDelivery(userUsecase)
	userGroup := api.Group("/user")
	userDelivery.MountProtectedRoutes(jwtMiddleware, userGroup)

	alamatRepository := repository.NewAlamatRepository(s.cfg)
	alamatUsecase := usecase.NewAlamatUsecase(alamatRepository)
	alamatDelivery := delivery.NewAlamatDelivery(alamatUsecase)
	alamatGroup := userGroup.Group("/alamat")
	alamatDelivery.MountProtectedRoutes(jwtMiddleware, alamatGroup)

	authDelivery := delivery.NewAuthDelivery(userUsecase)
	authGroup := api.Group("/auth")
	authDelivery.MountUnprotectedRoutes(authGroup)

	categoryRepository := repository.NewCategoryRepository(s.cfg)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepository)
	categoryDelivery := delivery.NewCategoryDelivery(categoryUsecase)
	categoryGroup := api.Group("/category")
	categoryDelivery.MountUnprotectedRoutes(categoryGroup)
	categoryDelivery.MountProtectedRoutes(jwtMiddleware, categoryGroup)

	fotoProdukRepository := repository.NewFotoProdukRepository(s.cfg)
	produkRepository := repository.NewProdukRepository(s.cfg)
	produkUsecase := usecase.NewProdukUsecase(
		produkRepository,
		fotoProdukRepository,
		tokoRepository,
		categoryRepository,
	)
	produkDelivery := delivery.NewProdukDelivery(produkUsecase)
	produkGroup := api.Group("/product")
	produkDelivery.MountUnprotectedRoutes(produkGroup)
	produkDelivery.MountProtectedRoutes(jwtMiddleware, produkGroup)

	logProdukRepository := repository.NewLogProdukRepository(s.cfg)

	detailTrxRepository := repository.NewDetailTrxRepository(s.cfg)

	trxRepository := repository.NewTrxRepository(s.cfg)
	trxUsecase := usecase.NewTrxUsecase(
		trxRepository,
		alamatRepository,
		detailTrxRepository,
		logProdukRepository,
		tokoRepository,
		categoryRepository,
		fotoProdukRepository,
		produkRepository,
	)
	trxDelivery := delivery.NewTrxDelivery(trxUsecase)
	trxGroup := api.Group("/trx")
	trxDelivery.MountProtectedRoutes(jwtMiddleware, trxGroup)

	if err := s.httpServer.Listen(fmt.Sprintf(":%d", s.cfg.ServicePort())); err != nil {
		log.Panic(err)
	}
}

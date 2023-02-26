package mysql

import (
	"log"
	"marketplace-api/model"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitGorm() *gorm.DB {
	connection := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(mysql.Open(connection))
	if err != nil {
		log.Panic(err)
	}
	db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Toko{},
		&model.Alamat{},
		&model.Produk{},
		&model.FotoProduk{},
		&model.LogProduk{},
		&model.DetailTrx{},
		&model.Trx{},
	)
	return db
}

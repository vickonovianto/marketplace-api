package config

import (
	"marketplace-api/config/mysql"
	"os"
	"strconv"

	"gorm.io/gorm"
)

type (
	config struct {
	}

	Config interface {
		ServicePort() int
		Database() *gorm.DB
	}
)

func NewConfig() Config {
	return &config{}
}

func (c *config) Database() *gorm.DB {
	return mysql.InitGorm()
}

func (c *config) ServicePort() int {
	v := os.Getenv("PORT")
	port, _ := strconv.Atoi(v)
	return port
}

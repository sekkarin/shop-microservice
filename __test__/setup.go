package whydoweneedtest

import (
	"github.com/sekkarin/shop-microservice/config"
)

func NewTestConfig() *config.Config {
	cfg := config.LoadConfig("../env/test/.env")
	return &cfg
}

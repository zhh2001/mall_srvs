package global

import (
	"gorm.io/gorm"

	"mall_srvs/order_srv/config"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig
)

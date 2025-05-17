package global

import (
	"gorm.io/gorm"

	"mall_srvs/inventory_srv/config"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig
)

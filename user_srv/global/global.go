package global

import (
	"gorm.io/gorm"

	"mall_srvs/user_srv/config"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig
)

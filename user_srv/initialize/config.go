package initialize

import (
	"fmt"

	"github.com/spf13/viper"

	"mall_srvs/user_srv/global"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	// 从配置文件中读取出对应的配置
	debug := GetEnvInfo("MALL_DEBUG")
	configFilePrefix := "config"
	configFilename := fmt.Sprintf("./user_srv/%s-pro.yaml", configFilePrefix)
	if debug {
		configFilename = fmt.Sprintf("./user_srv/%s-debug.yaml", configFilePrefix)
	}

	v := viper.New()
	v.SetConfigFile(configFilename)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		panic(err)
	}
}

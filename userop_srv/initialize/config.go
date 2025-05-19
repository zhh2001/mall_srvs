package initialize

import (
	"encoding/json"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"mall_srvs/userop_srv/global"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	// 从配置文件中读取出对应的配置
	debug := GetEnvInfo("MALL_DEBUG")
	configFilePrefix := "config"
	configFilename := fmt.Sprintf("./userop_srv/%s-pro.yaml", configFilePrefix)
	if debug {
		configFilename = fmt.Sprintf("./userop_srv/%s-debug.yaml", configFilePrefix)
	}

	v := viper.New()
	v.SetConfigFile(configFilename)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(&global.NacosConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息：%v", global.NacosConfig)

	// 从 Nacos 中读取配置信息
	sc := []constant.ServerConfig{
		{
			IpAddr:      global.NacosConfig.Host,
			Port:        global.NacosConfig.Port,
			ContextPath: "/nacos",
			Scheme:      "http",
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})
	if err != nil {
		zap.S().Fatalf("读取 Nacos 配置失败：%s", err.Error())
	}
	fmt.Println(content)
	err = json.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		panic(err)
	}
	fmt.Println(global.ServerConfig)
}

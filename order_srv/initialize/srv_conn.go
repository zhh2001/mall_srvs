package initialize

import (
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"mall_srvs/order_srv/global"
	"mall_srvs/order_srv/proto"
)

func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo
	goodsConn, err := grpc.NewClient(
		fmt.Sprintf(
			"consul://%s:%d/%s?wait=14s",
			consulInfo.Host,
			consulInfo.Port,
			global.ServerConfig.GoodsSrvInfo.Name,
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接【商品服务】失败")
		return
	}

	global.GoodsSrvClient = proto.NewGoodsClient(goodsConn)

	// 初始化库存服务连接
	invConn, err := grpc.NewClient(
		fmt.Sprintf(
			"consul://%s:%d/%s?wait=14s",
			consulInfo.Host,
			consulInfo.Port,
			global.ServerConfig.InventorySrvInfo.Name,
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接【库存服务】失败")
		return
	}

	global.InventorySrvClient = proto.NewInventoryClient(invConn)
}

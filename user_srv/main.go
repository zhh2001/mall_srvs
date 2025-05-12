package main

import (
	"flag"
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"mall_srvs/user_srv/handler"
	"mall_srvs/user_srv/initialize"
	"mall_srvs/user_srv/proto"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "IPv4地址")
	Port := flag.Int("port", 50051, "端口号")

	// 初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()

	flag.Parse()
	zap.S().Info("IP:", *IP)
	zap.S().Info("Port:", *Port)

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen: " + err.Error())
	}
	err = server.Serve(lis)
	if err != nil {
		panic("failed to serve: " + err.Error())
	}
}

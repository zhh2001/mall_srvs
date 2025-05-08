package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"net"

	"mall_srvs/user_srv/handler"
	"mall_srvs/user_srv/proto"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "IPv4地址")
	Port := flag.Int("port", 50051, "端口号")

	flag.Parse()
	fmt.Println("IP:", *IP)
	fmt.Println("Port:", *Port)

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

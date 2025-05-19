package main

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"mall_srvs/userop_srv/proto"
)

var (
	userFavClient proto.UserFavClient
	messageClient proto.MessageClient
	addressClient proto.AddressClient
	conn          *grpc.ClientConn
)

func TestAddressList() {
	_, err := addressClient.GetAddressList(context.Background(), &proto.AddressRequest{
		UserId: 1,
	})
	if err != nil {
		panic(err)
	}
}

func TestMessageList() {
	_, err := messageClient.MessageList(context.Background(), &proto.MessageRequest{
		UserId: 1,
	})
	if err != nil {
		panic(err)
	}
}

func TestUserFav() {
	_, err := userFavClient.GetFavList(context.Background(), &proto.UserFavRequest{
		UserId: 1,
	})
	if err != nil {
		panic(err)
	}
}

func Init() {
	var err error
	conn, err = grpc.NewClient("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	userFavClient = proto.NewUserFavClient(conn)
	messageClient = proto.NewMessageClient(conn)
	addressClient = proto.NewAddressClient(conn)
}

func main() {
	Init()

	TestAddressList()
	TestMessageList()
	TestUserFav()

	err := conn.Close()
	if err != nil {
		panic(err)
	}
}

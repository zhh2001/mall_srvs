package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"mall_srvs/goods_srv/proto"
)

var (
	brandClient proto.GoodsClient
	conn        *grpc.ClientConn
)

func TestGetCategoryBrandList() {
	rsp, err := brandClient.GetCategoryBrandList(context.Background(), &proto.CategoryInfoRequest{
		Id: 135475,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	fmt.Println(rsp.Data)
}

func Init() {
	var err error
	conn, err = grpc.NewClient("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	brandClient = proto.NewGoodsClient(conn)
}

func main() {
	Init()
	defer func() {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}()

	TestGetCategoryBrandList()
}

package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"mall_srvs/goods_srv/proto"
)

var (
	brandClient proto.GoodsClient
	conn        *grpc.ClientConn
)

func TestGetBrandList() {
	rsp, err := brandClient.BrandList(context.Background(), &proto.BrandFilterRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, brand := range rsp.Data {
		fmt.Println(brand.Name)
	}
}

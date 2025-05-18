package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"mall_srvs/order_srv/proto"
)

var (
	orderClient proto.OrderClient
	conn        *grpc.ClientConn
)

func TestCreateCartItem(userId int32, nums int32, goodsId int32) {
	rsp, err := orderClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  userId,
		Nums:    nums,
		GoodsId: goodsId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.GetId())
}

func TestCartItemList(userId int32) {
	rsp, err := orderClient.CartItemList(context.Background(), &proto.UserInfo{
		Id: userId,
	})
	if err != nil {
		panic(err)
	}
	for _, item := range rsp.Data {
		fmt.Println(item.GetId(), item.GetGoodsId(), item.GetNums())
	}
}

func TestUpdateCartItem(id int32, userId int32, goodsId int32) {
	_, err := orderClient.UpdateCartItem(context.Background(), &proto.CartItemRequest{
		Id:      id,
		UserId:  userId,
		GoodsId: goodsId,
		Checked: true,
	})
	if err != nil {
		panic(err)
	}
}

func TestCreateOrder(userId int32) {
	_, err := orderClient.CreateOrder(context.Background(), &proto.OrderRequest{
		UserId:  userId,
		Address: "上海市",
		Name:    "Zhang",
		Mobile:  "13877779999",
		Post:    "请尽快发货",
	})
	if err != nil {
		panic(err)
	}
}

func TestGetOrderDetail(orderId int32) {
	rsp, err := orderClient.OrderDetail(context.Background(), &proto.OrderRequest{
		Id: orderId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.OrderInfo.OrderSn)
	for _, goods := range rsp.GetGoods() {
		fmt.Println(goods.GetGoodsName())
	}
}

func TestGetOrderList() {
	rsp, err := orderClient.OrderList(context.Background(), &proto.OrderFilterRequest{})
	if err != nil {
		panic(err)
	}
	for _, order := range rsp.GetData() {
		fmt.Println(order.GetOrderSn())
	}
}

func Init() {
	var err error
	conn, err = grpc.NewClient("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	orderClient = proto.NewOrderClient(conn)
}

func main() {
	Init()

	//TestCreateCartItem(1, 1, 422)
	//TestCartItemList(1)
	//TestUpdateCartItem(1, 1, 421)
	//TestCreateOrder(1)
	TestGetOrderDetail(3)
	TestGetOrderList()

	err := conn.Close()
	if err != nil {
		panic(err)
	}
}

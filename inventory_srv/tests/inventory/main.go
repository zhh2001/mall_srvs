package main

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"mall_srvs/inventory_srv/proto"
)

var (
	invClient proto.InventoryClient
	conn      *grpc.ClientConn
)

func TestSetInv(goodsId int32, num int32) {
	_, err := invClient.SetInv(context.Background(), &proto.GoodsInvInfo{
		GoodsId: goodsId,
		Num:     num,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("设置库存成功")
}

func TestInvDetail(goodsId int32) {
	rsp, err := invClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: goodsId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Num)
}

func TestSell() {
	/*
		1. 第一件扣减成功；第二件扣减失败
		2. 两件都扣减成功
	*/
	_, err := invClient.Sell(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{GoodsId: 421, Num: 10},
			{GoodsId: 422, Num: 30},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("库存扣减成功")
}

func TestConSell(wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := invClient.Sell(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{GoodsId: 421, Num: 1},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("库存扣减成功")
}

func TestReback() {
	_, err := invClient.Reback(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{GoodsId: 421, Num: 10},
			{GoodsId: 422, Num: 30},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("库存归还成功")
}

func Init() {
	var err error
	conn, err = grpc.NewClient("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	invClient = proto.NewInventoryClient(conn)
}

func Test() {
	TestSetInv(421, 100)
	TestSetInv(422, 40)
	TestInvDetail(421)
	TestSell()
	TestReback()
}

func Reset() {
	// 为所有商品设置 100 库存
	var i int32
	for i = 421; i <= 840; i++ {
		TestSetInv(i, 100)
	}
}

func ConSell() {
	TestSetInv(421, 50)

	var wg sync.WaitGroup
	for i := 0; i < 80; i++ {
		wg.Add(1)
		go TestConSell(&wg)
	}

	wg.Wait()
}

func main() {
	Init()

	//Test()

	//Reset()

	ConSell()

	err := conn.Close()
	if err != nil {
		panic(err)
	}
}

package handler

import (
	"mall_srvs/goods_srv/proto"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

//func (g *GoodsServer) GoodsList(context.Context, *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {}
//BatchGetGoods(context.Context, *BatchGoodsIdInfo) (*GoodsListResponse, error)
//CreateGoods(context.Context, *CreateGoodsInfo) (*GoodsInfoResponse, error)
//DeleteGoods(context.Context, *DeleteGoodsInfo) (*emptypb.Empty, error)
//UpdateGoods(context.Context, *CreateGoodsInfo) (*emptypb.Empty, error)
//GetGoodsDetail(context.Context, *GoodInfoRequest) (*GoodsInfoResponse, error)

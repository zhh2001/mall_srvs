package handler

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"mall_srvs/userop_srv/global"
	"mall_srvs/userop_srv/model"
	"mall_srvs/userop_srv/proto"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

func (s *InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	// 设置库存
	var inv model.Inventory
	global.DB.Where(&model.Inventory{Goods: req.GetGoodsId()}).First(&inv)
	if inv.Goods == 0 {
		inv.Goods = req.GetGoodsId()
	}
	inv.Stocks = req.GetNum()

	global.DB.Save(&inv)
	return &emptypb.Empty{}, nil
}

func (s *InventoryServer) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if result := global.DB.Where(&model.Inventory{Goods: req.GetGoodsId()}).First(&inv); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有库存信息")
	}
	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}

var m sync.Mutex

func (s *InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// 扣减库存，本地事务，数据一致性
	//m.Lock() // 获取锁

	client := goredislib.NewClient(&goredislib.Options{
		Addr: "10.120.21.77:6379",
	})
	pool := goredis.NewPool(client)
	rs := redsync.New(pool)

	tx := global.DB.Begin()
	for _, goodsInfo := range req.GoodsInfo {
		//for {
		var inv model.Inventory
		//if result := tx.Clauses(clause.Locking{
		//	Strength: clause.LockingStrengthUpdate,
		//}).Where(&model.Inventory{
		//	Goods: goodsInfo.GetGoodsId(),
		//}).First(&inv); result.RowsAffected == 0 {
		//	tx.Rollback()
		//	return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		//}

		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodsInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}

		if result := global.DB.Where(&model.Inventory{
			Goods: goodsInfo.GetGoodsId(),
		}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		// 判断库存是否充足
		if inv.Stocks < goodsInfo.Num {
			tx.Rollback()
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		// 扣减
		inv.Stocks = inv.Stocks - goodsInfo.Num
		tx.Save(&inv)

		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
		}

		//if result := tx.Model(&model.Inventory{}).
		//	Select("stocks", "version").
		//	Where("goods = ? AND version = ?", goodsInfo.GoodsId, inv.Version).
		//	Updates(model.Inventory{
		//		Stocks:  inv.Stocks,
		//		Version: inv.Version + 1,
		//	}); result.RowsAffected == 0 {
		//	zap.S().Infof("库存扣减失败")
		//} else {
		//	break
		//}
		//tx.Save(&inv)
		//}
	}
	tx.Commit()
	//m.Unlock() // 释放锁
	return &emptypb.Empty{}, nil
}

func (s *InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// 库存归还：1. 订单超时归还；2. 订单创建失败；3. 手动归还
	tx := global.DB.Begin()
	for _, goodsInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := global.DB.Where(&model.Inventory{Goods: goodsInfo.GetGoodsId()}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}

		inv.Stocks = inv.Stocks + goodsInfo.Num
		tx.Save(&inv)
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}

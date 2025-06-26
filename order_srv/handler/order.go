package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"mall_srvs/order_srv/global"
	"mall_srvs/order_srv/model"
	"mall_srvs/order_srv/proto"
)

type OrderServer struct {
	proto.UnimplementedOrderServer
}

func GenerateOrderSn(userId int32) string {
	// 订单号的生成规则
	/*
		年月日时分秒 + 用户ID + 2位随机数
	*/
	now := time.Now()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(),
		userId, r.Intn(90)+10,
	)
	return orderSn
}

func (orderServer *OrderServer) CartItemList(ctx context.Context, req *proto.UserInfo) (*proto.CartItemListResponse, error) {
	// 获取用户的购物车列表
	var shopCarts []model.ShoppingCart
	var rsp proto.CartItemListResponse

	if result := global.DB.Where(&model.ShoppingCart{
		User: req.GetId(),
	}).Find(&shopCarts); result.RowsAffected == 0 {
		return nil, result.Error
	} else {
		rsp.Total = int32(result.RowsAffected)
	}

	for _, shopCart := range shopCarts {
		rsp.Data = append(rsp.Data, &proto.ShopCartInfoResponse{
			Id:      shopCart.ID,
			UserId:  shopCart.User,
			GoodsId: shopCart.Goods,
			Nums:    shopCart.Nums,
			Checked: shopCart.Checked,
		})
	}

	return &rsp, nil
}

func (orderServer *OrderServer) CreateCartItem(ctx context.Context, req *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	// 将商品添加到购物车 1. 购物车中原本没有这件商品 - 新建一个记录 2. 这个商品之前添加到了购物车 - 合并
	var shopCart model.ShoppingCart

	if result := global.DB.Where(&model.ShoppingCart{
		Goods: req.GoodsId,
		User:  req.UserId,
	}).First(&shopCart); result.RowsAffected == 1 {
		// 如果记录已经存在，则合并购物车记录, 更新操作
		shopCart.Nums = shopCart.Nums + req.Nums
	} else {
		// 插入操作
		shopCart.User = req.UserId
		shopCart.Goods = req.GoodsId
		shopCart.Nums = req.Nums
		shopCart.Checked = false
	}

	global.DB.Save(&shopCart)
	return &proto.ShopCartInfoResponse{Id: shopCart.ID}, nil
}

func (orderServer *OrderServer) UpdateCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	// 更新购物车记录，更新数量和选中状态
	var shopCart model.ShoppingCart

	if result := global.DB.
		Where("goods = ? AND user = ?", req.GoodsId, req.UserId).
		First(&shopCart); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}

	shopCart.Checked = req.Checked
	if req.Nums > 0 {
		shopCart.Nums = req.Nums
	}
	global.DB.Save(&shopCart)

	return &emptypb.Empty{}, nil
}

func (orderServer *OrderServer) DeleteCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	if result := global.DB.
		Where("goods = ? AND user = ?", req.GoodsId, req.UserId).
		Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}
	return &emptypb.Empty{}, nil
}

func (orderServer *OrderServer) OrderList(ctx context.Context, req *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	var orders []model.OrderInfo
	var rsp proto.OrderListResponse

	var total int64
	global.DB.Where(&model.OrderInfo{
		User: req.UserId,
	}).Count(&total)
	rsp.Total = int32(total)

	// 分页
	global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).
		Where(&model.OrderInfo{User: req.UserId}).Find(&orders)
	for _, order := range orders {
		orderInfo := proto.OrderInfoResponse{
			Id:      order.ID,
			UserId:  order.User,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status:  order.Status,
			Post:    order.Post,
			Total:   order.OrderMount,
			Address: order.Address,
			Name:    order.SignerName,
			Mobile:  order.SingerMobile,
			AddTime: order.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		rsp.Data = append(rsp.Data, &orderInfo)
	}
	return &rsp, nil
}

func (orderServer *OrderServer) OrderDetail(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	var order model.OrderInfo
	var rsp proto.OrderInfoDetailResponse

	if result := global.DB.Where(&model.OrderInfo{
		BaseModel: model.BaseModel{
			ID: req.Id,
		},
		User: req.UserId,
	}).First(&order); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}

	var orderInfo proto.OrderInfoResponse
	orderInfo.Id = order.ID
	orderInfo.UserId = order.User
	orderInfo.OrderSn = order.OrderSn
	orderInfo.PayType = order.PayType
	orderInfo.Status = order.Status
	orderInfo.Post = order.Post
	orderInfo.Total = order.OrderMount
	orderInfo.Address = order.Address
	orderInfo.Name = order.SignerName
	orderInfo.Mobile = order.SingerMobile

	rsp.OrderInfo = &orderInfo

	var orderGoods []model.OrderGoods
	if result := global.DB.Where(&model.OrderGoods{
		Order: order.ID,
	}).Find(&orderGoods); result.Error != nil {
		return nil, result.Error
	}

	for _, orderGood := range orderGoods {
		orderItem := proto.OrderItemResponse{
			GoodsId:    orderGood.Goods,
			GoodsName:  orderGood.GoodsName,
			GoodsPrice: orderGood.GoodsPrice,
			GoodsImage: orderGood.GoodsImage,
			Nums:       orderGood.Nums,
		}
		rsp.Goods = append(rsp.Goods, &orderItem)
	}

	return &rsp, nil
}

type OrderListener struct {
	Code        codes.Code
	Detail      string
	ID          int32
	OrderAmount float32
}

func (orderListener *OrderListener) ExecuteLocalTransaction(message *primitive.Message) primitive.LocalTransactionState {
	var orderInfo model.OrderInfo
	err := json.Unmarshal(message.Body, &orderInfo)

	var goodsIds []int32
	var shopCarts []model.ShoppingCart
	goodsNumsMap := make(map[int32]int32)
	if result := global.DB.Where(&model.ShoppingCart{
		User:    orderInfo.User,
		Checked: true,
	}).Find(&shopCarts); result.RowsAffected == 0 {
		orderListener.Code = codes.InvalidArgument
		orderListener.Detail = "没有选择结算的商品"
		return primitive.RollbackMessageState
	}

	for _, shopCart := range shopCarts {
		goodsIds = append(goodsIds, shopCart.Goods)
		goodsNumsMap[shopCart.Goods] = shopCart.Nums
	}

	// 跨服务调用 - 商品微服务
	goods, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: goodsIds,
	})
	if err != nil {
		orderListener.Code = codes.Internal
		orderListener.Detail = "批量查询商品信息失败"
		return primitive.RollbackMessageState
	}

	var orderAmount float32
	var orderGoods []*model.OrderGoods
	var goodsInvInfo []*proto.GoodsInvInfo
	for _, good := range goods.Data {
		orderAmount = orderAmount + good.GetShopPrice()*float32(goodsNumsMap[good.GetId()])
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      good.GetId(),
			GoodsName:  good.GetName(),
			GoodsImage: good.GetGoodsFrontImage(),
			GoodsPrice: good.GetShopPrice(),
			Nums:       goodsNumsMap[good.GetId()],
		})

		goodsInvInfo = append(goodsInvInfo, &proto.GoodsInvInfo{
			GoodsId: good.GetId(),
			Num:     goodsNumsMap[good.GetId()],
		})
	}

	// 跨服务调用 - 库存微服务
	if _, err = global.InventorySrvClient.Sell(context.Background(), &proto.SellInfo{
		GoodsInfo: goodsInvInfo,
		OrderSn:   orderInfo.OrderSn,
	}); err != nil {
		orderListener.Code = codes.ResourceExhausted
		orderListener.Detail = "扣减库存失败"
		return primitive.RollbackMessageState
	}

	// 生成订单表
	// 20250518xxxx
	tx := global.DB.Begin()
	orderInfo.OrderMount = orderAmount
	if result := tx.Save(&orderInfo); result.RowsAffected == 0 {
		tx.Rollback()
		orderListener.Code = codes.Internal
		orderListener.Detail = "创建订单失败"
		return primitive.CommitMessageState
	}

	orderListener.OrderAmount = orderAmount
	orderListener.ID = orderInfo.ID
	for _, orderGood := range orderGoods {
		orderGood.Order = orderInfo.ID
	}

	// 批量插入 orderGoods
	if result := tx.CreateInBatches(orderGoods, 100); result.RowsAffected == 0 {
		tx.Rollback()
		orderListener.Code = codes.Internal
		orderListener.Detail = "批量插入订单商品失败"
		return primitive.CommitMessageState
	}

	if result := tx.Where(&model.ShoppingCart{
		User:    orderInfo.User,
		Checked: true,
	}).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		tx.Rollback()
		orderListener.Code = codes.Internal
		orderListener.Detail = "删除购物车记录失败"
		return primitive.CommitMessageState
	}

	tx.Commit()
	orderListener.Code = codes.OK
	return primitive.UnknowState
}

func (orderListener *OrderListener) CheckLocalTransaction(message *primitive.MessageExt) primitive.LocalTransactionState {
	var orderInfo model.OrderInfo
	_ = json.Unmarshal(message.Body, &orderInfo)

	if result := global.DB.Where(model.OrderInfo{
		OrderSn: orderInfo.OrderSn,
	}).First(&orderInfo); result.RowsAffected == 0 {
		return primitive.CommitMessageState
	}

	return primitive.RollbackMessageState
}

func (orderServer *OrderServer) CreateOrder(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	/*
		新建订单
			1. 从购物车中获取到选中的商品
			2. 商品的价格自己查询 - 访问商品服务 (跨微服务)
			3. 库存的扣减 - 访问库存服务 (跨微服务)
			4. 订单的基本信息表 - 订单的商品信息表
			5. 从购物车中删除已购买的记录
	*/
	orderListener := OrderListener{}
	p, err := rocketmq.NewTransactionProducer(&orderListener, producer.WithNameServer([]string{"10.120.21.77:9876"}))
	if err != nil {
		zap.S().Errorf("生成producer失败：%s", err.Error())
		return nil, err
	}
	if err = p.Start(); err != nil {
		zap.S().Errorf("启动producer失败：%s", err.Error())
		return nil, err
	}

	order := model.OrderInfo{
		OrderSn:      GenerateOrderSn(req.GetUserId()),
		Address:      req.GetAddress(),
		SignerName:   req.GetName(),
		SingerMobile: req.GetMobile(),
		Post:         req.GetPost(),
		User:         req.GetUserId(),
	}

	jsonString, _ := json.Marshal(order)

	_, err = p.SendMessageInTransaction(
		context.Background(),
		primitive.NewMessage("order_reback", jsonString),
	)
	if err != nil {
		fmt.Printf("发送失败：%s\n", err)
		return nil, status.Errorf(codes.Internal, "发送消息失败")
	}
	if orderListener.Code != codes.OK {
		return nil, status.Errorf(orderListener.Code, orderListener.Detail)
	}

	return &proto.OrderInfoResponse{
		Id:      orderListener.ID,
		OrderSn: order.OrderSn,
		Total:   orderListener.OrderAmount,
	}, nil
}

func (orderServer *OrderServer) UpdateOrderStatus(ctx context.Context, req *proto.OrderStatus) (*emptypb.Empty, error) {
	if result := global.DB.Model(&model.OrderInfo{}).Where("order_sn = ?", req.GetOrderSn()).Update("status", req.GetStatus()); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}
	return &emptypb.Empty{}, nil
}

package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"mall_srvs/userop_srv/global"
	"mall_srvs/userop_srv/model"
	"mall_srvs/userop_srv/proto"
)

func (userOpServer *UserOpServer) GetAddressList(ctx context.Context, req *proto.AddressRequest) (*proto.AddressListResponse, error) {
	var addresses []model.Address
	var rsp proto.AddressListResponse
	var addressResponse []*proto.AddressResponse

	if result := global.DB.Where(&model.Address{User: req.GetUserId()}).Find(&addresses); result.RowsAffected != 0 {
		rsp.Total = int32(result.RowsAffected)
	}

	for _, address := range addresses {
		addressResponse = append(addressResponse, &proto.AddressResponse{
			Id:           address.ID,
			UserId:       address.User,
			Province:     address.Province,
			City:         address.City,
			District:     address.District,
			Address:      address.Address,
			SignerName:   address.SignerName,
			SignerMobile: address.SignerMobile,
		})
	}
	rsp.Data = addressResponse

	return &rsp, nil
}

func (userOpServer *UserOpServer) CreateAddress(ctx context.Context, req *proto.AddressRequest) (*proto.AddressResponse, error) {
	var address model.Address

	address.User = req.GetUserId()
	address.Province = req.GetProvince()
	address.City = req.GetCity()
	address.District = req.GetDistrict()
	address.Address = req.GetAddress()
	address.SignerName = req.GetSignerName()
	address.SignerMobile = req.GetSignerMobile()

	global.DB.Save(&address)

	return &proto.AddressResponse{Id: address.ID}, nil
}

func (userOpServer *UserOpServer) DeleteAddress(ctx context.Context, req *proto.AddressRequest) (*emptypb.Empty, error) {
	if result := global.DB.Where(
		"id = ? AND user = ?",
		req.GetId(),
		req.GetUserId(),
	).Delete(&model.Address{}); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "收货地址不存在")
	}
	return &emptypb.Empty{}, nil
}

func (userOpServer *UserOpServer) UpdateAddress(ctx context.Context, req *proto.AddressRequest) (*emptypb.Empty, error) {
	var address model.Address

	if result := global.DB.Where(
		"id = ? AND user = ?",
		req.GetId(),
		req.GetUserId(),
	).First(&address); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}

	if address.Province != "" {
		address.Province = req.GetProvince()
	}

	if address.City != "" {
		address.City = req.GetCity()
	}

	if address.District != "" {
		address.District = req.GetDistrict()
	}

	if address.Address != "" {
		address.Address = req.GetAddress()
	}

	if address.SignerName != "" {
		address.SignerName = req.GetSignerName()
	}

	if address.SignerMobile != "" {
		address.SignerMobile = req.GetSignerMobile()
	}

	global.DB.Save(&address)

	return &emptypb.Empty{}, nil
}

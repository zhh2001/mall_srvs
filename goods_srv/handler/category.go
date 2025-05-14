/* 商品分类 */

package handler

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"mall_srvs/goods_srv/global"
	"mall_srvs/goods_srv/model"
	"mall_srvs/goods_srv/proto"
)

func (s *GoodsServer) GetAllCategoryList(ctx context.Context, req *emptypb.Empty) (*proto.CategoryListResponse, error) {
	/*
		[
			{
				"id":     xxx,
				"name":   "",
				"level":  1,
				"is_tab": false,
				"parent": xxx,
				"sub_category": [
					"id":     xxx,
					"name":   "",
					"level":  1,
					"is_tab": false,
					"sub_category": []
				]
			}
		]
	*/
	var categories []model.Category
	global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categories)
	b, _ := json.Marshal(&categories)
	return &proto.CategoryListResponse{
		JsonData: string(b),
	}, nil
}

func (s *GoodsServer) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	categoryListResponse := proto.SubCategoryListResponse{}
	var category model.Category
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "商品分类不存在")
	}

	categoryListResponse.Info = &proto.CategoryInfoResponse{
		Id:             category.ID,
		Name:           category.Name,
		Level:          category.Level,
		IsTab:          category.IsTab,
		ParentCategory: category.ParentCategoryID,
	}

	var subCategories []model.Category
	var subCategoryResponse []*proto.CategoryInfoResponse
	var preloads string
	if category.Level == 1 {
		preloads = "SubCategory.SubCategory"
	} else {
		preloads = "SubCategory"
	}
	global.DB.Where(&model.Category{ParentCategoryID: req.Id}).Preload(preloads).Find(&subCategories)
	for _, subCategory := range subCategories {
		subCategoryResponse = append(subCategoryResponse, &proto.CategoryInfoResponse{
			Id:             subCategory.ID,
			Name:           subCategory.Name,
			Level:          subCategory.Level,
			IsTab:          subCategory.IsTab,
			ParentCategory: subCategory.ParentCategoryID,
		})
	}

	categoryListResponse.SubCategory = subCategoryResponse
	return &categoryListResponse, nil
}

//CreateCategory(context.Context, *CategoryInfoRequest) (*CategoryInfoResponse, error)
//DeleteCategory(context.Context, *DeleteCategoryRequest) (*emptypb.Empty, error)
//UpdateCategory(context.Context, *CategoryInfoRequest) (*emptypb.Empty, error)

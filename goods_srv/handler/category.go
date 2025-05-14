/* 商品分类 */

package handler

import (
	"context"
	"encoding/json"
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

//GetSubCategory(context.Context, *CategoryListRequest) (*SubCategoryListResponse, error)
//CreateCategory(context.Context, *CategoryInfoRequest) (*CategoryInfoResponse, error)
//DeleteCategory(context.Context, *DeleteCategoryRequest) (*emptypb.Empty, error)
//UpdateCategory(context.Context, *CategoryInfoRequest) (*emptypb.Empty, error)

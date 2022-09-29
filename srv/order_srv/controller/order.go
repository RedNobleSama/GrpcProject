/**
    @auther: oreki
    @date: 2022年09月28日 3:42 PM
    @note: 图灵老祖保佑,永无BUG
**/

package controller

import (
	"context"
	"google.golang.org/grpc/profiling/proto"
	"order_srv/db"
	in "order_srv/interface"
	"order_srv/model"
)

type OrderServer struct {
	proto.UnimplementedProfilingServer
}

func (o OrderServer) CartItemList(ctx context.Context, req *in.UserInfo) (*in.CartItemListResponse, error) {
	// 获取当前用户的购物车列表
	var shopCarts []model.ShoppingCart
	result := db.DB.Where(&model.ShoppingCart{User: req.Id}).Find(&shopCarts)
	if result.Error != nil {
		return nil, result.Error
	}

	rsp := in.CartItemListResponse{}
}

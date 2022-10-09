/**
    @auther: oreki
    @date: 2022年09月28日 3:42 PM
    @note: 图灵老祖保佑,永无BUG
**/

package controller

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/profiling/proto"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"math/rand"
	"order_srv/db"
	"order_srv/global"
	in "order_srv/interface"
	"order_srv/model"
	"time"
)

type OrderServer struct {
	proto.UnimplementedProfilingServer
}

func GenerateOrderSn(userId int32) string {
	// 订单号生成规则
	//年月日时分秒+用户ID+2位随机数
	now := time.Now()
	rand.Seed(time.Now().UnixNano())
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(),
		userId, rand.Intn(90)+10,
	)
	return orderSn
}

func (*OrderServer) SetInv(ctx context.Context, info *in.GoodsInvInfo) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (*OrderServer) InvDetail(ctx context.Context, info *in.GoodsInvInfo) (*in.GoodsInvInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (*OrderServer) Sell(ctx context.Context, info *in.SellInfo) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (*OrderServer) Reback(ctx context.Context, info *in.SellInfo) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (*OrderServer) CartItemList(ctx context.Context, req *in.UserInfo) (*in.CartItemListResponse, error) {
	// 获取当前用户的购物车列表
	var shopCarts []model.ShoppingCart
	result := db.DB.Where(&model.ShoppingCart{User: req.Id}).Find(&shopCarts)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := in.CartItemListResponse{
		Total: int32(result.RowsAffected),
	}

	for _, shopCart := range shopCarts {
		rsp.Data = append(rsp.Data, &in.ShopCartInfoResponse{
			Id:      shopCart.ID,
			UserId:  shopCart.User,
			GoodsId: shopCart.Goods,
			Nums:    shopCart.Nums,
			Checked: shopCart.Checked,
		})
	}
	return &rsp, nil
}

func (*OrderServer) CreateCartItem(ctx context.Context, req *in.CartItemRequest) (*in.ShopCartInfoResponse, error) {
	var shopCart model.ShoppingCart

	result := db.DB.Where(&model.ShoppingCart{
		Goods: req.GoodsId,
		User:  req.UserId,
	}).First(&shopCart)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 1 {
		// 如果记录已经存在，则合并记录
		shopCart.Nums += req.Nums
	} else {
		shopCart.User = req.UserId
		shopCart.Goods = req.GoodsId
		shopCart.Nums = req.Nums
		shopCart.Checked = false
	}
	db.DB.Save(&shopCart)
	return &in.ShopCartInfoResponse{Id: shopCart.ID}, nil
}

func (*OrderServer) UpdateCartItem(ctx context.Context, req *in.CartItemRequest) (*in.ShopCartInfoResponse, error) {
	var shopCart model.ShoppingCart

	result := db.DB.First(&shopCart, req.Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		// 如果记录不存在
		return nil, status.Error(codes.NotFound, "记录不存在")
	} else {
		shopCart.Nums = req.Nums
		shopCart.Checked = req.Checked
	}
	db.DB.Save(&shopCart)
	return &in.ShopCartInfoResponse{Id: shopCart.ID}, nil
}

func (*OrderServer) DeleteCartItem(ctx context.Context, req *in.CartItemRequest) (*emptypb.Empty, error) {
	result := db.DB.Delete(&model.ShoppingCart{}, req.Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "记录不存在")
	}
	return &emptypb.Empty{}, nil
}

func (*OrderServer) OrderList(ctx context.Context, req *in.OrderFilterRequest) (*in.OrderListResponse, error) {
	var orders []model.OrderInfo
	var rsp in.OrderListResponse
	var total int64
	db.DB.Where(&model.OrderInfo{User: req.UserId}).Count(&total)
	rsp.Total = int32(total)
	//分页
	db.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Where(&model.OrderInfo{User: req.UserId}).Find(&orders)
	for _, order := range orders {
		rsp.Data = append(rsp.Data, &in.OrderInfoResponse{
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
		})
	}
	return &rsp, nil
}

func (*OrderServer) OrderDetail(ctx context.Context, req *in.OrderRequest) (*in.OrderInfoDetailResponse, error) {
	var order model.OrderInfo
	var rsp in.OrderInfoDetailResponse

	//这个订单的id是否是当前用户的订单， 如果在web层用户传递过来一个id的订单， web层应该先查询一下订单id是否是当前用户的
	//在个人中心可以这样做，但是如果是后台管理系统，web层如果是后台管理系统 那么只传递order的id，如果是电商系统还需要一个用户的id
	if result := db.DB.Where(&model.OrderInfo{BaseModel: model.BaseModel{ID: req.Id}, User: req.UserId}).First(&order); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}

	orderInfo := in.OrderInfoResponse{}
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
	if result := db.DB.Where(&model.OrderGoods{Order: order.ID}).Find(&orderGoods); result.Error != nil {
		return nil, result.Error
	}

	for _, orderGood := range orderGoods {
		rsp.Goods = append(rsp.Goods, &in.OrderItemResponse{
			GoodsId:    orderGood.Goods,
			GoodsName:  orderGood.GoodsName,
			GoodsPrice: orderGood.GoodsPrice,
			GoodsImage: orderGood.GoodsImage,
			Nums:       orderGood.Nums,
		})
	}

	return &rsp, nil
}

func (*OrderServer) CreateOrder(ctx context.Context, req *in.OrderRequest) (*in.OrderInfoResponse, error) {
	// 从购物车中获取到选中的商品
	// 商品金额自己查询 - 访问商品服务 （跨微服务）
	// 商品的扣减 - 访问库存服务 （跨微服务）
	// 订单的基本信息表 - 订单的商品信息表
	// 从购物车中删除已购买的记录

	var goodsIds []int32
	var shopCarts []model.ShoppingCart
	goodsNumsMap := make(map[int32]int32)
	result := db.DB.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Find(&shopCarts)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "没有选中结算的商品 ")
	}
	for _, shopCart := range shopCarts {
		goodsIds = append(goodsIds, shopCart.Goods)
		goodsNumsMap[shopCart.Goods] = shopCart.Nums
	}

	// 跨服务调用商品微服务
	goods, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &in.BatchGoodsIdInfo{
		Id: goodsIds,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "批量查询商品信息失败")
	}
	var orderAmount float32
	var orderGoods []*model.OrderGoods
	var goodsInvInfo []*in.GoodsInvInfo
	for _, good := range goods.Data {
		orderAmount += good.ShopPrice * float32(goodsNumsMap[good.Id]) // 通过商品ID查询商品数量
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      good.Id,
			GoodsName:  good.Name,
			GoodsImage: good.GoodsFrontImage,
			GoodsPrice: good.ShopPrice,
			Nums:       goodsNumsMap[good.Id],
		})

		goodsInvInfo = append(goodsInvInfo, &in.GoodsInvInfo{
			GoodsId: good.Id,
			Num:     goodsNumsMap[good.Id],
		})
	}

	// 跨服务调用库存微服务进行库存扣减
	_, err = global.InventorySrvClient.Sell(context.Background(), &in.SellInfo{GoodsInfo: goodsInvInfo})
	if err != nil {
		return nil, status.Errorf(codes.ResourceExhausted, "扣减库存失败")
	}

	tx := db.DB.Begin() //开启事务
	//生成订单表
	// 20210308xxx
	order := model.OrderInfo{
		OrderSn:      GenerateOrderSn(req.UserId),
		OrderMount:   orderAmount,
		Address:      req.Address,
		SignerName:   req.Name,
		SingerMobile: req.Mobile,
		Post:         req.Post,
	}
	if result := tx.Save(&order); result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "创建订单失败")
	}

	for _, orderGood := range orderGoods {
		orderGood.Order = order.ID
	}

	// 批量插入orderGoods
	if result := tx.CreateInBatches(orderGoods, 100); result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "创建订单失败")
	}

	if result := tx.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Delete(model.ShoppingCart{}); result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "创建订单失败")
	}
	orderinfo := in.OrderInfoResponse{Id: order.ID, OrderSn: order.OrderSn, Total: order.OrderMount}
	return &orderinfo, nil
}

func (*OrderServer) UpdateOrderStatus(ctx context.Context, req *in.OrderStatus) (*emptypb.Empty, error) {
	//先查询，再更新 实际上有两条sql执行， select 和 update语句
	if result := db.DB.Model(&model.OrderInfo{}).Where("order_sn = ?", req.OrderSn).Update("status", req.Status); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}
	return &emptypb.Empty{}, nil
}

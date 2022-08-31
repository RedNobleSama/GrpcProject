/**
    @auther: oreki
    @date: 2022/5/23
    @note: 图灵老祖保佑,永无BUG
**/

package controller

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm/clause"
	"srv/inventory_srv/db"
	"srv/inventory_srv/model"
	in "srv/inventory_srv/rpc"
)

type InventoryServer struct {
}

func (i *InventoryServer) SetInv(ctx context.Context, req *in.GoodsInvInfo) (*emptypb.Empty, error) {
	//设置库存，更新库存
	var inv model.Inventory
	db.DB.First(&inv, req.GoodsId)
	inv.Goods = req.GoodsId
	inv.Stocks = req.Num
	db.DB.Save(&inv)
	return &emptypb.Empty{}, nil
}

func (i *InventoryServer) InvDetail(ctx context.Context, req *in.GoodsInvInfo) (*in.GoodsInvInfo, error) {
	//获取库存详情
	fmt.Println("InvDetail:", req.GoodsId)
	var inv model.Inventory
	if result := db.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv); result.RowsAffected == 0 {
		fmt.Println(result.Error)
		return nil, status.Errorf(codes.NotFound, "没有库存信息")
	}
	return &in.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}

func (i *InventoryServer) Sell(ctx context.Context, req *in.SellInfo) (*emptypb.Empty, error) {
	// 减少库存
	// 数据库基本的一个应用场景: 数据库事务
	// 并发情况之下，可能出现超卖的 情况, 会出现数据不一致的问题，需要分布式锁

	fmt.Println("Sell:", req.GoodsInfo)
	tx := db.DB.Begin() //开启事务
	for _, goodinfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Inventory{Goods: goodinfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Errorf(codes.NotFound, "没有库存信息")
		}
		if inv.Stocks < goodinfo.Num {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}

		//扣减
		inv.Stocks -= goodinfo.Num
		tx.Save(&inv)
	}
	tx.Commit() //手动提交事务操作
	return &emptypb.Empty{}, nil
}

func (i *InventoryServer) Reback(ctx context.Context, req *in.SellInfo) (*emptypb.Empty, error) {
	// 库存归还 1. 订单超时归还 2. 订单创建失败 3. 手动归还
	tx := db.DB.Begin()
	for _, goodinfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := db.DB.First(&inv, goodinfo.GoodsId); result.RowsAffected == 0 {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Errorf(codes.NotFound, "没有库存信息")
		}
		// 库存增加
		inv.Stocks += goodinfo.Num
		tx.Save(&inv)
	}
	tx.Commit() //手动提交事务操作
	return &emptypb.Empty{}, nil
}

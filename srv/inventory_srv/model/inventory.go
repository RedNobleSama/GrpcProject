/**
    @auther: oreki
    @date: 2022/5/15
    @note: 图灵老祖保佑,永无BUG
**/

package model

type Inventory struct {
	BaseModel
	Goods   int32 `gorm:"column:goods;type:int(11);not null;index" json:"goods"` //index 设置索引
	Stocks  int32 `gorm:"column:stocks;type:int(11);not null" json:"stocks"`     //库存
	Version int32 `gorm:"column:version;type:int(11);not null" json:"version"`   //分布式锁的乐观锁
}

type InventoryHistory struct {
	User    int32  `gorm:"column:user;type:int(11);not null" json:"user"`             //用户id
	Goods   int32  `gorm:"column:goods;type:int(11);not null" json:"goods"`           //商品id
	Nums    int32  `gorm:"column:nums;type:int(11);not null" json:"nums"`             //操作数量
	OrderSn string `gorm:"column:order_sn;type:varchar(64);not null" json:"order_sn"` //订单号
	Status  int32  //1. 表示库存是预扣减（幂等性） 2. 表示已经支付
}

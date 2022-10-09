/**
    @auther: oreki
    @date: 2022年10月09日 7:46 PM
    @note: 图灵老祖保佑,永无BUG
**/

package init

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	in "order_srv/interface"
)

var (
	GoodsSrvClient     in.GoodsClient
	InventorySrvClient in.InventoryClient
)

// GetFromConsul 从注册中心获取用户服务的信息
func GetFromConsul() (string, int, string, int) {
	GoodsSrvHost := ""
	GoodsSrvPort := 0
	InventorySrvHost := ""
	InventorySrvPort := 0
	cfg := api.DefaultConfig()
	cfg.Address = "127.0.0.1:8500"
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err.Error())
	}
	Goodsdata, err := client.Agent().ServicesWithFilter(`Service == "goods_srv"`)
	if err != nil {
		panic(err)
	}
	for _, value := range Goodsdata {
		GoodsSrvHost = value.Address
		GoodsSrvPort = value.Port
		break
	}

	Inventorydata, err := client.Agent().ServicesWithFilter(`Service == "inventory_srv"`)
	if err != nil {
		panic(err)
	}
	for _, value := range Inventorydata {
		InventorySrvHost = value.Address
		InventorySrvPort = value.Port
		break
	}

	return GoodsSrvHost, GoodsSrvPort, InventorySrvHost, InventorySrvPort
}

// 初始化第三方微服务的client
func InitSrvClient() {
	goodsSrvHost, goodsSrvPort, inventorySrvHost, inventorySrvPort := GetFromConsul()
	zap.S().Info("goos_srv地址", goodsSrvHost, goodsSrvPort)
	zap.S().Info("inventory_srv地址", inventorySrvHost, inventorySrvPort)
	//拨号连接用户grpc服务
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", goodsSrvHost, goodsSrvPort),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), //使用分布式调用方法为轮询调用
	)
	if err != nil {
		zap.S().Errorw("连接商品服务失败", "msg", err.Error())
		return
	}
	// 调用接口
	GoodsSrvClient = in.NewGoodsClient(goodsConn)
	// 将用户服务客户端对象赋值给全局变量

	inventoryConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", inventorySrvHost, inventorySrvPort),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), //使用分布式调用方法为轮询调用
	)
	if err != nil {
		zap.S().Errorw("连接库存服务失败", "msg", err.Error())
		return
	}
	// 调用接口
	InventorySrvClient = in.NewInventoryClient(inventoryConn)
	// 将用户服务客户端对象赋值给全局变量
}

/**
    @auther: oreki
    @date: 2022年08月31日 3:35 PM
    @note: 图灵老祖保佑,永无BUG
**/

package init

import (
	"fmt"
	"go.uber.org/zap"
	"goods_api/global"
	in "goods_api/interface"
	"google.golang.org/grpc"
)

// SrvInit 获取用户服务的grpc客户端
func SrvInit() {
	SrvHost, SrvPort := GetFromConsul()
	//UserSrvInfo := "user_srv"
	zap.S().Info("goods_srv地址", SrvHost, SrvPort)
	//拨号连接用户grpc服务
	userConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", SrvHost, SrvPort),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), //使用分布式调用方法为轮询调用
	)
	if err != nil {
		zap.S().Errorw("连接用户服务失败", "msg", err.Error())
		return
	}
	// 调用接口
	GoodsSrvClient := in.NewGoodsClient(userConn)
	global.GoodsSrvClient = GoodsSrvClient // 将用户服务客户端对象赋值给全局变量
}

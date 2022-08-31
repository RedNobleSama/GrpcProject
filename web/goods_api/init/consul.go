/**
    @auther: oreki
    @date: 2022年08月31日 10:31 AM
    @note: 图灵老祖保佑,永无BUG
**/

package init

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

// ConsulInit gin服务注册到consul
func ConsulInit() (*api.Client, string) {
	// gin服务注册到consul
	cfg := api.DefaultConfig()
	cfg.Address = "127.0.0.1:8500" // consul地址
	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().Errorw("consul服务初始化失败", err.Error())
	}

	// 生成consul注册对象
	reginstration := new(api.AgentServiceRegistration)
	reginstration.Name = "goods_api"
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	reginstration.ID = serviceID
	reginstration.Port = 50052
	reginstration.Tags = []string{"grpc", "api", "goods"}
	reginstration.Address = "127.0.0.1"

	// 生成对应的健康检查
	check := &api.AgentServiceCheck{
		HTTP:                           "http://127.0.0.1:50054/health",
		Timeout:                        "60s",
		Interval:                       "60s",
		DeregisterCriticalServiceAfter: "120s",
	}
	reginstration.Check = check

	err = client.Agent().ServiceRegister(reginstration)
	if err != nil {
		zap.S().Errorw("注册服务失败", err.Error())
	}

	return client, serviceID
}

// GetFromConsul 从注册中心获取用户服务的信息
func GetFromConsul() (string, int) {
	cfg := api.DefaultConfig()
	cfg.Address = "127.0.0.1:8500"
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err.Error())
	}
	data, err := client.Agent().ServicesWithFilter(`Service == "goods_srv"`)
	if err != nil {
		panic(err)
	}
	SrvHost := "" // 用户服务地址
	SrvPort := 0  // 用户服务端口
	for _, value := range data {
		SrvHost = value.Address
		SrvPort = value.Port
		break
	}

	return SrvHost, SrvPort
}

//// SrvInit 获取用户服务的grpc客户端
//func SrvInit() {
//	SrvHost, SrvPort := GetFromConsul()
//	//UserSrvInfo := "user_srv"
//	zap.S().Info("goods_srv地址", SrvHost, SrvPort)
//	//拨号连接用户grpc服务
//	userConn, err := grpc.Dial(
//		fmt.Sprintf("%s:%d", SrvHost, SrvPort),
//		grpc.WithInsecure(),
//		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), //使用分布式调用方法为轮询调用
//	)
//	if err != nil {
//		zap.S().Errorw("连接用户服务失败", "msg", err.Error())
//		return
//	}
//	// 调用接口
//	GoodsSrvClient := in.NewGoodsClient(userConn)
//	global.GoodsSrvClient = GoodsSrvClient // 将用户服务客户端对象赋值给全局变量
//}

/**
    @auther: oreki
    @date: 2022/4/26
    @note: 图灵老祖保佑,永无BUG
**/

package init

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"user_api/global"
	in "user_api/interface"
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
	reginstration.Name = "user_api"
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	reginstration.ID = serviceID
	reginstration.Port = 50052
	reginstration.Tags = []string{"grpc", "api", "user"}
	reginstration.Address = "127.0.0.1"

	// 生成对应的健康检查
	check := &api.AgentServiceCheck{
		HTTP:                           "http://127.0.0.1:50052/health",
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
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
	data, err := client.Agent().ServicesWithFilter(`Service == "user_srv"`)
	if err != nil {
		panic(err)
	}
	userSrvHost := "" // 用户服务地址
	userSrvPort := 0  // 用户服务端口
	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}

	return userSrvHost, userSrvPort
}

// SrvInit 获取用户服务的grpc客户端
func SrvInit() {
	userSrvHost, userSrvPort := GetFromConsul()
	//UserSrvInfo := "user_srv"
	zap.S().Info("user_srv地址", userSrvHost, userSrvPort)
	//拨号连接用户grpc服务
	userConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", userSrvHost, userSrvPort),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), //使用分布式调用方法为轮询调用
	)
	if err != nil {
		zap.S().Errorw("连接用户服务失败", "msg", err.Error())
		return
	}
	// 调用接口
	userSrvClient := in.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient // 将用户服务客户端对象赋值给全局变量
}

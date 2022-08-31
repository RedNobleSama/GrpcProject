/**
    @auther: oreki
    @date: 2022/4/27
    @note: 图灵老祖保佑,永无BUG
**/

package init

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// ConsulInit gin服务注册到consul
func ConsulInit(server *grpc.Server) (*api.Client, string) {
	zap.S().Info("consul init")
	grpc_health_v1.RegisterHealthServer(server, health.NewServer()) // 注册服务健康检查

	// grpc服务注册到consul
	cfg := api.DefaultConfig()
	cfg.Address = "127.0.0.1:8500"
	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().Errorw("consul init error", "err", err.Error())
		panic(err.Error())
	}
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("127.0.0.1:%d", 50051),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}
	// 生成consul注册对象
	reginstration := new(api.AgentServiceRegistration)
	reginstration.Name = "user_srv"
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	reginstration.ID = serviceID
	reginstration.Port = 50051
	reginstration.Tags = []string{"grpc", "srv", "user"}
	reginstration.Address = "127.0.0.1"
	reginstration.Check = check

	err = client.Agent().ServiceRegister(reginstration)
	if err != nil {
		zap.S().Errorw("consul注册失败，请检查consul是否启动", "err", err.Error())
		panic(err.Error)
	}

	return client, serviceID
}

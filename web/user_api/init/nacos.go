/**
    @auther: oreki
    @date: 2022/4/27
    @note: 图灵老祖保佑,永无BUG
**/

package init

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func NacosInit() {
	ns := constant.ServerConfig{
		IpAddr: "",
		Port:   80,
	}

	cc := constant.ClientConfig{
		NamespaceId:         "e525eafa-f7d7-4029-83d9-008937f9d468", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "./config/nacos/log",
		CacheDir:            "./config/nacos/cache",
		LogLevel:            "debug",
		Username:            "",
		Password:            "",
	}

	// 创建动态配置客户端
	configClient, _ := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": ns,
		"clientConfig":  cc,
	})

	// 获取配置：GetConfig
	_, _ = configClient.GetConfig(vo.ConfigParam{
		DataId: "dataId",
		Group:  "group"})

	// 监听配置变化：ListenConfig
	_ = configClient.ListenConfig(vo.ConfigParam{
		DataId: "dataId",
		Group:  "group",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
		},
	})
}

/**
    @auther: oreki
    @date: 2022/4/26
    @note: 图灵老祖保佑,永无BUG
**/

package main

import (
	"flag"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	initialize "user_api/init"
)

func main() {
	IP := flag.String("ip", "192.168.100.166", "ip地址")
	PORT := flag.String("port", "50052", "端口号")
	flag.Parse()

	initialize.LoggerInit()                            // 初始化logger
	Router := initialize.RouterInit()                  // 初始化路由
	initialize.SrvInit()                               // 连接rpc服务器
	ConsulClient, serviceID := initialize.ConsulInit() // 初始化consul注册
	go func() {
		zap.S().Info(fmt.Sprintf("启动服务 %s:%s", *IP, *PORT)) // 打印日志
		err := Router.Run(fmt.Sprintf("%s:%s", *IP, *PORT))
		if err != nil {
			zap.S().Panic("启动服务失败", zap.Error(err))
		} // 监听端口
	}()

	// 接收服务终止信号，并且在consul中注销服务
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	err := ConsulClient.Agent().ServiceDeregister(serviceID)
	if err != nil {
		zap.S().Panic("consul注销服务失败", zap.Error(err))
	}
	zap.S().Info("consul注销服务成功")
}

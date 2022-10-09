/**
    @auther: oreki
    @date: 2022/6/20
    @note: 图灵老祖保佑,永无BUG
**/

package main

import (
	"fmt"
	initialize "order_srv/init"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	initialize.LoggerInit() //初始化日志

	listen, server := initialize.ServerInit()                // 初始化rpc服务
	initialize.InitSrvClient()                               // 连接外部微服务
	ConsulClient, serviceID := initialize.ConsulInit(server) //初始化consul,注册服务

	go func() {
		err := server.Serve(listen)
		if err != nil {
			panic("failed to GRPC:" + err.Error())
		}
	}()

	// 接收服务终止信号，并且在consul中注销服务
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err := ConsulClient.Agent().ServiceDeregister(serviceID); err != nil {
		fmt.Println("注销失败")
	}
	fmt.Println("注销成功")
}

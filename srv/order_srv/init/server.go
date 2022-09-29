/**
    @auther: oreki
    @date: 2022年09月28日 3:39 PM
    @note: 图灵老祖保佑,永无BUG
**/

package init

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"order_srv/controller"
	in "order_srv/interface"
)

func ServerInit() (net.Listener, *grpc.Server) {
	zap.S().Info("init server")
	server := grpc.NewServer()
	//in.RegisterGoodsServer(server, &controller.GoodsServer{})
	in.RegisterInventoryServer(server, &controller.OrderServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", 50053))
	if err != nil {
		zap.S().Errorw("failed to listen:", err.Error())
		panic("failed to listen: " + err.Error())
	}
	return lis, server
}

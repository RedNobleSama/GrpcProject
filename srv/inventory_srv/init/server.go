/**
    @auther: oreki
    @date: 2022/5/14
    @note: 图灵老祖保佑,永无BUG
**/

package init

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"srv/inventory_srv/controller"
	in "srv/inventory_srv/rpc"
)

func ServerInit() (net.Listener, *grpc.Server) {
	zap.S().Info("init server")
	server := grpc.NewServer()
	//in.RegisterGoodsServer(server, &controller.GoodsServer{})
	in.RegisterInventoryServer(server, &controller.InventoryServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", 50053))
	if err != nil {
		zap.S().Errorw("failed to listen:", err.Error())
		panic("failed to listen: " + err.Error())
	}
	return lis, server
}

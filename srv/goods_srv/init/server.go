/**
    @auther: oreki
    @date: 2022/5/14
    @note: 图灵老祖保佑,永无BUG
**/

package init

import (
	"fmt"
	"go.uber.org/zap"
	"goods_srv/controller"
	in "goods_srv/interface"
	"google.golang.org/grpc"
	"net"
)

func ServerInit() (net.Listener, *grpc.Server) {
	zap.S().Info("init server")
	server := grpc.NewServer()
	in.RegisterGoodsServer(server, &controller.GoodsServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", 50053))
	if err != nil {
		zap.S().Errorw("failed to listen:", err.Error())
		panic("failed to listen: " + err.Error())
	}
	return lis, server
}

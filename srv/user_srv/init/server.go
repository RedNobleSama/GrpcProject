/**
    @auther: oreki
    @date: 2022/4/27
    @note: 图灵老祖保佑,永无BUG
**/

package init

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"user_srv/controller"
	in "user_srv/interface"
)

func ServerInit() (net.Listener, *grpc.Server) {
	zap.S().Info("init server")
	server := grpc.NewServer()
	in.RegisterUserServer(server, &controller.UserServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", 50051))
	if err != nil {
		zap.S().Errorw("failed to listen:", err.Error())
		panic("failed to listen: " + err.Error())
	}
	return lis, server
}

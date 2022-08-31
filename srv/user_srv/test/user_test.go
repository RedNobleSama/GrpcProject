/**
    @auther: oreki
    @date: 2022/4/25
    @note: 图灵老祖保佑,永无BUG
**/

package test

import (
	"context"
	"fmt"
	"github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc"
	"testing"
	in "user_srv/interface"
)

var UserSrv in.UserClient

// SrvInit 获取用户服务的grpc客户端
func init() {
	userSrvHost := "127.0.0.1"
	userSrvPort := 50051
	//拨号连接用户grpc服务
	userConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", userSrvHost, userSrvPort),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), //使用分布式调用方法为轮询调用
	)
	if err != nil {
		fmt.Println("连接用户服务失败", "msg", err.Error())
		return
	}
	// 调用接口
	userSrvClient := in.NewUserClient(userConn)
	UserSrv = userSrvClient // 将用户服务客户端对象赋值给全局变量
	fmt.Println("=== 连接用户服务成功")
}

func TestGetUserList(t *testing.T) {

}

func TestCreateUser(t *testing.T) {
	convey.Convey("创建用户测试用例", t, func() {
		convey.Convey("创建用户", func() {

			user := in.CreateUserInfo{
				NickName: "cyt",
				Password: "123456",
				Mobile:   "18970849456",
			}
			_, err := UserSrv.CreateUser(context.Background(), &user)
			convey.So(err, convey.ShouldBeNil)
			//convey.So(createUser, convey.ShouldNotBeNil, "创建用户失败")
		})
	})
}

func TestUpdateUser(t *testing.T) {

}

/**
    @auther: oreki
    @date: 2022/4/26
    @note: 图灵老祖保佑,永无BUG
**/

package router

import (
	"github.com/gin-gonic/gin"
	"user_api/api"
)

func InitUserRouter(router *gin.RouterGroup) {
	UserRouter := router.Group("user")
	{
		UserRouter.GET("/list", api.GetUserList)
	}
}

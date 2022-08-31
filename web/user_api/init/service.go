/**
    @auther: oreki
    @date: 2022/4/26
    @note: 图灵老祖保佑,永无BUG
**/

package init

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"user_api/router"
)

// RouterInit 初始化路由
func RouterInit() *gin.Engine {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "api health",
		})
	})
	ApiGroup := r.Group("/api")
	router.InitUserRouter(ApiGroup)
	return r
}

/**
    @auther: oreki
    @date: 2022年08月30日 11:46 PM
    @note: 图灵老祖保佑,永无BUG
**/

package init

import (
	"github.com/gin-gonic/gin"
	"goods_api/router"
	"net/http"
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
	router.InitGoodsRouter(ApiGroup)
	router.InitBannerRouter(ApiGroup)
	router.InitBrandRouter(ApiGroup)
	router.InitCategoryRouter(ApiGroup)
	return r
}

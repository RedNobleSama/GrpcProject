/**
    @auther: oreki
    @date: 2022年08月31日 11:38 PM
    @note: 图灵老祖保佑,永无BUG
**/

package router

import (
	"github.com/gin-gonic/gin"
	"goods_api/api/goods"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("goods")
	{
		GoodsRouter.GET("", goods.List)              //商品列表
		GoodsRouter.POST("", goods.New)              //改接口需要管理员权限
		GoodsRouter.GET("/:id", goods.Detail)        //获取商品的详情
		GoodsRouter.DELETE("/:id", goods.Delete)     //删除商品
		GoodsRouter.GET("/:id/stocks", goods.Stocks) //获取商品的库存

		GoodsRouter.PUT("/:id", goods.Update)
		GoodsRouter.PATCH("/:id", goods.UpdateStatus)
	}
}

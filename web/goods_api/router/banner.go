/**
    @auther: oreki
    @date: 2022年08月31日 11:37 PM
    @note: 图灵老祖保佑,永无BUG
**/

package router

import (
	"github.com/gin-gonic/gin"
	"goods_api/api/banner"
)

func InitBannerRouter(Router *gin.RouterGroup) {
	BannerRouter := Router.Group("banners")
	{
		BannerRouter.GET("", banner.List)          // 轮播图列表页
		BannerRouter.DELETE("/:id", banner.Delete) // 删除轮播图
		BannerRouter.POST("", banner.New)          //新建轮播图
		BannerRouter.PUT("/:id", banner.Update)    //修改轮播图信息
	}

	//BannerRouter := Router.Group("banners").Use(middlewares.Trace())
	//{
	//	BannerRouter.GET("", banners.List)                                                            // 轮播图列表页
	//	BannerRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), banners.Delete) // 删除轮播图
	//	BannerRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), banners.New)          //新建轮播图
	//	BannerRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), banners.Update)    //修改轮播图信息
	//}
}

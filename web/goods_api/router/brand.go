/**
    @auther: oreki
    @date: 2022年08月31日 11:37 PM
    @note: 图灵老祖保佑,永无BUG
**/

package router

import (
	"github.com/gin-gonic/gin"
	"goods_api/api/brands"
)

// 1. 商品的api接口开发完成
// 2. 图片的坑
func InitBrandRouter(Router *gin.RouterGroup) {
	CategoryBrandRouter := Router.Group("categorybrands")
	{
		CategoryBrandRouter.GET("", brands.CategoryBrandList)          // 类别品牌列表页
		CategoryBrandRouter.DELETE("/:id", brands.DeleteCategoryBrand) // 删除类别品牌
		CategoryBrandRouter.POST("", brands.NewCategoryBrand)          //新建类别品牌
		CategoryBrandRouter.PUT("/:id", brands.UpdateCategoryBrand)    //修改类别品牌
		CategoryBrandRouter.GET("/:id", brands.GetCategoryBrandList)   //获取分类的品牌
	}
}

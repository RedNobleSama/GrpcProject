/**
    @auther: oreki
    @date: 2022年08月30日 11:34 PM
    @note: 图灵老祖保佑,永无BUG
**/

package goods

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goods_api/api"
	"goods_api/global"
	in "goods_api/interface"
	"net/http"
	"strconv"
)

func List(c *gin.Context) {
	// 商品的列表
	request := &in.GoodsFilterRequest{}
	priceMin, _ := strconv.Atoi(c.DefaultQuery("pmin", "0"))
	priceMax, _ := strconv.Atoi(c.DefaultQuery("pmax", "0"))
	isHot := c.DefaultQuery("ishot", "0")
	if isHot == "1" {
		request.IsHot = true
	}
	isNew := c.DefaultQuery("isnew", "0")
	if isNew == "1" {
		request.IsNew = true
	}
	isTab := c.DefaultQuery("istab", "0")
	if isTab == "1" {
		request.IsTab = true
	}
	categoryId, _ := strconv.Atoi(c.DefaultQuery("cid", "0"))
	pages, _ := strconv.Atoi(c.DefaultQuery("p", "0"))
	pageNums, _ := strconv.Atoi(c.DefaultQuery("pn", "0"))
	keywords := c.DefaultQuery("key", "")
	brandId, _ := strconv.Atoi(c.DefaultQuery("bid", "0"))

	request.PriceMax = int32(priceMax)
	request.PriceMin = int32(priceMin)
	request.TopCategory = int32(categoryId)
	request.Pages = int32(pages)
	request.PagePerNums = int32(pageNums)
	request.Brand = int32(brandId)
	request.KeyWords = keywords

	rsp, err := global.GoodsSrvClient.GoodsList(context.Background(), request)
	if err != nil {
		zap.S().Errorw("[List] 查询 [商品列表] 失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	reMap := map[string]interface{}{
		"total": rsp.Total,
		"data":  rsp.Data,
	}
	c.JSON(http.StatusOK, reMap)
}

func New(c *gin.Context) {

}

func Detail(c *gin.Context) {

}

func Delete(c *gin.Context) {

}

func Stocks(c *gin.Context) {

}

func UpdateStatus(c *gin.Context) {

}

func Update(c *gin.Context) {

}

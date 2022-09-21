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
	"goods_api/struct/form"
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
	}

	goodsList := make([]interface{}, 0)
	for _, value := range rsp.Data {
		goodsList = append(goodsList, map[string]interface{}{
			"id":          value.Id,
			"name":        value.Name,
			"goods_brief": value.GoodsBrief,
			"desc":        value.GoodsDesc,
			"ship_free":   value.ShipFree,
			"images":      value.Images,
			"desc_images": value.DescImages,
			"front_image": value.GoodsFrontImage,
			"shop_price":  value.ShopPrice,
			"category": map[string]interface{}{
				"id":   value.Category.Id,
				"name": value.Category.Name,
			},
			"brand": map[string]interface{}{
				"id":   value.Brand.Id,
				"name": value.Brand.Name,
				"logo": value.Brand.Logo,
			},
			"is_hot":  value.IsHot,
			"is_new":  value.IsNew,
			"on_sale": value.OnSale,
		})
	}

	c.JSON(http.StatusOK, reMap)
}

func New(c *gin.Context) {
	// 新建商品
	goodsForm := form.GoodsForm{}
	if err := c.ShouldBindJSON(&goodsForm); err != nil {
		api.HandleValidatorError(c, err)
		return
	}
	goodsClient := global.GoodsSrvClient
	rsp, err := goodsClient.CreateGoods(context.Background(), &in.CreateGoodsInfo{
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	//如何设置库存
	//TODO 商品的库存 - 分布式事务
	c.JSON(http.StatusOK, rsp)

}

func Detail(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	r, err := global.GoodsSrvClient.GetGoodsDetail(context.WithValue(context.Background(), "ginContext", c), &in.GoodInfoRequest{
		Id: int32(i),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	rsp := map[string]interface{}{
		"id":          r.Id,
		"name":        r.Name,
		"goods_brief": r.GoodsBrief,
		"desc":        r.GoodsDesc,
		"ship_free":   r.ShipFree,
		"images":      r.Images,
		"desc_images": r.DescImages,
		"front_image": r.GoodsFrontImage,
		"shop_price":  r.ShopPrice,
		"ctegory": map[string]interface{}{
			"id":   r.Category.Id,
			"name": r.Category.Name,
		},
		"brand": map[string]interface{}{
			"id":   r.Brand.Id,
			"name": r.Brand.Name,
			"logo": r.Brand.Logo,
		},
		"is_hot":  r.IsHot,
		"is_new":  r.IsNew,
		"on_sale": r.OnSale,
	}
	c.JSON(http.StatusOK, rsp)
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsSrvClient.DeleteGoods(context.Background(), &in.DeleteGoodsInfo{Id: int32(i)})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.Status(http.StatusOK)
	return
}

func Stocks(c *gin.Context) {
	id := c.Param("id")
	_, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	//TODO 商品的库存
	return
}

func UpdateStatus(c *gin.Context) {
	goodsStatusForm := form.GoodsStatusForm{}
	if err := c.ShouldBindJSON(&goodsStatusForm); err != nil {
		api.HandleValidatorError(c, err)
		return
	}

	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if _, err = global.GoodsSrvClient.UpdateGoods(context.Background(), &in.CreateGoodsInfo{
		Id:     int32(i),
		IsHot:  *goodsStatusForm.IsHot,
		IsNew:  *goodsStatusForm.IsNew,
		OnSale: *goodsStatusForm.OnSale,
	}); err != nil {
		api.HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "修改成功",
	})
}

func Update(c *gin.Context) {
	goodsForm := form.GoodsForm{}
	if err := c.ShouldBindJSON(&goodsForm); err != nil {
		api.HandleValidatorError(c, err)
		return
	}

	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if _, err = global.GoodsSrvClient.UpdateGoods(context.Background(), &in.CreateGoodsInfo{
		Id:              int32(i),
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	}); err != nil {
		api.HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}

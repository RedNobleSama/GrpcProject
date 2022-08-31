/**
    @auther: oreki
    @date: 2022年08月31日 4:14 PM
    @note: 图灵老祖保佑,永无BUG
**/

package form

type BannerForm struct {
	Image string `form:"image" json:"image" binding:"url"`
	Index int    `form:"index" json:"index" binding:"required"`
	Url   string `form:"url" json:"url" binding:"url"`
}

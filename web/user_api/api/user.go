/**
    @auther: oreki
    @date: 2022/4/26
    @note: 图灵老祖保佑,永无BUG
**/

package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
	"user_api/global"
	"user_api/global/HandleGrpcError"
	"user_api/global/response"
	in "user_api/interface"
)

func GetUserList(c *gin.Context) {
	zap.S().Info("获取用户列表页")

	pNum, _ := strconv.Atoi(c.DefaultQuery("pNum", "0"))
	pSize, _ := strconv.Atoi(c.DefaultQuery("pSize", "0"))

	userlist, err := global.UserSrvClient.GetUserList(context.Background(), &in.PageInfo{
		PNum:  int64(pNum),
		PSize: int64(pSize),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询用户列表失败", "msg:", err.Error())
		HandleGrpcError.HandleGrpcErrorToHttp(err, c)
		panic(err)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range userlist.Data {
		userRes := response.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			Mobile:   value.Mobile,
			BirthDay: response.JsonTime(time.Unix(value.Birthday, 0)),
			Gender:   value.Gender,
			Role:     value.Role,
		}

		result = append(result, userRes)
	}
	c.JSON(http.StatusOK, result)
}

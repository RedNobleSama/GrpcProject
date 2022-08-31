/**
    @auther: oreki
    @date: 2022/4/26
    @note: 图灵老祖保佑,永无BUG
**/

package response

import (
	"fmt"
	"time"
)

type JsonTime time.Time

func (J JsonTime) MarshalJSON() ([]byte, error) {
	var stmp = fmt.Sprintf("\"%s\"", time.Time(J).Format("2006-01-02"))
	return []byte(stmp), nil
}

type UserResponse struct {
	Id       int64    `json:"id"`
	Password string   `json:"password"`
	Mobile   string   `json:"mobile"`
	NickName string   `json:"nickName"`
	BirthDay JsonTime `json:"birthDay"`
	Gender   string   `json:"gender"`
	Role     int64    `json:"role"`
}

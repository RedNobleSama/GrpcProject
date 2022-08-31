/**
    @auther: oreki
    @date: 2022/5/14
    @note: 图灵老祖保佑,永无BUG
**/

package utils

import "time"

func GetNowTime() string {
	currentTime := time.Now().Format("2006-01-02")
	return currentTime
}

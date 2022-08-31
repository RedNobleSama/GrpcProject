/**
    @auther: oreki
    @date: 2022年08月31日 10:29 AM
    @note: 图灵老祖保佑,永无BUG
**/

package utils

import "time"

func GetNowTime() string {
	currentTime := time.Now().Format("2006-01-02")
	return currentTime
}

/**
    @auther: oreki
    @date: 2022年09月28日 3:41 PM
    @note: 图灵老祖保佑,永无BUG
**/

package utils

import "time"

func GetNowTime() string {
	currentTime := time.Now().Format("2006-01-02")
	return currentTime
}

/**
    @auther: oreki
    @date: 2022/5/12
    @note: 图灵老祖保佑,永无BUG
**/

package controller

import "gorm.io/gorm"

// Paginate gorm 通用分页
func Paginate(pageSize, pageNumber int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageNumber == 0 {
			pageNumber = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (pageNumber - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

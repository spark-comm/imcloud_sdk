package db

import "gorm.io/gorm"

func SqlDataLimit(pageSize, pageNum int) func(db *gorm.DB) *gorm.DB {
	if pageSize == 0 {
		pageSize = 10
	}
	if pageNum == 0 {
		pageNum = 1
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(pageSize).Offset(pageNum - 1*pageSize)
	}
}

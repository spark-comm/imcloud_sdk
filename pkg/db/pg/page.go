package pg

import "gorm.io/gorm"

// Page 分页参数
type Page struct {
	Size int64 `json:"size"` // 页码大小，最大100
	NO   int64 `json:"no"`   // 页码，从1开始
}

func (r *Page) Fix() {
	if r.NO <= 0 {
		r.NO = 1
	}

	if r.Size <= 0 {
		r.Size = 10
	} else if r.Size > 100 {
		r.Size = 100
	}
}
func Operation(r *Page) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		r.Fix()
		offset := int((r.NO - 1) * r.Size)
		return db.Offset(offset).Limit(int(r.Size))
	}
}

package scope

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

func ScopeCategoryByName(name string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(name)+"%")
	}
}

func ScopeCategoryCreatedAfter(t time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("created_at >= ?", t)
	}
}

func ScopeCategoryCreatedBefore(t time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("created_at <= ?", t)
	}
}

func ScopeCategoryWithProducts() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload("Products")
	}
}

func ScopeCategoryByMinProducts(min int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Joins("LEFT JOIN products p ON p.category_id = categories.id").
			Group("categories.id").
			Having("COUNT(p.id) >= ?", min)
	}
}

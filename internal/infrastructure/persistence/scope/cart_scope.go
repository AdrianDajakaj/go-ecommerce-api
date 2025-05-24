package scope

import (
	"time"

	"gorm.io/gorm"
)

func ScopeCartByUser(userID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID)
	}
}

func ScopeCartByTotalRange(min, max float64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("total BETWEEN ? AND ?", min, max)
	}
}

func ScopeCartWithItems() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload("Items.Product")
	}
}

func ScopeCartCreatedAfter(t time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("created_at >= ?", t)
	}
}

func ScopeCartCreatedBefore(t time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("created_at <= ?", t)
	}
}

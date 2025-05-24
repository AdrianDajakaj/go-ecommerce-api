package scope

import (
	"time"

	"gorm.io/gorm"
)

func ScopeByUser(userID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID)
	}
}

func ScopeByStatus(status string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", status)
	}
}

func ScopeByTotalRange(min, max float64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("total BETWEEN ? AND ?", min, max)
	}
}

func ScopeCreatedAfter(t time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("created_at >= ?", t)
	}
}

func ScopeCreatedBefore(t time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("created_at <= ?", t)
	}
}

func ScopeWithAssociations() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Preload("User").
			Preload("ShippingAddress").
			Preload("Items")
	}
}

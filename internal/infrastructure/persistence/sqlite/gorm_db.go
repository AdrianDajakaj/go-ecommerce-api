package sqlite

import (
	"go-ecommerce-api/internal/domain/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewGormDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	models := []interface{}{
		&model.User{},
		&model.Address{},
		&model.Category{},
		&model.Product{},
		&model.ProductImage{},
		&model.Cart{},
		&model.CartItem{},
		&model.Order{},
		&model.OrderItem{},
	}

	if err := db.AutoMigrate(models...); err != nil {
		return nil, err
	}

	return db, nil
}

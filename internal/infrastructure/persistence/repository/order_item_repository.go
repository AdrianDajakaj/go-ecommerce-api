package repository

import (
	"go-ecommerce-api/internal/domain/model"
	"go-ecommerce-api/internal/domain/repository"

	"gorm.io/gorm"
)

type orderItemRepository struct {
	db *gorm.DB
}

func NewOrderItemRepository(db *gorm.DB) repository.OrderItemRepository {
	return &orderItemRepository{db: db}
}

func (r *orderItemRepository) FindByOrderID(orderID uint) ([]model.OrderItem, error) {
	var items []model.OrderItem
	if err := r.db.Where("order_id = ?", orderID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

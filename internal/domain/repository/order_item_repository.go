package repository

import "go-ecommerce-api/internal/domain/model"

type OrderItemRepository interface {
	FindByOrderID(orderID uint) ([]model.OrderItem, error)
}

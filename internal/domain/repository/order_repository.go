package repository

import "go-ecommerce-api/internal/domain/model"

type OrderRepository interface {
	FindByID(id uint) (*model.Order, error)
	FindByUserID(userID uint) ([]model.Order, error)
	FindAll() ([]model.Order, error)
	FindWithFilters(filters map[string]string) ([]model.Order, error)
	Create(order *model.Order) error
	Update(order *model.Order) error
}

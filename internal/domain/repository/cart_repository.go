package repository

import "go-ecommerce-api/internal/domain/model"

type CartRepository interface {
	FindByUserID(userID uint) (*model.Cart, error)
	FindWithFilters(filters map[string]string) ([]model.Cart, error)
	Create(cart *model.Cart) error
	Update(cart *model.Cart) error
	Delete(cartID uint) error
}

package repository

import "go-ecommerce-api/internal/domain/model"

type CartItemRepository interface {
	FindByID(id uint) (*model.CartItem, error)
	FindByCartID(cartID uint) ([]model.CartItem, error)
	AddItem(item *model.CartItem) error
	UpdateItem(item *model.CartItem) error
	DeleteItem(id uint) error
	ClearCart(cartID uint) error
}

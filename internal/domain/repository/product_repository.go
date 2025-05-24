package repository

import "go-ecommerce-api/internal/domain/model"

type ProductRepository interface {
	FindByID(id uint) (*model.Product, error)
	FindAll() ([]model.Product, error)
	FindWithFilters(filters map[string]string) ([]model.Product, error)
	Create(product *model.Product) error
	Update(product *model.Product) error
	Delete(id uint) error
}

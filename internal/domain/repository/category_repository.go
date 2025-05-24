package repository

import "go-ecommerce-api/internal/domain/model"

type CategoryRepository interface {
	FindByID(id uint) (*model.Category, error)
	FindAll() ([]model.Category, error)
	FindWithFilters(filters map[string]string) ([]model.Category, error)
	Create(category *model.Category) error
	Update(category *model.Category) error
	Delete(id uint) error
}

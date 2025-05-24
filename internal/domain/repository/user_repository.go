package repository

import "go-ecommerce-api/internal/domain/model"

type UserRepository interface {
	FindByID(id uint) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindAll() ([]model.User, error)
	FindWithFilters(filters map[string]string) ([]model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id uint) error
}

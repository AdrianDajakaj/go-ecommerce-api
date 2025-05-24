package repository

import "go-ecommerce-api/internal/domain/model"

type AddressRepository interface {
	FindByID(id uint) (*model.Address, error)
	Create(address *model.Address) error
	Update(address *model.Address) error
	Delete(id uint) error
}

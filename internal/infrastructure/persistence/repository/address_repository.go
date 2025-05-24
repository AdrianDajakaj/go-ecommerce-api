package repository

import (
	"errors"
	"go-ecommerce-api/internal/domain/model"
	"go-ecommerce-api/internal/domain/repository"

	"gorm.io/gorm"
)

type addressRepo struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) repository.AddressRepository {
	return &addressRepo{db: db}
}

func (r *addressRepo) FindByID(id uint) (*model.Address, error) {
	var addr model.Address
	if err := r.db.First(&addr, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &addr, nil
}

func (r *addressRepo) Create(address *model.Address) error {
	return r.db.Create(address).Error
}

func (r *addressRepo) Update(address *model.Address) error {
	return r.db.Save(address).Error
}

func (r *addressRepo) Delete(id uint) error {
	return r.db.Delete(&model.Address{}, id).Error
}

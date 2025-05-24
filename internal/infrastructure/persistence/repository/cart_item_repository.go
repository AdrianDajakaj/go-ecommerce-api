package repository

import (
	"errors"
	"go-ecommerce-api/internal/domain/model"
	"go-ecommerce-api/internal/domain/repository"

	"gorm.io/gorm"
)

type cartItemRepository struct {
	db *gorm.DB
}

func NewCartItemRepository(db *gorm.DB) repository.CartItemRepository {
	return &cartItemRepository{db}
}

func (r *cartItemRepository) FindByID(id uint) (*model.CartItem, error) {
	var item model.CartItem
	if err := r.db.Preload("Product").First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *cartItemRepository) FindByCartID(cartID uint) ([]model.CartItem, error) {
	var items []model.CartItem
	if err := r.db.Preload("Product").
		Where("cart_id = ?", cartID).
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *cartItemRepository) AddItem(item *model.CartItem) error {
	return r.db.Create(item).Error
}

func (r *cartItemRepository) UpdateItem(item *model.CartItem) error {
	result := r.db.Save(item)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *cartItemRepository) DeleteItem(id uint) error {
	result := r.db.Delete(&model.CartItem{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *cartItemRepository) ClearCart(cartID uint) error {
	return r.db.Where("cart_id = ?", cartID).Delete(&model.CartItem{}).Error
}

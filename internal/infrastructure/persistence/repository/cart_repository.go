package repository

import (
	"errors"
	"go-ecommerce-api/internal/domain/model"
	"go-ecommerce-api/internal/domain/repository"
	"go-ecommerce-api/internal/infrastructure/persistence/scope"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) repository.CartRepository {
	return &cartRepository{db}
}

func (r *cartRepository) FindByUserID(userID uint) (*model.Cart, error) {
	var cart model.Cart
	if err := r.db.Preload("Items.Product").
		Where("user_id = ?", userID).
		First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) FindByCartID(cartID uint) (*model.Cart, error) {
	var cart model.Cart
	if err := r.db.Preload("Items.Product").
		Where("id = ?", cartID).
		First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) FindWithFilters(filters map[string]string) ([]model.Cart, error) {
	db := r.db.Model(&model.Cart{})
	db = db.Scopes(scope.ScopeCartWithItems())

	db = r.applyUserFilter(db, filters)
	db = r.applyTotalRangeFilter(db, filters)
	db = r.applyCreatedAfterFilter(db, filters)
	db = r.applyCreatedBeforeFilter(db, filters)

	var carts []model.Cart
	if err := db.Find(&carts).Error; err != nil {
		return nil, err
	}

	return carts, nil
}

func (r *cartRepository) applyUserFilter(db *gorm.DB, filters map[string]string) *gorm.DB {
	if v, ok := filters["user_id"]; ok {
		if id, err := strconv.ParseUint(v, 10, 64); err == nil {
			db = db.Scopes(scope.ScopeCartByUser(uint(id)))
		}
	}
	return db
}

func (r *cartRepository) applyTotalRangeFilter(db *gorm.DB, filters map[string]string) *gorm.DB {
	vMin, okMin := filters["total_min"]
	vMax, okMax := filters["total_max"]

	if okMin && okMax {
		if min, err1 := strconv.ParseFloat(vMin, 64); err1 == nil {
			if max, err2 := strconv.ParseFloat(vMax, 64); err2 == nil {
				db = db.Scopes(scope.ScopeCartByTotalRange(min, max))
			}
		}
	}
	return db
}

func (r *cartRepository) applyCreatedAfterFilter(db *gorm.DB, filters map[string]string) *gorm.DB {
	if v, ok := filters["created_after"]; ok {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			db = db.Scopes(scope.ScopeCartCreatedAfter(t))
		}
	}
	return db
}

func (r *cartRepository) applyCreatedBeforeFilter(db *gorm.DB, filters map[string]string) *gorm.DB {
	if v, ok := filters["created_before"]; ok {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			db = db.Scopes(scope.ScopeCartCreatedBefore(t))
		}
	}
	return db
}

func (r *cartRepository) Create(cart *model.Cart) error {
	return r.db.Create(cart).Error
}

func (r *cartRepository) Update(cart *model.Cart) error {
	result := r.db.Save(cart)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *cartRepository) Delete(cartID uint) error {
	result := r.db.Delete(&model.Cart{}, cartID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

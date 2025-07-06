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

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) repository.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) FindByID(id uint) (*model.Order, error) {
	var order model.Order
	if err := r.db.Preload("User").
		Preload("ShippingAddress").
		Preload("Items").
		First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindByUserID(userID uint) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Preload("ShippingAddress").
		Preload("Items").
		Where("user_id = ?", userID).
		Find(&orders).Error
	return orders, err
}

func (r *orderRepository) FindAll() ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Preload("User").
		Preload("ShippingAddress").
		Preload("Items").
		Find(&orders).Error
	return orders, err
}

func (r *orderRepository) FindWithFilters(filters map[string]string) ([]model.Order, error) {
	db := r.db.Model(&model.Order{}).
		Scopes(scope.ScopeWithAssociations())

	r.applyUserFilter(db, filters)
	r.applyStatusFilter(db, filters)
	r.applyTotalRangeFilter(db, filters)
	r.applyTimeFilters(db, filters)

	var orders []model.Order
	if err := db.Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) applyUserFilter(db *gorm.DB, filters map[string]string) {
	if v, ok := filters["user_id"]; ok {
		if id, err := strconv.ParseUint(v, 10, 64); err == nil {
			db.Scopes(scope.ScopeByUser(uint(id)))
		}
	}
}

func (r *orderRepository) applyStatusFilter(db *gorm.DB, filters map[string]string) {
	if v, ok := filters["status"]; ok {
		db.Scopes(scope.ScopeByStatus(v))
	}
}

func (r *orderRepository) applyTotalRangeFilter(db *gorm.DB, filters map[string]string) {
	vMin, okMin := filters["total_min"]
	vMax, okMax := filters["total_max"]

	if okMin && okMax {
		min, err1 := strconv.ParseFloat(vMin, 64)
		max, err2 := strconv.ParseFloat(vMax, 64)

		if err1 == nil && err2 == nil {
			db.Scopes(scope.ScopeByTotalRange(min, max))
		}
	}
}

func (r *orderRepository) applyTimeFilters(db *gorm.DB, filters map[string]string) {
	if v, ok := filters["created_after"]; ok {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			db.Scopes(scope.ScopeCreatedAfter(t))
		}
	}
	if v, ok := filters["created_before"]; ok {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			db.Scopes(scope.ScopeCreatedBefore(t))
		}
	}
}

func (r *orderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) Update(order *model.Order) error {
	result := r.db.Save(order)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

package repository

import (
	"errors"
	"go-ecommerce-api/internal/domain/model"
	"go-ecommerce-api/internal/domain/repository"
	"go-ecommerce-api/internal/infrastructure/persistence/scope"
	"strconv"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) repository.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindByID(id uint) (*model.Product, error) {
	var prod model.Product
	if err := r.db.
		Preload("Category").
		Preload("Images").
		First(&prod, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &prod, nil
}

func (r *productRepository) FindAll() ([]model.Product, error) {
	var prods []model.Product
	err := r.db.
		Preload("Category").
		Preload("Images").
		Find(&prods).Error
	return prods, err
}

func (r *productRepository) FindWithFilters(filters map[string]string) ([]model.Product, error) {
	db := r.db.Model(&model.Product{}).
		Preload("Category").
		Preload("Images")

	if v, ok := filters["category_id"]; ok {
		if id, err := strconv.Atoi(v); err == nil {
			db = db.Scopes(scope.ScopeProductByCategory(uint(id)))
		}
	}
	if v, ok := filters["name"]; ok {
		db = db.Scopes(scope.ScopeProductByName(v))
	}
	if v, ok := filters["is_active"]; ok && (v == "true" || v == "false") {
		db = db.Scopes(scope.ScopeProductByIsActive(v == "true"))
	}
	if min, ok1 := filters["price_min"]; ok1 {
		if max, ok2 := filters["price_max"]; ok2 {
			if fmin, err1 := strconv.ParseFloat(min, 64); err1 == nil {
				if fmax, err2 := strconv.ParseFloat(max, 64); err2 == nil {
					db = db.Scopes(scope.ScopeProductByPriceRange(fmin, fmax))
				}
			}
		}
	}

	var products []model.Product
	if err := db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) Update(product *model.Product) error {
	result := r.db.Save(product)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *productRepository) Delete(id uint) error {
	result := r.db.Delete(&model.Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

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

	r.applyCategoryFilter(db, filters)
	r.applyNameFilter(db, filters)
	r.applyActiveFilter(db, filters)
	r.applyPriceRangeFilter(db, filters)

	var products []model.Product
	if err := db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) applyCategoryFilter(db *gorm.DB, filters map[string]string) {
	if v, ok := filters["category_id"]; ok {
		if id, err := strconv.Atoi(v); err == nil {
			db.Scopes(scope.ScopeProductByCategory(uint(id)))
		}
	}
}

func (r *productRepository) applyNameFilter(db *gorm.DB, filters map[string]string) {
	if v, ok := filters["name"]; ok {
		db.Scopes(scope.ScopeProductByName(v))
	}
}

func (r *productRepository) applyActiveFilter(db *gorm.DB, filters map[string]string) {
	if v, ok := filters["is_active"]; ok && (v == "true" || v == "false") {
		db.Scopes(scope.ScopeProductByIsActive(v == "true"))
	}
}

func (r *productRepository) applyPriceRangeFilter(db *gorm.DB, filters map[string]string) {
	min, okMin := filters["price_min"]
	max, okMax := filters["price_max"]

	if okMin && okMax {
		fmin, err1 := strconv.ParseFloat(min, 64)
		fmax, err2 := strconv.ParseFloat(max, 64)

		if err1 == nil && err2 == nil {
			db.Scopes(scope.ScopeProductByPriceRange(fmin, fmax))
		}
	}
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

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

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) repository.CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) FindByID(id uint) (*model.Category, error) {
	var category model.Category
	if err := r.db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) FindAll() ([]model.Category, error) {
	var categories []model.Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepository) FindWithFilters(filters map[string]string) ([]model.Category, error) {
	db := r.db.Model(&model.Category{})

	if _, ok := filters["with_products"]; ok {
		db = db.Scopes(scope.ScopeCategoryWithProducts())
	}

	if v, ok := filters["name"]; ok {
		db = db.Scopes(scope.ScopeCategoryByName(v))
	}
	if v, ok := filters["created_after"]; ok {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			db = db.Scopes(scope.ScopeCategoryCreatedAfter(t))
		}
	}
	if v, ok := filters["created_before"]; ok {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			db = db.Scopes(scope.ScopeCategoryCreatedBefore(t))
		}
	}
	if v, ok := filters["min_products"]; ok {
		if n, err := strconv.Atoi(v); err == nil {
			db = db.Scopes(scope.ScopeCategoryByMinProducts(n))
		}
	}

	var cats []model.Category
	if err := db.Find(&cats).Error; err != nil {
		return nil, err
	}
	return cats, nil
}

func (r *categoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) Update(category *model.Category) error {
	result := r.db.Save(category)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *categoryRepository) Delete(id uint) error {
	result := r.db.Delete(&model.Category{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

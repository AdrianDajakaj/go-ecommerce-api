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
	if err := r.db.Preload("Subcategories").
		Preload("Products").
		Preload("ParentCategory").
		First(&category, id).Error; err != nil {

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

	r.applyBooleanFilters(db, filters)
	r.applyStringFilters(db, filters)
	r.applyTimeFilters(db, filters)
	r.applyNumericFilters(db, filters)

	var cats []model.Category
	if err := db.Preload("Subcategories").Find(&cats).Error; err != nil {
		return nil, err
	}
	return cats, nil
}

func (r *categoryRepository) applyBooleanFilters(db *gorm.DB, filters map[string]string) {
	if _, ok := filters["with_products"]; ok {
		db.Scopes(scope.ScopeCategoryWithProducts())
	}
	if _, ok := filters["with_subcategories"]; ok {
		db.Scopes(scope.ScopeCategoryWithSubcategories())
	}
}

func (r *categoryRepository) applyStringFilters(db *gorm.DB, filters map[string]string) {
	if v, ok := filters["name"]; ok {
		db.Scopes(scope.ScopeCategoryByName(v))
	}
}

func (r *categoryRepository) applyTimeFilters(db *gorm.DB, filters map[string]string) {
	if v, ok := filters["created_after"]; ok {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			db.Scopes(scope.ScopeCategoryCreatedAfter(t))
		}
	}
	if v, ok := filters["created_before"]; ok {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			db.Scopes(scope.ScopeCategoryCreatedBefore(t))
		}
	}
}

func (r *categoryRepository) applyNumericFilters(db *gorm.DB, filters map[string]string) {
	if v, ok := filters["min_products"]; ok {
		if n, err := strconv.Atoi(v); err == nil {
			db.Scopes(scope.ScopeCategoryByMinProducts(n))
		}
	}
	if v, ok := filters["parent_id"]; ok {
		if id, err := strconv.Atoi(v); err == nil {
			db.Scopes(scope.ScopeCategoryByParentID(uint(id)))
		}
	}
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

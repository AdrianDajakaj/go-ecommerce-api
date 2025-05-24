package usecase

import (
	"errors"

	"go-ecommerce-api/internal/domain/model"
	"go-ecommerce-api/internal/domain/repository"

	"gorm.io/gorm"
)

type CategoryUsecase interface {
	GetByID(id uint) (*model.Category, error)
	GetAll() ([]model.Category, error)
	GetWithFilters(filters map[string]string) ([]model.Category, error)
	Create(category *model.Category) (*model.Category, error)
	Update(category *model.Category) (*model.Category, error)
	Delete(id uint) error
}

type categoryUsecase struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryUsecase(categoryRepo repository.CategoryRepository) CategoryUsecase {
	return &categoryUsecase{
		categoryRepo: categoryRepo,
	}
}

func (u *categoryUsecase) GetByID(id uint) (*model.Category, error) {
	category, err := u.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return category, nil
}

func (u *categoryUsecase) GetAll() ([]model.Category, error) {
	return u.categoryRepo.FindAll()
}

func (u *categoryUsecase) GetWithFilters(filters map[string]string) ([]model.Category, error) {
	return u.categoryRepo.FindWithFilters(filters)
}

func (u *categoryUsecase) Create(category *model.Category) (*model.Category, error) {
	if category == nil || category.Name == "" {
		return nil, errors.New("invalid category data")
	}
	if err := u.categoryRepo.Create(category); err != nil {
		return nil, err
	}
	return u.categoryRepo.FindByID(category.ID)
}

func (u *categoryUsecase) Update(category *model.Category) (*model.Category, error) {
	if category == nil || category.ID == 0 {
		return nil, errors.New("invalid category")
	}
	if err := u.categoryRepo.Update(category); err != nil {
		return nil, err
	}
	return u.categoryRepo.FindByID(category.ID)
}

func (u *categoryUsecase) Delete(id uint) error {
	category, err := u.categoryRepo.FindByID(id)
	if err != nil {
		return err
	}
	if category == nil {
		return gorm.ErrRecordNotFound
	}
	return u.categoryRepo.Delete(id)
}

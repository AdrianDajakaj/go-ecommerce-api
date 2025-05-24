package usecase

import (
	"errors"

	"go-ecommerce-api/internal/domain/model"
	"go-ecommerce-api/internal/domain/repository"

	"gorm.io/gorm"
)

type ProductUsecase interface {
	GetByID(id uint) (*model.Product, error)
	GetAll() ([]model.Product, error)
	GetWithFilters(filters map[string]string) ([]model.Product, error)
	Create(product *model.Product) (*model.Product, error)
	Update(product *model.Product) (*model.Product, error)
	Delete(id uint) error
}

type productUsecase struct {
	productRepo repository.ProductRepository
}

func NewProductUsecase(productRepo repository.ProductRepository) ProductUsecase {
	return &productUsecase{productRepo: productRepo}
}

func (u *productUsecase) GetByID(id uint) (*model.Product, error) {
	prod, err := u.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if prod == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return prod, nil
}

func (u *productUsecase) GetAll() ([]model.Product, error) {
	return u.productRepo.FindAll()
}

func (u *productUsecase) GetWithFilters(filters map[string]string) ([]model.Product, error) {
	return u.productRepo.FindWithFilters(filters)
}

func (u *productUsecase) Create(product *model.Product) (*model.Product, error) {
	if product == nil || product.Name == "" {
		return nil, errors.New("invalid product data")
	}
	if err := u.productRepo.Create(product); err != nil {
		return nil, err
	}
	return u.productRepo.FindByID(product.ID)
}

func (u *productUsecase) Update(product *model.Product) (*model.Product, error) {
	if product == nil || product.ID == 0 {
		return nil, errors.New("invalid product")
	}
	if err := u.productRepo.Update(product); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return u.productRepo.FindByID(product.ID)
}

func (u *productUsecase) Delete(id uint) error {
	prod, err := u.productRepo.FindByID(id)
	if err != nil {
		return err
	}
	if prod == nil {
		return gorm.ErrRecordNotFound
	}
	return u.productRepo.Delete(id)
}

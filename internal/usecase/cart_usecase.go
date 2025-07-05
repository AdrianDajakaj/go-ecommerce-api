package usecase

import (
	"errors"
	"fmt"
	"go-ecommerce-api/internal/domain/model"
	"go-ecommerce-api/internal/domain/repository"

	"gorm.io/gorm"
)

type CartUsecase interface {
	GetByUserID(userID uint) (*model.Cart, error)
	GetWithFilters(filters map[string]string) ([]model.Cart, error)
	AddProduct(userID, productID uint, quantity int) (*model.Cart, error)
	UpdateItem(itemID uint, quantity int) (*model.Cart, error)
	RemoveItem(itemID uint) (*model.Cart, error)
	ClearCart(userID uint) (*model.Cart, error)
}

type cartUsecase struct {
	cartRepo     repository.CartRepository
	cartItemRepo repository.CartItemRepository
	productRepo  repository.ProductRepository
}

func NewCartUsecase(
	cartRepo repository.CartRepository,
	cartItemRepo repository.CartItemRepository,
	productRepo repository.ProductRepository,
) CartUsecase {
	return &cartUsecase{cartRepo, cartItemRepo, productRepo}
}

func (u *cartUsecase) GetByUserID(userID uint) (*model.Cart, error) {
	cart, err := u.cartRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return cart, nil
}

func (u *cartUsecase) GetWithFilters(filters map[string]string) ([]model.Cart, error) {
	return u.cartRepo.FindWithFilters(filters)
}

func (u *cartUsecase) AddProduct(userID, productID uint, quantity int) (*model.Cart, error) {
	if quantity <= 0 {
		return nil, errors.New("invalid quantity")
	}

	cart, err := u.cartRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		cart = &model.Cart{UserID: userID, Total: 0}
		if err := u.cartRepo.Create(cart); err != nil {
			return nil, err
		}
	}

	prod, err := u.productRepo.FindByID(productID)
	if err != nil {
		return nil, err
	}
	if prod == nil {
		return nil, gorm.ErrRecordNotFound
	}

	item := &model.CartItem{
		CartID:    cart.ID,
		ProductID: prod.ID,
		Quantity:  quantity,
		UnitPrice: prod.Price,
		Subtotal:  prod.Price * float64(quantity),
	}
	if err := u.cartItemRepo.AddItem(item); err != nil {
		return nil, err
	}

	cart.Total += item.Subtotal
	if err := u.cartRepo.Update(cart); err != nil {
		return nil, err
	}
	return u.cartRepo.FindByUserID(userID)
}

func (u *cartUsecase) UpdateItem(itemID uint, quantity int) (*model.Cart, error) {
	item, err := u.cartItemRepo.FindByID(itemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, gorm.ErrRecordNotFound
	}
	fmt.Println(item.CartID)
	cart, err := u.cartRepo.FindByCartID(item.CartID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, gorm.ErrRecordNotFound
	}

	switch {
	case quantity < 0:
		return nil, errors.New("invalid quantity")
	case quantity == 0:
		cart.Total -= item.Subtotal
		if err := u.cartRepo.Update(cart); err != nil {
			return nil, err
		}
		if err := u.cartItemRepo.DeleteItem(itemID); err != nil {
			return nil, err
		}
	default:
		old := item.Subtotal
		item.Quantity = quantity
		item.Subtotal = item.UnitPrice * float64(quantity)
		if err := u.cartItemRepo.UpdateItem(item); err != nil {
			return nil, err
		}
		cart.Total += item.Subtotal - old
		if err := u.cartRepo.Update(cart); err != nil {
			return nil, err
		}
	}

	return u.cartRepo.FindByUserID(cart.UserID)
}

func (u *cartUsecase) RemoveItem(itemID uint) (*model.Cart, error) {
	item, err := u.cartItemRepo.FindByID(itemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, gorm.ErrRecordNotFound
	}

	cart, err := u.cartRepo.FindByCartID(item.CartID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, gorm.ErrRecordNotFound
	}

	cart.Total -= item.Subtotal
	if err := u.cartRepo.Update(cart); err != nil {
		return nil, err
	}
	if err := u.cartItemRepo.DeleteItem(itemID); err != nil {
		return nil, err
	}

	return u.cartRepo.FindByUserID(cart.UserID)
}

func (u *cartUsecase) ClearCart(userID uint) (*model.Cart, error) {
	cart, err := u.cartRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, gorm.ErrRecordNotFound
	}

	if err := u.cartItemRepo.ClearCart(cart.ID); err != nil {
		return nil, err
	}
	cart.Total = 0
	if err := u.cartRepo.Update(cart); err != nil {
		return nil, err
	}

	return u.cartRepo.FindByUserID(userID)
}

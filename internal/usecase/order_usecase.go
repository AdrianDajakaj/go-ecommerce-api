package usecase

import (
	"errors"
	"fmt"
	"time"

	"go-ecommerce-api/internal/domain/model"
	"go-ecommerce-api/internal/domain/repository"

	"gorm.io/gorm"
)

type OrderUsecase interface {
	GetByID(id uint) (*model.Order, error)
	GetByUserID(userID uint) ([]model.Order, error)
	GetAll() ([]model.Order, error)
	GetWithFilters(filters map[string]string) ([]model.Order, error)
	CreateFromCart(userID uint, paymentMethod model.PaymentMethod, shippingAddressID uint) (*model.Order, error)
	UpdateStatus(id uint, status model.OrderStatus) (*model.Order, error)
	CancelOrder(id uint) (*model.Order, error)
}

type orderUsecase struct {
	orderRepo    repository.OrderRepository
	cartRepo     repository.CartRepository
	cartItemRepo repository.CartItemRepository
	productRepo  repository.ProductRepository
	userRepo     repository.UserRepository
	addressRepo  repository.AddressRepository
}

func NewOrderUsecase(
	orderRepo repository.OrderRepository,
	cartRepo repository.CartRepository,
	cartItemRepo repository.CartItemRepository,
	productRepo repository.ProductRepository,
	userRepo repository.UserRepository,
	addressRepo repository.AddressRepository,
) OrderUsecase {
	return &orderUsecase{
		orderRepo:    orderRepo,
		cartRepo:     cartRepo,
		cartItemRepo: cartItemRepo,
		productRepo:  productRepo,
		userRepo:     userRepo,
		addressRepo:  addressRepo,
	}
}

func (uc *orderUsecase) GetByID(id uint) (*model.Order, error) {
	order, err := uc.orderRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return order, nil
}

func (uc *orderUsecase) GetByUserID(userID uint) ([]model.Order, error) {
	orders, err := uc.orderRepo.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	return orders, nil
}

func (uc *orderUsecase) GetAll() ([]model.Order, error) {
	orders, err := uc.orderRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	return orders, nil
}

func (u *orderUsecase) GetWithFilters(filters map[string]string) ([]model.Order, error) {
	return u.orderRepo.FindWithFilters(filters)
}

func (uc *orderUsecase) CreateFromCart(userID uint, paymentMethod model.PaymentMethod, shippingAddressID uint) (*model.Order, error) {
	cart, err := uc.cartRepo.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}
	if cart == nil || len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	address, err := uc.addressRepo.FindByID(shippingAddressID)
	if err != nil {
		return nil, fmt.Errorf("failed to get address: %w", err)
	}
	if address == nil {
		return nil, errors.New("shipping address not found")
	}

	var orderItems []model.OrderItem
	total := 0.0

	for _, item := range cart.Items {
		product, err := uc.productRepo.FindByID(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to get product: %w", err)
		}
		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("not enough stock for product %s", product.Name)
		}

		product.Stock -= item.Quantity
		if err := uc.productRepo.Update(product); err != nil {
			return nil, fmt.Errorf("failed to update product stock: %w", err)
		}

		orderItems = append(orderItems, model.OrderItem{
			ProductID: product.ID,
			Name:      product.Name,
			UnitPrice: product.Price,
			Quantity:  item.Quantity,
			Subtotal:  product.Price * float64(item.Quantity),
		})
		total += product.Price * float64(item.Quantity)
	}

	order := &model.Order{
		UserID:            userID,
		Status:            model.StatusPending,
		PaymentMethod:     paymentMethod,
		ShippingAddressID: shippingAddressID,
		Items:             orderItems,
		Total:             total,
	}

	if err := uc.orderRepo.Create(order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	if err := uc.cartItemRepo.ClearCart(cart.ID); err != nil {
		return nil, fmt.Errorf("failed to clear cart: %w", err)
	}

	cart.Total = 0
	if err := uc.cartRepo.Update(cart); err != nil {
		return nil, fmt.Errorf("failed to update cart: %w", err)
	}

	return order, nil
}

func (uc *orderUsecase) UpdateStatus(id uint, status model.OrderStatus) (*model.Order, error) {
	order, err := uc.orderRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return nil, gorm.ErrRecordNotFound
	}

	order.Status = status
	if status == model.StatusPaid {
		now := time.Now()
		order.PaidAt = &now
	} else if status == model.StatusShipped {
		now := time.Now()
		order.ShippedAt = &now
	}

	if err := uc.orderRepo.Update(order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}
	return order, nil
}

func (uc *orderUsecase) CancelOrder(id uint) (*model.Order, error) {
	order, err := uc.orderRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return nil, gorm.ErrRecordNotFound
	}

	if order.Status == model.StatusCancelled {
		return order, nil
	}

	for _, item := range order.Items {
		product, err := uc.productRepo.FindByID(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to get product: %w", err)
		}
		if product == nil {
			continue
		}

		product.Stock += item.Quantity
		if err := uc.productRepo.Update(product); err != nil {
			return nil, fmt.Errorf("failed to restore product stock: %w", err)
		}
	}

	order.Status = model.StatusCancelled
	now := time.Now()
	order.CancelledAt = &now

	if err := uc.orderRepo.Update(order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return order, nil
}

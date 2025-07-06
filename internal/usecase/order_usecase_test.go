package usecase

import (
	"errors"
	"testing"
	"time"

	"go-ecommerce-api/internal/domain/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Test constants
const (
	dbError            = "database error"
	dbConnectionFailed = "database connection failed"
	invalidFilter      = "invalid filter"
	productNotFound    = "product not found"
	updateFailed       = "update failed"
	orderUpdateFailed  = "order update failed"
	testProduct1Name   = "Product 1"
	testProduct2Name   = "Product 2"
	testStreet         = "Test St"
	testCity           = "Test City"
	testNumber         = "123"
	modelProduct       = "*model.Product"
	modelOrder         = "*model.Order"
	modelCart          = "*model.Cart"
)

// Mock repositories
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) FindByID(id uint) (*model.Order, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderRepository) FindByUserID(userID uint) ([]model.Order, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Order), args.Error(1)
}

func (m *MockOrderRepository) FindAll() ([]model.Order, error) {
	args := m.Called()
	return args.Get(0).([]model.Order), args.Error(1)
}

func (m *MockOrderRepository) FindWithFilters(filters map[string]string) ([]model.Order, error) {
	args := m.Called(filters)
	return args.Get(0).([]model.Order), args.Error(1)
}

func (m *MockOrderRepository) Create(order *model.Order) error {
	args := m.Called(order)
	// Set ID for created order
	if args.Error(0) == nil && order != nil {
		order.ID = 1
	}
	return args.Error(0)
}

func (m *MockOrderRepository) Update(order *model.Order) error {
	args := m.Called(order)
	// Simulate update operation
	return args.Error(0)
}

type MockCartRepository struct {
	mock.Mock
}

func (m *MockCartRepository) FindByID(id uint) (*model.Cart, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Cart), args.Error(1)
}

func (m *MockCartRepository) FindByUserID(userID uint) (*model.Cart, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Cart), args.Error(1)
}

func (m *MockCartRepository) FindByCartID(cartID uint) (*model.Cart, error) {
	args := m.Called(cartID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Cart), args.Error(1)
}

func (m *MockCartRepository) FindWithFilters(filters map[string]string) ([]model.Cart, error) {
	args := m.Called(filters)
	return args.Get(0).([]model.Cart), args.Error(1)
}

func (m *MockCartRepository) Create(cart *model.Cart) error {
	args := m.Called(cart)
	// Set ID for created cart
	if args.Error(0) == nil && cart != nil {
		cart.ID = 1
	}
	return args.Error(0)
}

func (m *MockCartRepository) Update(cart *model.Cart) error {
	args := m.Called(cart)
	// Simulate update operation
	return args.Error(0)
}

func (m *MockCartRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockCartItemRepository struct {
	mock.Mock
}

func (m *MockCartItemRepository) FindByID(id uint) (*model.CartItem, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CartItem), args.Error(1)
}

func (m *MockCartItemRepository) FindByCartID(cartID uint) ([]model.CartItem, error) {
	args := m.Called(cartID)
	return args.Get(0).([]model.CartItem), args.Error(1)
}

func (m *MockCartItemRepository) FindByCartAndProduct(cartID, productID uint) (*model.CartItem, error) {
	args := m.Called(cartID, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CartItem), args.Error(1)
}

func (m *MockCartItemRepository) AddItem(item *model.CartItem) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockCartItemRepository) Create(item *model.CartItem) error {
	args := m.Called(item)
	// Set ID for created item
	if args.Error(0) == nil && item != nil {
		item.ID = 1
	}
	return args.Error(0)
}

func (m *MockCartItemRepository) UpdateItem(item *model.CartItem) error {
	args := m.Called(item)
	// Simulate update operation for item
	if args.Error(0) == nil && item != nil {
		item.UpdatedAt = time.Now()
	}
	return args.Error(0)
}

func (m *MockCartItemRepository) Update(item *model.CartItem) error {
	args := m.Called(item)
	// Simulate repository update with validation
	if item != nil && args.Error(0) == nil {
		item.UpdatedAt = time.Now()
		return nil
	}
	return args.Error(0)
}

func (m *MockCartItemRepository) DeleteItem(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCartItemRepository) Delete(id uint) error {
	args := m.Called(id)
	// Simulate repository delete with validation
	if id > 0 && args.Error(0) == nil {
		return nil
	}
	return args.Error(0)
}

func (m *MockCartItemRepository) ClearCart(cartID uint) error {
	args := m.Called(cartID)
	return args.Error(0)
}

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) FindByID(id uint) (*model.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Product), args.Error(1)
}

func (m *MockProductRepository) FindAll() ([]model.Product, error) {
	args := m.Called()
	return args.Get(0).([]model.Product), args.Error(1)
}

func (m *MockProductRepository) FindWithFilters(filters map[string]string) ([]model.Product, error) {
	args := m.Called(filters)
	return args.Get(0).([]model.Product), args.Error(1)
}

func (m *MockProductRepository) FindByCategory(categoryID uint) ([]model.Product, error) {
	args := m.Called(categoryID)
	return args.Get(0).([]model.Product), args.Error(1)
}

func (m *MockProductRepository) Create(product *model.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) Update(product *model.Product) error {
	args := m.Called(product)
	// Simulate update operation for product
	if args.Error(0) == nil && product != nil {
		product.UpdatedAt = time.Now()
	}
	return args.Error(0)
}

func (m *MockProductRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindAll() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) FindWithFilters(filters map[string]string) ([]model.User, error) {
	args := m.Called(filters)
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(user *model.User) error {
	args := m.Called(user)
	// Simulate update operation for user
	if args.Error(0) == nil && user != nil {
		user.UpdatedAt = time.Now()
	}
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockAddressRepository struct {
	mock.Mock
}

func (m *MockAddressRepository) FindByID(id uint) (*model.Address, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Address), args.Error(1)
}

func (m *MockAddressRepository) Create(address *model.Address) error {
	args := m.Called(address)
	return args.Error(0)
}

func (m *MockAddressRepository) Update(address *model.Address) error {
	args := m.Called(address)
	// Simulate update operation for address
	if args.Error(0) == nil && address != nil {
		address.UpdatedAt = time.Now()
	}
	return args.Error(0)
}

func (m *MockAddressRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupOrderUsecase() (*orderUsecase, *MockOrderRepository, *MockCartRepository, *MockCartItemRepository, *MockProductRepository, *MockUserRepository, *MockAddressRepository) {
	mockOrderRepo := new(MockOrderRepository)
	mockCartRepo := new(MockCartRepository)
	mockCartItemRepo := new(MockCartItemRepository)
	mockProductRepo := new(MockProductRepository)
	mockUserRepo := new(MockUserRepository)
	mockAddressRepo := new(MockAddressRepository)

	uc := &orderUsecase{
		orderRepo:    mockOrderRepo,
		cartRepo:     mockCartRepo,
		cartItemRepo: mockCartItemRepo,
		productRepo:  mockProductRepo,
		userRepo:     mockUserRepo,
		addressRepo:  mockAddressRepo,
	}

	return uc, mockOrderRepo, mockCartRepo, mockCartItemRepo, mockProductRepo, mockUserRepo, mockAddressRepo
}

func TestNewOrderUsecase(t *testing.T) {
	mockOrderRepo := new(MockOrderRepository)
	mockCartRepo := new(MockCartRepository)
	mockCartItemRepo := new(MockCartItemRepository)
	mockProductRepo := new(MockProductRepository)
	mockUserRepo := new(MockUserRepository)
	mockAddressRepo := new(MockAddressRepository)

	uc := NewOrderUsecase(mockOrderRepo, mockCartRepo, mockCartItemRepo, mockProductRepo, mockUserRepo, mockAddressRepo)

	// Assertion 94: NewOrderUsecase should return a non-nil usecase instance
	assert.NotNil(t, uc)
	// Assertion 95: NewOrderUsecase should return an OrderUsecase interface
	assert.Implements(t, (*OrderUsecase)(nil), uc)
}

func TestOrderUsecaseGetByIDSuccess(t *testing.T) {
	uc, mockOrderRepo, _, _, _, _, _ := setupOrderUsecase()

	expectedOrder := &model.Order{
		ID:     1,
		UserID: 1,
		Status: model.StatusPending,
		Total:  100.0,
	}

	mockOrderRepo.On("FindByID", uint(1)).Return(expectedOrder, nil)

	result, err := uc.GetByID(1)

	// Assertion 96: GetByID should not return an error for valid order ID
	assert.NoError(t, err)
	// Assertion 97: GetByID should return the expected order
	assert.Equal(t, expectedOrder, result)
	// Assertion 98: GetByID should return an order with correct ID
	assert.Equal(t, uint(1), result.ID)
	// Assertion 99: GetByID should return an order with correct UserID
	assert.Equal(t, uint(1), result.UserID)
	// Assertion 100: GetByID should return an order with correct Status
	assert.Equal(t, model.StatusPending, result.Status)
	// Assertion 101: GetByID should return an order with correct Total
	assert.Equal(t, 100.0, result.Total)

	mockOrderRepo.AssertExpectations(t)
}

func TestOrderUsecaseGetByIDNotFound(t *testing.T) {
	uc, mockOrderRepo, _, _, _, _, _ := setupOrderUsecase()

	mockOrderRepo.On("FindByID", uint(999)).Return(nil, nil)

	result, err := uc.GetByID(999)

	// Assertion 102: GetByID should return gorm.ErrRecordNotFound for non-existent order
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	// Assertion 103: GetByID should return nil order for non-existent order
	assert.Nil(t, result)

	mockOrderRepo.AssertExpectations(t)
}

func TestOrderUsecaseGetByIDRepositoryError(t *testing.T) {
	uc, mockOrderRepo, _, _, _, _, _ := setupOrderUsecase()

	mockOrderRepo.On("FindByID", uint(1)).Return(nil, errors.New(dbError))

	result, err := uc.GetByID(1)

	// Assertion 104: GetByID should return error when repository fails
	assert.Error(t, err)
	// Assertion 105: GetByID should return nil order when repository fails
	assert.Nil(t, result)
	// Assertion 106: GetByID should wrap repository error message
	assert.Contains(t, err.Error(), "failed to get order")

	mockOrderRepo.AssertExpectations(t)
}

func TestOrderUsecaseGetByUserIDSuccess(t *testing.T) {
	uc, mockOrderRepo, _, _, _, _, _ := setupOrderUsecase()

	expectedOrders := []model.Order{
		{ID: 1, UserID: 1, Status: model.StatusPending, Total: 100.0},
		{ID: 2, UserID: 1, Status: model.StatusPaid, Total: 200.0},
	}

	mockOrderRepo.On("FindByUserID", uint(1)).Return(expectedOrders, nil)

	result, err := uc.GetByUserID(1)

	// Assertion 107: GetByUserID should not return an error for valid user ID
	assert.NoError(t, err)
	// Assertion 108: GetByUserID should return the expected orders slice
	assert.Equal(t, expectedOrders, result)
	// Assertion 109: GetByUserID should return correct number of orders
	assert.Len(t, result, 2)
	// Assertion 110: GetByUserID should return orders with correct user ID
	assert.Equal(t, uint(1), result[0].UserID)
	// Assertion 111: GetByUserID should return orders with correct user ID for second order
	assert.Equal(t, uint(1), result[1].UserID)

	mockOrderRepo.AssertExpectations(t)
}

func TestOrderUsecaseGetByUserIDEmptyResult(t *testing.T) {
	uc, mockOrderRepo, _, _, _, _, _ := setupOrderUsecase()

	mockOrderRepo.On("FindByUserID", uint(999)).Return([]model.Order{}, nil)

	result, err := uc.GetByUserID(999)

	// Assertion 112: GetByUserID should not return an error for user with no orders
	assert.NoError(t, err)
	// Assertion 113: GetByUserID should return empty slice for user with no orders
	assert.Empty(t, result)
	// Assertion 114: GetByUserID should return slice with zero length for user with no orders
	assert.Len(t, result, 0)

	mockOrderRepo.AssertExpectations(t)
}

func TestOrderUsecaseGetByUserIDRepositoryError(t *testing.T) {
	uc, mockOrderRepo, _, _, _, _, _ := setupOrderUsecase()

	mockOrderRepo.On("FindByUserID", uint(1)).Return([]model.Order{}, errors.New(dbError))

	result, err := uc.GetByUserID(1)

	// Assertion 115: GetByUserID should return error when repository fails
	assert.Error(t, err)
	// Assertion 116: GetByUserID should return empty slice when repository fails
	assert.Empty(t, result)
	// Assertion 117: GetByUserID should wrap repository error message
	assert.Contains(t, err.Error(), "failed to get orders")

	mockOrderRepo.AssertExpectations(t)
}

func TestOrderUsecaseGetAllSuccess(t *testing.T) {
	uc, mockOrderRepo, _, _, _, _, _ := setupOrderUsecase()

	expectedOrders := []model.Order{
		{ID: 1, UserID: 1, Status: model.StatusPending, Total: 100.0},
		{ID: 2, UserID: 2, Status: model.StatusPaid, Total: 200.0},
		{ID: 3, UserID: 3, Status: model.StatusShipped, Total: 300.0},
	}

	mockOrderRepo.On("FindAll").Return(expectedOrders, nil)

	result, err := uc.GetAll()

	// Assertion 118: GetAll should not return an error when repository succeeds
	assert.NoError(t, err)
	// Assertion 119: GetAll should return the expected orders slice
	assert.Equal(t, expectedOrders, result)
	// Assertion 120: GetAll should return correct number of orders
	assert.Len(t, result, 3)
	// Assertion 121: GetAll should return orders with different user IDs
	assert.NotEqual(t, result[0].UserID, result[1].UserID)

	mockOrderRepo.AssertExpectations(t)
}

func TestOrderUsecaseGetAllEmptyResult(t *testing.T) {
	uc, mockOrderRepo, _, _, _, _, _ := setupOrderUsecase()

	mockOrderRepo.On("FindAll").Return([]model.Order{}, nil)

	result, err := uc.GetAll()

	// Assertion 122: GetAll should not return an error when no orders exist
	assert.NoError(t, err)
	// Assertion 123: GetAll should return empty slice when no orders exist
	assert.Empty(t, result)
	// Assertion 124: GetAll should return slice with zero length when no orders exist
	assert.Len(t, result, 0)

	mockOrderRepo.AssertExpectations(t)
}

func TestOrderUsecaseGetAllRepositoryError(t *testing.T) {
	uc, mockOrderRepo, _, _, _, _, _ := setupOrderUsecase()

	mockOrderRepo.On("FindAll").Return([]model.Order{}, errors.New(dbConnectionFailed))

	result, err := uc.GetAll()

	// Assertion 125: GetAll should return error when repository fails
	assert.Error(t, err)
	// Assertion 126: GetAll should return empty slice when repository fails
	assert.Empty(t, result)
	// Assertion 127: GetAll should wrap repository error message
	assert.Contains(t, err.Error(), "failed to get orders")

	mockOrderRepo.AssertExpectations(t)
}

func TestOrderUsecaseGetWithFiltersSuccess(t *testing.T) {
	uc, mockOrderRepo, _, _, _, _, _ := setupOrderUsecase()

	filters := map[string]string{
		"status":  "PENDING",
		"user_id": "1",
	}

	expectedOrders := []model.Order{
		{ID: 1, UserID: 1, Status: model.StatusPending, Total: 100.0},
	}

	mockOrderRepo.On("FindWithFilters", filters).Return(expectedOrders, nil)

	result, err := uc.GetWithFilters(filters)

	// Assertion 128: GetWithFilters should not return an error when repository succeeds
	assert.NoError(t, err)
	// Assertion 129: GetWithFilters should return the expected filtered orders
	assert.Equal(t, expectedOrders, result)
	// Assertion 130: GetWithFilters should return correct number of filtered orders
	assert.Len(t, result, 1)
	// Assertion 131: GetWithFilters should return orders matching the filter status
	assert.Equal(t, model.StatusPending, result[0].Status)

	mockOrderRepo.AssertExpectations(t)
}

func TestOrderUsecaseGetWithFiltersRepositoryError(t *testing.T) {
	uc, mockOrderRepo, _, _, _, _, _ := setupOrderUsecase()

	filters := map[string]string{"status": "INVALID"}

	mockOrderRepo.On("FindWithFilters", filters).Return([]model.Order{}, errors.New(invalidFilter))

	result, err := uc.GetWithFilters(filters)

	// Assertion 132: GetWithFilters should return error when repository fails
	assert.Error(t, err)
	// Assertion 133: GetWithFilters should return empty slice when repository fails
	assert.Empty(t, result)

	mockOrderRepo.AssertExpectations(t)
}

func TestOrderUsecaseCreateFromCartSuccess(t *testing.T) {
	uc, mockOrderRepo, mockCartRepo, mockCartItemRepo, mockProductRepo, _, mockAddressRepo := setupOrderUsecase()

	cart := &model.Cart{
		ID:     1,
		UserID: 1,
		Total:  150.0,
		Items: []model.CartItem{
			{ID: 1, CartID: 1, ProductID: 1, Quantity: 2},
			{ID: 2, CartID: 1, ProductID: 2, Quantity: 1},
		},
	}

	product1 := &model.Product{ID: 1, Name: testProduct1Name, Price: 50.0, Stock: 10}
	product2 := &model.Product{ID: 2, Name: testProduct2Name, Price: 50.0, Stock: 5}

	address := &model.Address{ID: 1, Street: testStreet, Number: testNumber, City: testCity}

	mockCartRepo.On("FindByUserID", uint(1)).Return(cart, nil)
	mockAddressRepo.On("FindByID", uint(1)).Return(address, nil)
	mockProductRepo.On("FindByID", uint(1)).Return(product1, nil)
	mockProductRepo.On("FindByID", uint(2)).Return(product2, nil)
	mockProductRepo.On("Update", mock.AnythingOfType(modelProduct)).Return(nil).Twice()
	mockOrderRepo.On("Create", mock.AnythingOfType(modelOrder)).Return(nil)
	mockCartItemRepo.On("ClearCart", uint(1)).Return(nil)
	mockCartRepo.On("Update", mock.AnythingOfType(modelCart)).Return(nil)

	result, err := uc.CreateFromCart(1, model.PaymentCard, 1)

	// Assertion 134: CreateFromCart should not return an error for valid cart and address
	assert.NoError(t, err)
	// Assertion 135: CreateFromCart should return a non-nil order
	assert.NotNil(t, result)
	// Assertion 136: CreateFromCart should set correct user ID on order
	assert.Equal(t, uint(1), result.UserID)
	// Assertion 137: CreateFromCart should set pending status on new order
	assert.Equal(t, model.StatusPending, result.Status)
	// Assertion 138: CreateFromCart should set correct payment method on order
	assert.Equal(t, model.PaymentCard, result.PaymentMethod)
	// Assertion 139: CreateFromCart should set correct shipping address ID on order
	assert.Equal(t, uint(1), result.ShippingAddressID)
	// Assertion 140: CreateFromCart should calculate correct total for order
	assert.Equal(t, 150.0, result.Total)
	// Assertion 141: CreateFromCart should create correct number of order items
	assert.Len(t, result.Items, 2)
	// Assertion 142: CreateFromCart should set correct product ID on first order item
	assert.Equal(t, uint(1), result.Items[0].ProductID)
	// Assertion 143: CreateFromCart should set correct quantity on first order item
	assert.Equal(t, 2, result.Items[0].Quantity)
	// Assertion 144: CreateFromCart should set correct subtotal on first order item
	assert.Equal(t, 100.0, result.Items[0].Subtotal)

	mockOrderRepo.AssertExpectations(t)
	mockCartRepo.AssertExpectations(t)
	mockCartItemRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
	mockAddressRepo.AssertExpectations(t)
}

func TestOrderUsecaseCreateFromCartEmptyCart(t *testing.T) {
	uc, _, mockCartRepo, _, _, _, _ := setupOrderUsecase()

	cart := &model.Cart{
		ID:     1,
		UserID: 1,
		Items:  []model.CartItem{},
	}

	mockCartRepo.On("FindByUserID", uint(1)).Return(cart, nil)

	result, err := uc.CreateFromCart(1, model.PaymentCard, 1)

	// Assertion 145: CreateFromCart should return error for empty cart
	assert.Error(t, err)
	// Assertion 146: CreateFromCart should return nil order for empty cart
	assert.Nil(t, result)
	// Assertion 147: CreateFromCart should return appropriate error message for empty cart
	assert.Contains(t, err.Error(), "cart is empty")

	mockCartRepo.AssertExpectations(t)
}

func TestOrderUsecaseCreateFromCartCartNotFound(t *testing.T) {
	uc, _, mockCartRepo, _, _, _, _ := setupOrderUsecase()

	mockCartRepo.On("FindByUserID", uint(999)).Return(nil, nil)

	result, err := uc.CreateFromCart(999, model.PaymentCard, 1)

	// Assertion 148: CreateFromCart should return error when cart not found
	assert.Error(t, err)
	// Assertion 149: CreateFromCart should return nil order when cart not found
	assert.Nil(t, result)
	// Assertion 150: CreateFromCart should return appropriate error message when cart not found
	assert.Contains(t, err.Error(), "cart is empty")

	mockCartRepo.AssertExpectations(t)
}

func TestOrderUsecaseCreateFromCartAddressNotFound(t *testing.T) {
	uc, _, mockCartRepo, _, _, _, mockAddressRepo := setupOrderUsecase()

	cart := &model.Cart{
		ID:     1,
		UserID: 1,
		Items: []model.CartItem{
			{ID: 1, CartID: 1, ProductID: 1, Quantity: 1},
		},
	}

	mockCartRepo.On("FindByUserID", uint(1)).Return(cart, nil)
	mockAddressRepo.On("FindByID", uint(999)).Return(nil, nil)

	result, err := uc.CreateFromCart(1, model.PaymentCard, 999)

	// Assertion 151: CreateFromCart should return error when shipping address not found
	assert.Error(t, err)
	// Assertion 152: CreateFromCart should return nil order when shipping address not found
	assert.Nil(t, result)
	// Assertion 153: CreateFromCart should return appropriate error message when shipping address not found
	assert.Contains(t, err.Error(), "shipping address not found")

	mockCartRepo.AssertExpectations(t)
	mockAddressRepo.AssertExpectations(t)
}

func TestOrderUsecaseCreateFromCartInsufficientStock(t *testing.T) {
	uc, _, mockCartRepo, _, mockProductRepo, _, mockAddressRepo := setupOrderUsecase()

	cart := &model.Cart{
		ID:     1,
		UserID: 1,
		Items: []model.CartItem{
			{ID: 1, CartID: 1, ProductID: 1, Quantity: 10},
		},
	}

	product := &model.Product{ID: 1, Name: testProduct1Name, Price: 50.0, Stock: 5}
	address := &model.Address{ID: 1, Street: testStreet, Number: testNumber, City: testCity}

	mockCartRepo.On("FindByUserID", uint(1)).Return(cart, nil)
	mockAddressRepo.On("FindByID", uint(1)).Return(address, nil)
	mockProductRepo.On("FindByID", uint(1)).Return(product, nil)

	result, err := uc.CreateFromCart(1, model.PaymentCard, 1)

	// Assertion 154: CreateFromCart should return error when insufficient stock
	assert.Error(t, err)
	// Assertion 155: CreateFromCart should return nil order when insufficient stock
	assert.Nil(t, result)
	// Assertion 156: CreateFromCart should return appropriate error message when insufficient stock
	assert.Contains(t, err.Error(), "not enough stock for product")

	mockCartRepo.AssertExpectations(t)
	mockAddressRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestOrderUsecaseCreateFromCartProductNotFound(t *testing.T) {
	uc, _, mockCartRepo, _, mockProductRepo, _, mockAddressRepo := setupOrderUsecase()

	cart := &model.Cart{
		ID:     1,
		UserID: 1,
		Items: []model.CartItem{
			{ID: 1, CartID: 1, ProductID: 999, Quantity: 1},
		},
	}

	address := &model.Address{ID: 1, Street: testStreet, Number: testNumber, City: testCity}

	mockCartRepo.On("FindByUserID", uint(1)).Return(cart, nil)
	mockAddressRepo.On("FindByID", uint(1)).Return(address, nil)
	mockProductRepo.On("FindByID", uint(999)).Return(nil, errors.New(productNotFound))

	result, err := uc.CreateFromCart(1, model.PaymentCard, 1)

	// Assertion 157: CreateFromCart should return error when product not found
	assert.Error(t, err)
	// Assertion 158: CreateFromCart should return nil order when product not found
	assert.Nil(t, result)
	// Assertion 159: CreateFromCart should wrap product repository error
	assert.Contains(t, err.Error(), "failed to get product")

	mockCartRepo.AssertExpectations(t)
	mockAddressRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestOrderUsecaseUpdateStatusSuccess(t *testing.T) {
	uc, mockOrderRepo, _, _, _, _, _ := setupOrderUsecase()

	order := &model.Order{
		ID:     1,
		UserID: 1,
		Status: model.StatusPending,
		Total:  100.0,
	}

	mockOrderRepo.On("FindByID", uint(1)).Return(order, nil)
	mockOrderRepo.On("Update", mock.AnythingOfType(modelOrder)).Return(nil)

	result, err := uc.UpdateStatus(1, model.StatusPaid)

	// Assertion 160: UpdateStatus should not return an error for valid order and status
	assert.NoError(t, err)
	// Assertion 161: UpdateStatus should return updated order
	assert.NotNil(t, result)
	// Assertion 162: UpdateStatus should update order status correctly
	assert.Equal(t, model.StatusPaid, result.Status)
	// Assertion 163: UpdateStatus should set PaidAt timestamp when status is PAID
	assert.NotNil(t, result.PaidAt)
	// Assertion 164: UpdateStatus should set PaidAt to current time approximately
	assert.WithinDuration(t, time.Now(), *result.PaidAt, time.Second)

	mockOrderRepo.AssertExpectations(t)
}

func TestOrderUsecaseUpdateStatusToShipped(t *testing.T) {
	uc, mockOrderRepo, _, _, _, _, _ := setupOrderUsecase()

	order := &model.Order{
		ID:     1,
		UserID: 1,
		Status: model.StatusPaid,
		Total:  100.0,
	}

	mockOrderRepo.On("FindByID", uint(1)).Return(order, nil)
	mockOrderRepo.On("Update", mock.AnythingOfType(modelOrder)).Return(nil)

	result, err := uc.UpdateStatus(1, model.StatusShipped)

	// Assertion 165: UpdateStatus should not return an error when updating to shipped
	assert.NoError(t, err)
	// Assertion 166: UpdateStatus should update status to shipped correctly
	assert.Equal(t, model.StatusShipped, result.Status)
	// Assertion 167: UpdateStatus should set ShippedAt timestamp when status is SHIPPED
	assert.NotNil(t, result.ShippedAt)
	// Assertion 168: UpdateStatus should set ShippedAt to current time approximately
	assert.WithinDuration(t, time.Now(), *result.ShippedAt, time.Second)

	mockOrderRepo.AssertExpectations(t)
}

func TestOrderUsecaseUpdateStatusOrderNotFound(t *testing.T) {
	uc, mockOrderRepo, _, _, _, _, _ := setupOrderUsecase()

	mockOrderRepo.On("FindByID", uint(999)).Return(nil, nil)

	result, err := uc.UpdateStatus(999, model.StatusPaid)

	// Assertion 169: UpdateStatus should return gorm.ErrRecordNotFound for non-existent order
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	// Assertion 170: UpdateStatus should return nil order for non-existent order
	assert.Nil(t, result)

	mockOrderRepo.AssertExpectations(t)
}

func TestOrderUsecaseUpdateStatusRepositoryError(t *testing.T) {
	uc, mockOrderRepo, _, _, _, _, _ := setupOrderUsecase()

	mockOrderRepo.On("FindByID", uint(1)).Return(nil, errors.New(dbError))

	result, err := uc.UpdateStatus(1, model.StatusPaid)

	// Assertion 171: UpdateStatus should return error when repository fails to find order
	assert.Error(t, err)
	// Assertion 172: UpdateStatus should return nil order when repository fails to find order
	assert.Nil(t, result)
	// Assertion 173: UpdateStatus should wrap repository error message
	assert.Contains(t, err.Error(), "failed to get order")

	mockOrderRepo.AssertExpectations(t)
}

func TestOrderUsecaseUpdateStatusUpdateRepositoryError(t *testing.T) {
	uc, mockOrderRepo, _, _, _, _, _ := setupOrderUsecase()

	order := &model.Order{
		ID:     1,
		UserID: 1,
		Status: model.StatusPending,
		Total:  100.0,
	}

	mockOrderRepo.On("FindByID", uint(1)).Return(order, nil)
	mockOrderRepo.On("Update", mock.AnythingOfType(modelOrder)).Return(errors.New(updateFailed))

	result, err := uc.UpdateStatus(1, model.StatusPaid)

	// Assertion 174: UpdateStatus should return error when repository fails to update
	assert.Error(t, err)
	// Assertion 175: UpdateStatus should return nil order when repository fails to update
	assert.Nil(t, result)
	// Assertion 176: UpdateStatus should wrap update repository error message
	assert.Contains(t, err.Error(), "failed to update order")

	mockOrderRepo.AssertExpectations(t)
}

func TestOrderUsecaseCancelOrderSuccess(t *testing.T) {
	uc, mockOrderRepo, _, _, mockProductRepo, _, _ := setupOrderUsecase()

	order := &model.Order{
		ID:     1,
		UserID: 1,
		Status: model.StatusPending,
		Items: []model.OrderItem{
			{ID: 1, ProductID: 1, Quantity: 2, UnitPrice: 50.0, Subtotal: 100.0},
			{ID: 2, ProductID: 2, Quantity: 1, UnitPrice: 30.0, Subtotal: 30.0},
		},
		Total: 130.0,
	}

	product1 := &model.Product{ID: 1, Name: testProduct1Name, Price: 50.0, Stock: 8}
	product2 := &model.Product{ID: 2, Name: "Product 2", Price: 30.0, Stock: 4}

	mockOrderRepo.On("FindByID", uint(1)).Return(order, nil)
	mockProductRepo.On("FindByID", uint(1)).Return(product1, nil)
	mockProductRepo.On("FindByID", uint(2)).Return(product2, nil)
	mockProductRepo.On("Update", mock.AnythingOfType(modelProduct)).Return(nil).Twice()
	mockOrderRepo.On("Update", mock.AnythingOfType(modelOrder)).Return(nil)

	result, err := uc.CancelOrder(1)

	// Assertion 177: CancelOrder should not return an error for valid order
	assert.NoError(t, err)
	// Assertion 178: CancelOrder should return updated order
	assert.NotNil(t, result)
	// Assertion 179: CancelOrder should update order status to cancelled
	assert.Equal(t, model.StatusCancelled, result.Status)
	// Assertion 180: CancelOrder should set CancelledAt timestamp
	assert.NotNil(t, result.CancelledAt)
	// Assertion 181: CancelOrder should set CancelledAt to current time approximately
	assert.WithinDuration(t, time.Now(), *result.CancelledAt, time.Second)

	mockOrderRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestOrderUsecaseCancelOrderAlreadyCancelled(t *testing.T) {
	uc, mockOrderRepo, _, _, _, _, _ := setupOrderUsecase()

	cancelledTime := time.Now().Add(-time.Hour)
	order := &model.Order{
		ID:          1,
		UserID:      1,
		Status:      model.StatusCancelled,
		CancelledAt: &cancelledTime,
		Total:       100.0,
	}

	mockOrderRepo.On("FindByID", uint(1)).Return(order, nil)

	result, err := uc.CancelOrder(1)

	// Assertion 182: CancelOrder should not return an error for already cancelled order
	assert.NoError(t, err)
	// Assertion 183: CancelOrder should return the order for already cancelled order
	assert.NotNil(t, result)
	// Assertion 184: CancelOrder should maintain cancelled status for already cancelled order
	assert.Equal(t, model.StatusCancelled, result.Status)
	// Assertion 185: CancelOrder should maintain original cancelled timestamp
	assert.Equal(t, &cancelledTime, result.CancelledAt)

	mockOrderRepo.AssertExpectations(t)
}

func TestOrderUsecaseCancelOrderOrderNotFound(t *testing.T) {
	uc, mockOrderRepo, _, _, _, _, _ := setupOrderUsecase()

	mockOrderRepo.On("FindByID", uint(999)).Return(nil, nil)

	result, err := uc.CancelOrder(999)

	// Assertion 186: CancelOrder should return gorm.ErrRecordNotFound for non-existent order
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	// Assertion 187: CancelOrder should return nil order for non-existent order
	assert.Nil(t, result)

	mockOrderRepo.AssertExpectations(t)
}

func TestOrderUsecaseCancelOrderProductNotFound(t *testing.T) {
	uc, mockOrderRepo, _, _, mockProductRepo, _, _ := setupOrderUsecase()

	order := &model.Order{
		ID:     1,
		UserID: 1,
		Status: model.StatusPending,
		Items: []model.OrderItem{
			{ID: 1, ProductID: 999, Quantity: 2, UnitPrice: 50.0, Subtotal: 100.0},
		},
		Total: 100.0,
	}

	mockOrderRepo.On("FindByID", uint(1)).Return(order, nil)
	mockProductRepo.On("FindByID", uint(999)).Return(nil, errors.New(productNotFound))

	result, err := uc.CancelOrder(1)

	// Assertion 188: CancelOrder should return error when product not found during cancellation
	assert.Error(t, err)
	// Assertion 189: CancelOrder should return nil order when product not found during cancellation
	assert.Nil(t, result)
	// Assertion 190: CancelOrder should wrap product repository error during cancellation
	assert.Contains(t, err.Error(), "failed to get product")

	mockOrderRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestOrderUsecaseCancelOrderProductUpdateError(t *testing.T) {
	uc, mockOrderRepo, _, _, mockProductRepo, _, _ := setupOrderUsecase()

	order := &model.Order{
		ID:     1,
		UserID: 1,
		Status: model.StatusPending,
		Items: []model.OrderItem{
			{ID: 1, ProductID: 1, Quantity: 2, UnitPrice: 50.0, Subtotal: 100.0},
		},
		Total: 100.0,
	}

	product := &model.Product{ID: 1, Name: testProduct1Name, Price: 50.0, Stock: 8}

	mockOrderRepo.On("FindByID", uint(1)).Return(order, nil)
	mockProductRepo.On("FindByID", uint(1)).Return(product, nil)
	mockProductRepo.On("Update", mock.AnythingOfType(modelProduct)).Return(errors.New(updateFailed))

	result, err := uc.CancelOrder(1)

	// Assertion 191: CancelOrder should return error when product update fails during cancellation
	assert.Error(t, err)
	// Assertion 192: CancelOrder should return nil order when product update fails during cancellation
	assert.Nil(t, result)
	// Assertion 193: CancelOrder should wrap product update error during cancellation
	assert.Contains(t, err.Error(), "failed to restore product stock")

	mockOrderRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestOrderUsecaseCancelOrderOrderUpdateError(t *testing.T) {
	uc, mockOrderRepo, _, _, mockProductRepo, _, _ := setupOrderUsecase()

	order := &model.Order{
		ID:     1,
		UserID: 1,
		Status: model.StatusPending,
		Items: []model.OrderItem{
			{ID: 1, ProductID: 1, Quantity: 2, UnitPrice: 50.0, Subtotal: 100.0},
		},
		Total: 100.0,
	}

	product := &model.Product{ID: 1, Name: testProduct1Name, Price: 50.0, Stock: 8}

	mockOrderRepo.On("FindByID", uint(1)).Return(order, nil)
	mockProductRepo.On("FindByID", uint(1)).Return(product, nil)
	mockProductRepo.On("Update", mock.AnythingOfType(modelProduct)).Return(nil)
	mockOrderRepo.On("Update", mock.AnythingOfType(modelOrder)).Return(errors.New(orderUpdateFailed))

	result, err := uc.CancelOrder(1)

	// Assertion 194: CancelOrder should return error when order update fails
	assert.Error(t, err)
	// Assertion 195: CancelOrder should return nil order when order update fails
	assert.Nil(t, result)
	// Assertion 196: CancelOrder should wrap order update error message
	assert.Contains(t, err.Error(), "failed to update order")

	mockOrderRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestOrderUsecaseCancelOrderProductNilHandling(t *testing.T) {
	uc, mockOrderRepo, _, _, mockProductRepo, _, _ := setupOrderUsecase()

	order := &model.Order{
		ID:     1,
		UserID: 1,
		Status: model.StatusPending,
		Items: []model.OrderItem{
			{ID: 1, ProductID: 1, Quantity: 2, UnitPrice: 50.0, Subtotal: 100.0},
		},
		Total: 100.0,
	}

	mockOrderRepo.On("FindByID", uint(1)).Return(order, nil)
	mockProductRepo.On("FindByID", uint(1)).Return(nil, nil)
	mockOrderRepo.On("Update", mock.AnythingOfType(modelOrder)).Return(nil)

	result, err := uc.CancelOrder(1)

	// Assertion 197: CancelOrder should not return an error when product is nil (deleted product)
	assert.NoError(t, err)
	// Assertion 198: CancelOrder should return updated order when product is nil (deleted product)
	assert.NotNil(t, result)
	// Assertion 199: CancelOrder should update order status to cancelled when product is nil
	assert.Equal(t, model.StatusCancelled, result.Status)
	// Assertion 200: CancelOrder should set CancelledAt timestamp when product is nil
	assert.NotNil(t, result.CancelledAt)

	mockOrderRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

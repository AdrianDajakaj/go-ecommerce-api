package usecase

import (
	"testing"

	"go-ecommerce-api/internal/domain/model"

	"gorm.io/gorm"
)

// Mock repositories for cart testing
type mockCartRepository struct {
	carts  []model.Cart
	nextID uint
}

type mockCartItemRepository struct {
	items  []model.CartItem
	nextID uint
}

func newMockCartRepository() *mockCartRepository {
	return &mockCartRepository{
		carts:  []model.Cart{},
		nextID: 1,
	}
}

func newMockCartItemRepository() *mockCartItemRepository {
	return &mockCartItemRepository{
		items:  []model.CartItem{},
		nextID: 1,
	}
}

func (m *mockCartRepository) FindByUserID(userID uint) (*model.Cart, error) {
	for _, cart := range m.carts {
		if cart.UserID == userID {
			return &cart, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockCartRepository) FindByCartID(cartID uint) (*model.Cart, error) {
	for _, cart := range m.carts {
		if cart.ID == cartID {
			return &cart, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockCartRepository) FindWithFilters(filters map[string]string) ([]model.Cart, error) {
	return m.carts, nil
}

func (m *mockCartRepository) Create(cart *model.Cart) error {
	cart.ID = m.nextID
	m.nextID++
	m.carts = append(m.carts, *cart)
	return nil
}

func (m *mockCartRepository) Update(cart *model.Cart) error {
	for i, c := range m.carts {
		if c.ID == cart.ID {
			m.carts[i] = *cart
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *mockCartRepository) Delete(id uint) error {
	for i, c := range m.carts {
		if c.ID == id {
			m.carts = append(m.carts[:i], m.carts[i+1:]...)
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *mockCartItemRepository) FindByID(id uint) (*model.CartItem, error) {
	for _, item := range m.items {
		if item.ID == id {
			return &item, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockCartItemRepository) FindByCartID(cartID uint) ([]model.CartItem, error) {
	var result []model.CartItem
	for _, item := range m.items {
		if item.CartID == cartID {
			result = append(result, item)
		}
	}
	return result, nil
}

func (m *mockCartItemRepository) AddItem(item *model.CartItem) error {
	item.ID = m.nextID
	m.nextID++
	m.items = append(m.items, *item)
	return nil
}

func (m *mockCartItemRepository) UpdateItem(item *model.CartItem) error {
	for i, it := range m.items {
		if it.ID == item.ID {
			m.items[i] = *item
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *mockCartItemRepository) DeleteItem(id uint) error {
	for i, it := range m.items {
		if it.ID == id {
			m.items = append(m.items[:i], m.items[i+1:]...)
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *mockCartItemRepository) ClearCart(cartID uint) error {
	var newItems []model.CartItem
	for _, item := range m.items {
		if item.CartID != cartID {
			newItems = append(newItems, item)
		}
	}
	m.items = newItems
	return nil
}

func TestCartUsecaseGetByUserID(t *testing.T) {
	cartRepo := newMockCartRepository()
	cartItemRepo := newMockCartItemRepository()
	productRepo := newMockProductRepository()
	usecase := NewCartUsecase(cartRepo, cartItemRepo, productRepo)

	// Test Case 17: Get cart for non-existent user
	cart, err := usecase.GetByUserID(999)
	// Assertion 61: Should return ErrRecordNotFound for non-existent user cart
	if err != gorm.ErrRecordNotFound {
		t.Errorf("Expected gorm.ErrRecordNotFound, got %v", err)
	}
	// Assertion 62: Cart should be nil when not found
	if cart != nil {
		t.Errorf("Expected nil cart, got %v", cart)
	}

	// Create test cart
	testCart := &model.Cart{
		UserID: 1,
		Total:  0.0,
	}
	cartRepo.Create(testCart)

	// Test Case 18: Get existing user cart
	cart, err = usecase.GetByUserID(1)
	// Assertion 63: No error should occur when getting existing cart
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	// Assertion 64: Cart should not be nil when found
	if cart == nil {
		t.Error("Expected cart, got nil")
		return
	}
	// Assertion 65: Cart UserID should match requested user
	if cart.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", cart.UserID)
	}
	// Assertion 66: Cart should have correct total
	if cart.Total != 0.0 {
		t.Errorf("Expected total 0.0, got %f", cart.Total)
	}
}

func TestCartUsecaseGetWithFilters(t *testing.T) {
	cartRepo := newMockCartRepository()
	cartItemRepo := newMockCartItemRepository()
	productRepo := newMockProductRepository()
	usecase := NewCartUsecase(cartRepo, cartItemRepo, productRepo)

	// Add test carts
	testCarts := []*model.Cart{
		{UserID: 1, Total: 50.0},
		{UserID: 2, Total: 75.0},
		{UserID: 3, Total: 100.0},
	}

	for _, c := range testCarts {
		cartRepo.Create(c)
	}

	// Test Case 19: Get all carts with filters
	filters := map[string]string{}
	carts, err := usecase.GetWithFilters(filters)
	// Assertion 67: No error should occur when getting carts with filters
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	// Assertion 68: Should return all carts
	if len(carts) != 3 {
		t.Errorf("Expected 3 carts, got %d", len(carts))
	}
	// Assertion 69: First cart should have correct UserID
	if carts[0].UserID != 1 {
		t.Errorf("Expected first cart UserID 1, got %d", carts[0].UserID)
	}
	// Assertion 70: Second cart should have correct total price
	if carts[1].Total != 75.0 {
		t.Errorf("Expected second cart total 75.0, got %f", carts[1].Total)
	}
}

func TestCartUsecaseAddProduct(t *testing.T) {
	cartRepo := newMockCartRepository()
	cartItemRepo := newMockCartItemRepository()
	productRepo := newMockProductRepository()
	usecase := NewCartUsecase(cartRepo, cartItemRepo, productRepo)

	// Setup test product
	testProduct := &model.Product{
		Name:     "Test Product",
		Price:    25.99,
		Currency: "USD",
		Stock:    10,
		IsActive: true,
	}
	productRepo.Create(testProduct)

	// Test Case 20: Add product with invalid quantity (zero)
	cart, err := usecase.AddProduct(1, 1, 0)
	// Assertion 71: Should return error for zero quantity
	if err == nil {
		t.Error("Expected error for zero quantity")
	}
	// Assertion 72: Cart should be nil for invalid quantity
	if cart != nil {
		t.Error("Expected nil cart for zero quantity")
	}

	// Test Case 21: Add product with negative quantity
	cart, err = usecase.AddProduct(1, 1, -1)
	// Assertion 73: Should return error for negative quantity
	if err == nil {
		t.Error("Expected error for negative quantity")
	}
	// Assertion 74: Cart should be nil for negative quantity
	if cart != nil {
		t.Error("Expected nil cart for negative quantity")
	}

	// Test Case 22: Add non-existent product
	_, err = usecase.AddProduct(1, 999, 1)
	// Assertion 75: Should return error for non-existent product
	if err == nil {
		t.Error("Expected error for non-existent product")
	}

	// Create a cart for the user first (as AddProduct expects cart to exist)
	testCart := &model.Cart{
		UserID: 1,
		Total:  0.0,
	}
	cartRepo.Create(testCart)

	// Test Case 23: Add product to existing cart
	cart, err = usecase.AddProduct(1, 1, 2)
	// Assertion 76: Should not return error when adding product to existing cart
	if err != nil {
		t.Errorf("Expected no error adding product to existing cart, got %v", err)
	}
	// Assertion 77: Cart should not be nil after adding product
	if cart == nil {
		t.Error("Expected cart after adding product, got nil")
	}
}

// Helper function to create test products for cart integration tests
func createCartTestProducts(productRepo *mockProductRepository) {
	products := []*model.Product{
		{Name: "Product 1", Price: 10.0, Currency: "USD", Stock: 100, IsActive: true},
		{Name: "Product 2", Price: 20.0, Currency: "USD", Stock: 50, IsActive: true},
		{Name: "Product 3", Price: 15.0, Currency: "USD", Stock: 25, IsActive: true},
	}

	for _, p := range products {
		productRepo.Create(p)
	}
}

// Helper function to setup cart with items
func setupCartWithItems(cartRepo *mockCartRepository, cartItemRepo *mockCartItemRepository, userID uint) {
	testCart := &model.Cart{
		UserID: userID,
		Total:  0.0,
	}
	cartRepo.Create(testCart)

	cartItem1 := &model.CartItem{
		CartID:    1,
		ProductID: 1,
		Quantity:  2,
		UnitPrice: 10.0,
		Subtotal:  20.0,
	}
	cartItemRepo.AddItem(cartItem1)

	cartItem2 := &model.CartItem{
		CartID:    1,
		ProductID: 2,
		Quantity:  1,
		UnitPrice: 20.0,
		Subtotal:  20.0,
	}
	cartItemRepo.AddItem(cartItem2)
}

func TestCartUsecaseIntegrationCreateAndRetrieve(t *testing.T) {
	cartRepo := newMockCartRepository()
	cartItemRepo := newMockCartItemRepository()
	productRepo := newMockProductRepository()
	usecase := NewCartUsecase(cartRepo, cartItemRepo, productRepo)

	createCartTestProducts(productRepo)
	userID := uint(1)

	// Initially no cart should exist
	_, err := usecase.GetByUserID(userID)
	// Assertion 78: Should return error for non-existent cart initially
	if err != gorm.ErrRecordNotFound {
		t.Errorf("Expected gorm.ErrRecordNotFound initially, got %v", err)
	}

	// Setup cart with items
	setupCartWithItems(cartRepo, cartItemRepo, userID)

	// Get cart and verify items
	cart, err := usecase.GetByUserID(userID)
	// Assertion 79: Should successfully get cart after adding items
	if err != nil {
		t.Errorf("Expected no error getting cart, got %v", err)
	}
	// Assertion 80: Cart should not be nil
	if cart == nil {
		t.Error("Expected cart, got nil")
		return
	}
	// Assertion 81: Cart should belong to correct user
	if cart.UserID != userID {
		t.Errorf("Expected cart UserID %d, got %d", userID, cart.UserID)
	}
}

func TestCartUsecaseIntegrationUpdateItems(t *testing.T) {
	cartRepo := newMockCartRepository()
	cartItemRepo := newMockCartItemRepository()
	productRepo := newMockProductRepository()

	createCartTestProducts(productRepo)
	userID := uint(1)
	setupCartWithItems(cartRepo, cartItemRepo, userID)

	// Test updating cart item quantity
	updatedItem := &model.CartItem{
		ID:        1,
		CartID:    1,
		ProductID: 1,
		Quantity:  3,
		UnitPrice: 10.0,
		Subtotal:  30.0,
	}
	cartItemRepo.UpdateItem(updatedItem)

	// Verify item was updated
	item, err := cartItemRepo.FindByID(1)
	// Assertion 82: Should successfully find updated item
	if err != nil {
		t.Errorf("Expected no error finding updated item, got %v", err)
	}
	// Assertion 83: Updated item quantity should match
	if item.Quantity != 3 {
		t.Errorf("Expected updated quantity 3, got %d", item.Quantity)
	}
	// Assertion 84: Updated item subtotal should match
	if item.Subtotal != 30.0 {
		t.Errorf("Expected updated subtotal 30.0, got %f", item.Subtotal)
	}
}

func TestCartUsecaseIntegrationRemoveItems(t *testing.T) {
	cartRepo := newMockCartRepository()
	cartItemRepo := newMockCartItemRepository()
	productRepo := newMockProductRepository()

	createCartTestProducts(productRepo)
	userID := uint(1)
	setupCartWithItems(cartRepo, cartItemRepo, userID)

	// Test removing item from cart
	err := cartItemRepo.DeleteItem(2)
	// Assertion 85: Should successfully delete cart item
	if err != nil {
		t.Errorf("Expected no error deleting item, got %v", err)
	}

	// Verify item was removed
	_, err = cartItemRepo.FindByID(2)
	// Assertion 86: Should return error for deleted item
	if err != gorm.ErrRecordNotFound {
		t.Errorf("Expected gorm.ErrRecordNotFound for deleted item, got %v", err)
	}
}

func TestCartUsecaseIntegrationClearCart(t *testing.T) {
	cartRepo := newMockCartRepository()
	cartItemRepo := newMockCartItemRepository()
	productRepo := newMockProductRepository()
	usecase := NewCartUsecase(cartRepo, cartItemRepo, productRepo)

	createCartTestProducts(productRepo)
	userID := uint(1)
	setupCartWithItems(cartRepo, cartItemRepo, userID)

	// Test clearing entire cart
	err := cartItemRepo.ClearCart(1)
	// Assertion 87: Should successfully clear cart items
	if err != nil {
		t.Errorf("Expected no error clearing cart, got %v", err)
	}

	// Verify all items were removed
	remainingItem, err := cartItemRepo.FindByID(1)
	// Assertion 88: Should return error for cleared cart items
	if err != gorm.ErrRecordNotFound {
		t.Errorf("Expected gorm.ErrRecordNotFound for cleared items, got %v", err)
	}
	// Assertion 89: No items should remain after clearing cart
	if remainingItem != nil {
		t.Error("Expected nil item after clearing cart")
	}

	// Verify cart still exists but is empty
	cart, err := usecase.GetByUserID(userID)
	// Assertion 90: Cart should still exist after clearing items
	if err != nil {
		t.Errorf("Expected no error getting cart after clearing, got %v", err)
	}
	// Assertion 91: Cart should not be nil after clearing items
	if cart != nil && cart.UserID != userID {
		t.Errorf("Expected cart to still belong to user %d", userID)
	}
}

func TestCartUsecaseIntegrationDeleteCart(t *testing.T) {
	cartRepo := newMockCartRepository()
	cartItemRepo := newMockCartItemRepository()
	productRepo := newMockProductRepository()
	usecase := NewCartUsecase(cartRepo, cartItemRepo, productRepo)

	createCartTestProducts(productRepo)
	userID := uint(1)
	setupCartWithItems(cartRepo, cartItemRepo, userID)

	// Test deleting entire cart
	err := cartRepo.Delete(1)
	// Assertion 92: Should successfully delete cart
	if err != nil {
		t.Errorf("Expected no error deleting cart, got %v", err)
	}

	// Verify cart was deleted
	_, err = usecase.GetByUserID(userID)
	// Assertion 93: Should return error for deleted cart
	if err != gorm.ErrRecordNotFound {
		t.Errorf("Expected gorm.ErrRecordNotFound for deleted cart, got %v", err)
	}
}

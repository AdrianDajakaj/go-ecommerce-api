package usecase

import (
	"testing"

	"go-ecommerce-api/internal/domain/model"

	"gorm.io/gorm"
)

// Test error message constants
const (
	errExpectedGormNotFound = "Expected gorm.ErrRecordNotFound, got %v"
	errExpectedNoError      = "Expected no error, got %v"
)

// Mock ProductRepository for testing
type mockProductRepository struct {
	products []model.Product
	nextID   uint
}

func newMockProductRepository() *mockProductRepository {
	return &mockProductRepository{
		products: []model.Product{},
		nextID:   1,
	}
}

func (m *mockProductRepository) FindByID(id uint) (*model.Product, error) {
	for _, product := range m.products {
		if product.ID == id {
			return &product, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockProductRepository) FindAll() ([]model.Product, error) {
	return m.products, nil
}

func (m *mockProductRepository) FindWithFilters(filters map[string]string) ([]model.Product, error) {
	result := []model.Product{}
	for _, product := range m.products {
		if categoryID, ok := filters["category_id"]; ok {
			if product.CategoryID != 1 && categoryID == "1" {
				continue
			}
		}
		if active, ok := filters["is_active"]; ok {
			if active == "true" && !product.IsActive {
				continue
			}
		}
		result = append(result, product)
	}
	return result, nil
}

func (m *mockProductRepository) Create(product *model.Product) error {
	product.ID = m.nextID
	m.nextID++
	m.products = append(m.products, *product)
	return nil
}

func (m *mockProductRepository) Update(product *model.Product) error {
	for i, p := range m.products {
		if p.ID == product.ID {
			m.products[i] = *product
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *mockProductRepository) Delete(id uint) error {
	for i, p := range m.products {
		if p.ID == id {
			m.products = append(m.products[:i], m.products[i+1:]...)
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func TestProductUsecaseGetByID(t *testing.T) {
	repo := newMockProductRepository()
	usecase := NewProductUsecase(repo)

	// Test Case 1: Get non-existent product
	product, err := usecase.GetByID(999)
	// Assertion 1: Error should be ErrRecordNotFound for non-existent product
	if err != gorm.ErrRecordNotFound {
		t.Errorf(errExpectedGormNotFound, err)
	}
	// Assertion 2: Product should be nil when not found
	if product != nil {
		t.Errorf("Expected nil product, got %v", product)
	}

	// Setup test data
	testProduct := &model.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		Currency:    "USD",
		Stock:       10,
		IsActive:    true,
		CategoryID:  1,
	}
	repo.Create(testProduct)

	// Test Case 2: Get existing product
	product, err = usecase.GetByID(1)
	// Assertion 3: No error should occur when getting existing product
	if err != nil {
		t.Errorf(errExpectedNoError, err)
	}
	// Assertion 4: Product should not be nil when found
	if product == nil {
		t.Error("Expected product, got nil")
		return
	}
	// Assertion 5: Product ID should match requested ID
	if product.ID != 1 {
		t.Errorf("Expected ID 1, got %d", product.ID)
	}
	// Assertion 6: Product name should match test data
	if product.Name != "Test Product" {
		t.Errorf("Expected name 'Test Product', got '%s'", product.Name)
	}
	// Assertion 7: Product price should match test data
	if product.Price != 99.99 {
		t.Errorf("Expected price 99.99, got %f", product.Price)
	}
}

func TestProductUsecaseGetAll(t *testing.T) {
	repo := newMockProductRepository()
	usecase := NewProductUsecase(repo)

	// Test Case 3: Get all products from empty repository
	products, err := usecase.GetAll()
	// Assertion 8: No error should occur when getting all products from empty repo
	if err != nil {
		t.Errorf(errExpectedNoError, err)
	}
	// Assertion 9: Products slice should be empty initially
	if len(products) != 0 {
		t.Errorf("Expected 0 products, got %d", len(products))
	}

	// Add test products
	testProducts := []*model.Product{
		{Name: "Product 1", Price: 10.0, Currency: "USD", Stock: 5, IsActive: true, CategoryID: 1},
		{Name: "Product 2", Price: 20.0, Currency: "USD", Stock: 3, IsActive: false, CategoryID: 2},
		{Name: "Product 3", Price: 30.0, Currency: "USD", Stock: 8, IsActive: true, CategoryID: 1},
	}

	for _, p := range testProducts {
		repo.Create(p)
	}

	// Test Case 4: Get all products with data
	products, err = usecase.GetAll()
	// Assertion 10: No error should occur when getting all products
	if err != nil {
		t.Errorf(errExpectedNoError, err)
	}
	// Assertion 11: Should return correct number of products
	if len(products) != 3 {
		t.Errorf("Expected 3 products, got %d", len(products))
	}
	// Assertion 12: First product name should match
	if products[0].Name != "Product 1" {
		t.Errorf("Expected first product name 'Product 1', got '%s'", products[0].Name)
	}
	// Assertion 13: Second product should be inactive
	if products[1].IsActive {
		t.Error("Expected second product to be inactive")
	}
	// Assertion 14: Third product price should match
	if products[2].Price != 30.0 {
		t.Errorf("Expected third product price 30.0, got %f", products[2].Price)
	}
}

func TestProductUsecaseGetWithFilters(t *testing.T) {
	repo := newMockProductRepository()
	usecase := NewProductUsecase(repo)

	// Add test products
	testProducts := []*model.Product{
		{Name: "Active Product 1", Price: 10.0, Currency: "USD", Stock: 5, IsActive: true, CategoryID: 1},
		{Name: "Inactive Product", Price: 20.0, Currency: "USD", Stock: 3, IsActive: false, CategoryID: 2},
		{Name: "Active Product 2", Price: 30.0, Currency: "USD", Stock: 8, IsActive: true, CategoryID: 1},
	}

	for _, p := range testProducts {
		repo.Create(p)
	}

	// Test Case 5: Filter by category
	filters := map[string]string{"category_id": "1"}
	products, err := usecase.GetWithFilters(filters)
	// Assertion 15: No error should occur when filtering by category
	if err != nil {
		t.Errorf(errExpectedNoError, err)
	}
	// Assertion 16: Should return products from category 1 only
	if len(products) != 2 {
		t.Errorf("Expected 2 products from category 1, got %d", len(products))
	}

	// Test Case 6: Filter by active status
	filters = map[string]string{"is_active": "true"}
	products, err = usecase.GetWithFilters(filters)
	// Assertion 17: No error should occur when filtering by active status
	if err != nil {
		t.Errorf(errExpectedNoError, err)
	}
	// Assertion 18: Should return only active products
	if len(products) != 2 {
		t.Errorf("Expected 2 active products, got %d", len(products))
	}
	// Assertion 19: All returned products should be active
	for _, p := range products {
		if !p.IsActive {
			t.Error("Expected all filtered products to be active")
		}
	}
}

func TestProductUsecaseCreate(t *testing.T) {
	repo := newMockProductRepository()
	usecase := NewProductUsecase(repo)

	// Test Case 7: Create product with nil input
	product, err := usecase.Create(nil)
	// Assertion 20: Should return error for nil product
	if err == nil {
		t.Error("Expected error for nil product")
	}
	// Assertion 21: Should return nil product for invalid input
	if product != nil {
		t.Error("Expected nil product for invalid input")
	}
	// Assertion 22: Error message should indicate invalid product data
	if err.Error() != "invalid product data" {
		t.Errorf("Expected 'invalid product data' error, got '%s'", err.Error())
	}

	// Test Case 8: Create product with empty name
	emptyProduct := &model.Product{Name: "", Price: 10.0}
	product, err = usecase.Create(emptyProduct)
	// Assertion 23: Should return error for empty name
	if err == nil {
		t.Error("Expected error for empty product name")
	}
	// Assertion 24: Should return nil product for empty name
	if product != nil {
		t.Error("Expected nil product for empty name")
	}

	// Test Case 9: Create valid product
	validProduct := &model.Product{
		Name:        "New Product",
		Description: "New Description",
		Price:       49.99,
		Currency:    "USD",
		Stock:       15,
		IsActive:    true,
		CategoryID:  1,
	}
	product, err = usecase.Create(validProduct)
	// Assertion 25: No error should occur for valid product creation
	if err != nil {
		t.Errorf(errExpectedNoError, err)
	}
	// Assertion 26: Created product should not be nil
	if product == nil {
		t.Error("Expected created product, got nil")
		return
	}
	// Assertion 27: Created product should have assigned ID
	if product.ID == 0 {
		t.Error("Expected non-zero ID for created product")
	}
	// Assertion 28: Created product name should match input
	if product.Name != "New Product" {
		t.Errorf("Expected name 'New Product', got '%s'", product.Name)
	}
	// Assertion 29: Created product price should match input
	if product.Price != 49.99 {
		t.Errorf("Expected price 49.99, got %f", product.Price)
	}
}

func TestProductUsecaseUpdate(t *testing.T) {
	repo := newMockProductRepository()
	usecase := NewProductUsecase(repo)

	// Test Case 10: Update with nil product
	product, err := usecase.Update(nil)
	// Assertion 30: Should return error for nil product update
	if err == nil {
		t.Error("Expected error for nil product")
	}
	// Assertion 31: Should return nil product for nil input
	if product != nil {
		t.Error("Expected nil product for nil input")
	}

	// Test Case 11: Update with zero ID
	zeroIDProduct := &model.Product{ID: 0, Name: "Test"}
	_, err = usecase.Update(zeroIDProduct)
	// Assertion 32: Should return error for zero ID
	if err == nil {
		t.Error("Expected error for zero ID")
	}
	// Assertion 33: Error message should indicate invalid product
	if err.Error() != "invalid product" {
		t.Errorf("Expected 'invalid product' error, got '%s'", err.Error())
	}

	// Setup existing product
	existingProduct := &model.Product{
		Name:       "Existing Product",
		Price:      25.0,
		Currency:   "USD",
		Stock:      5,
		IsActive:   true,
		CategoryID: 1,
	}
	repo.Create(existingProduct)

	// Test Case 12: Update non-existent product
	nonExistentProduct := &model.Product{ID: 999, Name: "Non-existent"}
	_, err = usecase.Update(nonExistentProduct)
	// Assertion 34: Should return ErrRecordNotFound for non-existent product
	if err != gorm.ErrRecordNotFound {
		t.Errorf(errExpectedGormNotFound, err)
	}

	// Test Case 13: Update existing product
	updateProduct := &model.Product{
		ID:         1,
		Name:       "Updated Product",
		Price:      35.0,
		Currency:   "USD",
		Stock:      10,
		IsActive:   false,
		CategoryID: 2,
	}
	product, err = usecase.Update(updateProduct)
	// Assertion 35: No error should occur for valid update
	if err != nil {
		t.Errorf(errExpectedNoError, err)
	}
	// Assertion 36: Updated product should not be nil
	if product == nil {
		t.Error("Expected updated product, got nil")
		return
	}
	// Assertion 37: Updated product name should match
	if product.Name != "Updated Product" {
		t.Errorf("Expected name 'Updated Product', got '%s'", product.Name)
	}
	// Assertion 38: Updated product price should match
	if product.Price != 35.0 {
		t.Errorf("Expected price 35.0, got %f", product.Price)
	}
	// Assertion 39: Updated product should be inactive
	if product.IsActive {
		t.Error("Expected updated product to be inactive")
	}
}

func TestProductUsecaseDelete(t *testing.T) {
	repo := newMockProductRepository()
	usecase := NewProductUsecase(repo)

	// Test Case 14: Delete non-existent product
	err := usecase.Delete(999)
	// Assertion 40: Should return ErrRecordNotFound for non-existent product
	if err != gorm.ErrRecordNotFound {
		t.Errorf(errExpectedGormNotFound, err)
	}

	// Setup existing product
	existingProduct := &model.Product{
		Name:       "Product to Delete",
		Price:      15.0,
		Currency:   "USD",
		Stock:      3,
		IsActive:   true,
		CategoryID: 1,
	}
	repo.Create(existingProduct)

	// Verify product exists before deletion
	product, err := usecase.GetByID(1)
	// Assertion 41: Product should exist before deletion
	if err != nil {
		t.Errorf("Expected no error before deletion, got %v", err)
	}
	// Assertion 42: Product should not be nil before deletion
	if product == nil {
		t.Error("Expected product to exist before deletion")
	}

	// Test Case 15: Delete existing product
	err = usecase.Delete(1)
	// Assertion 43: No error should occur for valid deletion
	if err != nil {
		t.Errorf("Expected no error for valid deletion, got %v", err)
	}

	// Verify product is deleted
	product, err = usecase.GetByID(1)
	// Assertion 44: Should return ErrRecordNotFound after deletion
	if err != gorm.ErrRecordNotFound {
		t.Errorf("Expected gorm.ErrRecordNotFound after deletion, got %v", err)
	}
	// Assertion 45: Product should be nil after deletion
	if product != nil {
		t.Error("Expected nil product after deletion")
	}

	// Test repository state after deletion
	allProducts, err := usecase.GetAll()
	// Assertion 46: No error should occur when getting all products after deletion
	if err != nil {
		t.Errorf("Expected no error getting all products, got %v", err)
	}
	// Assertion 47: Repository should be empty after deletion
	if len(allProducts) != 0 {
		t.Errorf("Expected 0 products after deletion, got %d", len(allProducts))
	}
}

// Helper function to create test products
func createTestProducts(usecase ProductUsecase) []*model.Product {
	products := []*model.Product{
		{Name: "Product A", Price: 100.0, Currency: "USD", Stock: 10, IsActive: true, CategoryID: 1},
		{Name: "Product B", Price: 200.0, Currency: "USD", Stock: 5, IsActive: true, CategoryID: 2},
		{Name: "Product C", Price: 150.0, Currency: "USD", Stock: 0, IsActive: false, CategoryID: 1},
	}

	var createdProducts []*model.Product
	for _, p := range products {
		created, _ := usecase.Create(p)
		createdProducts = append(createdProducts, created)
	}
	return createdProducts
}

func TestProductUsecaseIntegrationCreateMultiple(t *testing.T) {
	repo := newMockProductRepository()
	usecase := NewProductUsecase(repo)

	products := createTestProducts(usecase)

	// Assertion 48: Each product creation should succeed
	for i, created := range products {
		if created == nil {
			t.Errorf("Product %d creation failed", i+1)
			continue
		}
		// Assertion 49: Each created product should have valid ID
		if created.ID == 0 {
			t.Error("Expected non-zero ID for created product")
		}
	}

	// Get all products and verify count
	allProducts, err := usecase.GetAll()
	// Assertion 50: No error should occur getting all products
	if err != nil {
		t.Errorf("Expected no error getting all products, got %v", err)
	}
	// Assertion 51: Should have correct number of products
	if len(allProducts) != 3 {
		t.Errorf("Expected 3 products, got %d", len(allProducts))
	}
}

func TestProductUsecaseIntegrationFilterActive(t *testing.T) {
	repo := newMockProductRepository()
	usecase := NewProductUsecase(repo)

	createTestProducts(usecase)

	// Filter active products
	activeFilters := map[string]string{"is_active": "true"}
	activeProducts, err := usecase.GetWithFilters(activeFilters)
	// Assertion 52: No error should occur filtering active products
	if err != nil {
		t.Errorf("Expected no error filtering products, got %v", err)
	}
	// Assertion 53: Should have 2 active products
	if len(activeProducts) != 2 {
		t.Errorf("Expected 2 active products, got %d", len(activeProducts))
	}
}

func TestProductUsecaseIntegrationUpdateProduct(t *testing.T) {
	repo := newMockProductRepository()
	usecase := NewProductUsecase(repo)

	createTestProducts(usecase)

	// Update first product
	updateData := &model.Product{
		ID:         1,
		Name:       "Updated Product A",
		Price:      120.0,
		Currency:   "USD",
		Stock:      15,
		IsActive:   true,
		CategoryID: 1,
	}
	updated, err := usecase.Update(updateData)
	// Assertion 54: No error should occur during update
	if err != nil {
		t.Errorf("Expected no error updating product, got %v", err)
	}
	// Assertion 55: Updated product name should match
	if updated.Name != "Updated Product A" {
		t.Errorf("Expected name 'Updated Product A', got '%s'", updated.Name)
	}
	// Assertion 56: Updated product price should match
	if updated.Price != 120.0 {
		t.Errorf("Expected price 120.0, got %f", updated.Price)
	}
}

func TestProductUsecaseIntegrationDeleteProduct(t *testing.T) {
	repo := newMockProductRepository()
	usecase := NewProductUsecase(repo)

	createTestProducts(usecase)

	// Delete second product
	err := usecase.Delete(2)
	// Assertion 57: No error should occur during deletion
	if err != nil {
		t.Errorf("Expected no error deleting product, got %v", err)
	}

	// Verify final state
	finalProducts, err := usecase.GetAll()
	// Assertion 58: No error should occur getting final products
	if err != nil {
		t.Errorf("Expected no error getting final products, got %v", err)
	}
	// Assertion 59: Should have 2 products remaining
	if len(finalProducts) != 2 {
		t.Errorf("Expected 2 products remaining, got %d", len(finalProducts))
	}

	// Verify deleted product is gone
	_, err = usecase.GetByID(2)
	// Assertion 60: Should return ErrRecordNotFound for deleted product
	if err != gorm.ErrRecordNotFound {
		t.Errorf("Expected gorm.ErrRecordNotFound for deleted product, got %v", err)
	}
}

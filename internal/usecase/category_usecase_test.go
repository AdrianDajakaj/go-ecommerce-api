package usecase

import (
	"errors"
	"testing"

	"go-ecommerce-api/internal/domain/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Test constants for repeated strings
const (
	testCategoryName       = "Test Category"
	newCategoryName        = "New Category"
	updatedCategoryName    = "Updated Category"
	updatedIntegrationName = "Updated Integration Category"
	findFailedError        = "find failed"
	updatedIntegrationURL  = "https://example.com/updated-integration.png"
)

// Mock CategoryRepository for testing
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) FindByID(id uint) (*model.Category, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Category), args.Error(1)
}

func (m *MockCategoryRepository) FindAll() ([]model.Category, error) {
	args := m.Called()
	return args.Get(0).([]model.Category), args.Error(1)
}

func (m *MockCategoryRepository) FindWithFilters(filters map[string]string) ([]model.Category, error) {
	args := m.Called(filters)
	return args.Get(0).([]model.Category), args.Error(1)
}

func (m *MockCategoryRepository) Create(category *model.Category) error {
	// Mock implementation for Create operation - validates and creates new category
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Update(category *model.Category) error {
	// Mock implementation for Update operation - validates and updates existing category
	args := m.Called(category)
	if category != nil && category.ID == 0 {
		return errors.New("invalid ID for update")
	}
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupCategoryUsecase() (*categoryUsecase, *MockCategoryRepository) {
	mockRepo := new(MockCategoryRepository)
	uc := &categoryUsecase{
		categoryRepo: mockRepo,
	}
	return uc, mockRepo
}

func TestNewCategoryUsecase(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	uc := NewCategoryUsecase(mockRepo)

	// Assertion 201: NewCategoryUsecase should return a non-nil usecase instance
	assert.NotNil(t, uc)
	// Assertion 202: NewCategoryUsecase should return a CategoryUsecase interface
	assert.Implements(t, (*CategoryUsecase)(nil), uc)
}

func TestCategoryUsecaseGetByIDSuccess(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	iconURL := "https://example.com/icon.png"
	expectedCategory := &model.Category{
		ID:      1,
		Name:    "Electronics",
		IconURL: &iconURL,
	}

	mockRepo.On("FindByID", uint(1)).Return(expectedCategory, nil)

	result, err := uc.GetByID(1)

	// Assertion 203: GetByID should not return an error for valid category ID
	assert.NoError(t, err)
	// Assertion 204: GetByID should return the expected category
	assert.Equal(t, expectedCategory, result)
	// Assertion 205: GetByID should return a category with correct ID
	assert.Equal(t, uint(1), result.ID)
	// Assertion 206: GetByID should return a category with correct Name
	assert.Equal(t, "Electronics", result.Name)
	// Assertion 207: GetByID should return a category with correct IconURL
	assert.Equal(t, &iconURL, result.IconURL)

	mockRepo.AssertExpectations(t)
}

func TestCategoryUsecaseGetByIDNotFound(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	mockRepo.On("FindByID", uint(999)).Return(nil, nil)

	result, err := uc.GetByID(999)

	// Assertion 208: GetByID should return gorm.ErrRecordNotFound for non-existent category
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	// Assertion 209: GetByID should return nil category for non-existent category
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestCategoryUsecaseGetByIDRepositoryError(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("database error"))

	result, err := uc.GetByID(1)

	// Assertion 210: GetByID should return error when repository fails
	assert.Error(t, err)
	// Assertion 211: GetByID should return nil category when repository fails
	assert.Nil(t, result)
	// Assertion 212: GetByID should return the exact repository error
	assert.EqualError(t, err, "database error")

	mockRepo.AssertExpectations(t)
}

func TestCategoryUsecaseGetAllSuccess(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	iconURL1 := "https://example.com/icon1.png"
	iconURL2 := "https://example.com/icon2.png"
	expectedCategories := []model.Category{
		{ID: 1, Name: "Electronics", IconURL: &iconURL1},
		{ID: 2, Name: "Clothing", IconURL: &iconURL2},
		{ID: 3, Name: "Books", IconURL: nil},
	}

	mockRepo.On("FindAll").Return(expectedCategories, nil)

	result, err := uc.GetAll()

	// Assertion 213: GetAll should not return an error when repository succeeds
	assert.NoError(t, err)
	// Assertion 214: GetAll should return the expected categories slice
	assert.Equal(t, expectedCategories, result)
	// Assertion 215: GetAll should return correct number of categories
	assert.Len(t, result, 3)
	// Assertion 216: GetAll should return categories with correct names
	assert.Equal(t, "Electronics", result[0].Name)
	// Assertion 217: GetAll should return categories with correct names for second category
	assert.Equal(t, "Clothing", result[1].Name)
	// Assertion 218: GetAll should handle categories with nil IconURL
	assert.Nil(t, result[2].IconURL)

	mockRepo.AssertExpectations(t)
}

func TestCategoryUsecaseGetAllEmptyResult(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	mockRepo.On("FindAll").Return([]model.Category{}, nil)

	result, err := uc.GetAll()

	// Assertion 219: GetAll should not return an error when no categories exist
	assert.NoError(t, err)
	// Assertion 220: GetAll should return empty slice when no categories exist
	assert.Empty(t, result)
	// Assertion 221: GetAll should return slice with zero length when no categories exist
	assert.Len(t, result, 0)

	mockRepo.AssertExpectations(t)
}

func TestCategoryUsecaseGetAllRepositoryError(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	mockRepo.On("FindAll").Return([]model.Category{}, errors.New("database connection failed"))

	result, err := uc.GetAll()

	// Assertion 222: GetAll should return error when repository fails
	assert.Error(t, err)
	// Assertion 223: GetAll should return empty slice when repository fails
	assert.Empty(t, result)
	// Assertion 224: GetAll should return the exact repository error
	assert.EqualError(t, err, "database connection failed")

	mockRepo.AssertExpectations(t)
}

func TestCategoryUsecaseGetWithFiltersSuccess(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	filters := map[string]string{
		"parent_id": "1",
		"name":      "Sub",
	}

	iconURL := "https://example.com/subcategory.png"
	parentID := uint(1)
	expectedCategories := []model.Category{
		{ID: 2, Name: "Subcategory 1", IconURL: &iconURL, ParentID: &parentID},
		{ID: 3, Name: "Subcategory 2", IconURL: nil, ParentID: &parentID},
	}

	mockRepo.On("FindWithFilters", filters).Return(expectedCategories, nil)

	result, err := uc.GetWithFilters(filters)

	// Assertion 225: GetWithFilters should not return an error when repository succeeds
	assert.NoError(t, err)
	// Assertion 226: GetWithFilters should return the expected filtered categories
	assert.Equal(t, expectedCategories, result)
	// Assertion 227: GetWithFilters should return correct number of filtered categories
	assert.Len(t, result, 2)
	// Assertion 228: GetWithFilters should return categories with correct parent ID
	assert.Equal(t, &parentID, result[0].ParentID)
	// Assertion 229: GetWithFilters should return categories with correct parent ID for second category
	assert.Equal(t, &parentID, result[1].ParentID)

	mockRepo.AssertExpectations(t)
}

func TestCategoryUsecaseGetWithFiltersEmptyResult(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	filters := map[string]string{"name": "NonExistent"}

	mockRepo.On("FindWithFilters", filters).Return([]model.Category{}, nil)

	result, err := uc.GetWithFilters(filters)

	// Assertion 230: GetWithFilters should not return an error when no matches found
	assert.NoError(t, err)
	// Assertion 231: GetWithFilters should return empty slice when no matches found
	assert.Empty(t, result)
	// Assertion 232: GetWithFilters should return slice with zero length when no matches found
	assert.Len(t, result, 0)

	mockRepo.AssertExpectations(t)
}

func TestCategoryUsecaseGetWithFiltersRepositoryError(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	filters := map[string]string{"invalid": "filter"}

	mockRepo.On("FindWithFilters", filters).Return([]model.Category{}, errors.New("invalid filter"))

	result, err := uc.GetWithFilters(filters)

	// Assertion 233: GetWithFilters should return error when repository fails
	assert.Error(t, err)
	// Assertion 234: GetWithFilters should return empty slice when repository fails
	assert.Empty(t, result)
	// Assertion 235: GetWithFilters should return the exact repository error
	assert.EqualError(t, err, "invalid filter")

	mockRepo.AssertExpectations(t)
}

func TestCategoryUsecaseCreateSuccess(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	iconURL := "https://example.com/new-category.png"
	newCategory := &model.Category{
		Name:    newCategoryName,
		IconURL: &iconURL,
	}

	createdCategory := &model.Category{
		ID:      1,
		Name:    newCategoryName,
		IconURL: &iconURL,
	}

	mockRepo.On("Create", newCategory).Return(nil).Run(func(args mock.Arguments) {
		cat := args.Get(0).(*model.Category)
		cat.ID = 1
	})
	mockRepo.On("FindByID", uint(1)).Return(createdCategory, nil)

	result, err := uc.Create(newCategory)

	// Assertion 236: Create should not return an error for valid category
	assert.NoError(t, err)
	// Assertion 237: Create should return the created category
	assert.Equal(t, createdCategory, result)
	// Assertion 238: Create should return a category with assigned ID
	assert.Equal(t, uint(1), result.ID)
	// Assertion 239: Create should return a category with correct name
	assert.Equal(t, newCategoryName, result.Name)
	// Assertion 240: Create should return a category with correct IconURL
	assert.Equal(t, &iconURL, result.IconURL)

	mockRepo.AssertExpectations(t)
}

func TestCategoryUsecaseCreateNilCategory(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	result, err := uc.Create(nil)

	// Assertion 241: Create should return error for nil category
	assert.Error(t, err)
	// Assertion 242: Create should return nil category for nil input
	assert.Nil(t, result)
	// Assertion 243: Create should return appropriate error message for nil category
	assert.EqualError(t, err, "invalid category data")

	mockRepo.AssertNotCalled(t, "Create")
}

func TestCategoryUsecaseCreateEmptyName(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	emptyCategory := &model.Category{Name: ""}

	result, err := uc.Create(emptyCategory)

	// Assertion 244: Create should return error for empty name
	assert.Error(t, err)
	// Assertion 245: Create should return nil category for empty name
	assert.Nil(t, result)
	// Assertion 246: Create should return appropriate error message for empty name
	assert.EqualError(t, err, "invalid category data")

	mockRepo.AssertNotCalled(t, "Create")
}

func TestCategoryUsecaseCreateRepositoryCreateError(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	newCategory := &model.Category{Name: testCategoryName}

	mockRepo.On("Create", newCategory).Return(errors.New("create failed"))

	result, err := uc.Create(newCategory)

	// Assertion 247: Create should return error when repository create fails
	assert.Error(t, err)
	// Assertion 248: Create should return nil category when repository create fails
	assert.Nil(t, result)
	// Assertion 249: Create should return the exact repository error
	assert.EqualError(t, err, "create failed")

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "FindByID")
}

func TestCategoryUsecaseCreateRepositoryFindError(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	newCategory := &model.Category{Name: testCategoryName}

	mockRepo.On("Create", newCategory).Return(nil).Run(func(args mock.Arguments) {
		cat := args.Get(0).(*model.Category)
		cat.ID = 1
	})
	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New(findFailedError))

	result, err := uc.Create(newCategory)

	// Assertion 250: Create should return error when repository find fails after create
	assert.Error(t, err)
	// Assertion 251: Create should return nil category when repository find fails after create
	assert.Nil(t, result)
	// Assertion 252: Create should return the exact find repository error
	assert.EqualError(t, err, findFailedError)

	mockRepo.AssertExpectations(t)
}

func TestCategoryUsecaseUpdateSuccess(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	iconURL := "https://example.com/updated-category.png"
	updateCategory := &model.Category{
		ID:      1,
		Name:    updatedCategoryName,
		IconURL: &iconURL,
	}

	updatedCategory := &model.Category{
		ID:      1,
		Name:    updatedCategoryName,
		IconURL: &iconURL,
	}

	mockRepo.On("Update", updateCategory).Return(nil)
	mockRepo.On("FindByID", uint(1)).Return(updatedCategory, nil)

	result, err := uc.Update(updateCategory)

	// Assertion 253: Update should not return an error for valid category
	assert.NoError(t, err)
	// Assertion 254: Update should return the updated category
	assert.Equal(t, updatedCategory, result)
	// Assertion 255: Update should return a category with correct ID
	assert.Equal(t, uint(1), result.ID)
	// Assertion 256: Update should return a category with updated name
	assert.Equal(t, updatedCategoryName, result.Name)
	// Assertion 257: Update should return a category with updated IconURL
	assert.Equal(t, &iconURL, result.IconURL)

	mockRepo.AssertExpectations(t)
}

func TestCategoryUsecaseUpdateNilCategory(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	result, err := uc.Update(nil)

	// Assertion 258: Update should return error for nil category
	assert.Error(t, err)
	// Assertion 259: Update should return nil category for nil input
	assert.Nil(t, result)
	// Assertion 260: Update should return appropriate error message for nil category
	assert.EqualError(t, err, "invalid category")

	mockRepo.AssertNotCalled(t, "Update")
}

func TestCategoryUsecaseUpdateZeroID(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	zeroIDCategory := &model.Category{ID: 0, Name: "Test"}

	result, err := uc.Update(zeroIDCategory)

	// Assertion 261: Update should return error for zero ID
	assert.Error(t, err)
	// Assertion 262: Update should return nil category for zero ID
	assert.Nil(t, result)
	// Assertion 263: Update should return appropriate error message for zero ID
	assert.EqualError(t, err, "invalid category")

	mockRepo.AssertNotCalled(t, "Update")
}

func TestCategoryUsecaseUpdateRepositoryUpdateError(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	updateCategory := &model.Category{ID: 1, Name: testCategoryName}

	mockRepo.On("Update", updateCategory).Return(errors.New("update failed"))

	result, err := uc.Update(updateCategory)

	// Assertion 264: Update should return error when repository update fails
	assert.Error(t, err)
	// Assertion 265: Update should return nil category when repository update fails
	assert.Nil(t, result)
	// Assertion 266: Update should return the exact repository error
	assert.EqualError(t, err, "update failed")

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "FindByID")
}

func TestCategoryUsecaseUpdateRepositoryFindError(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	updateCategory := &model.Category{ID: 1, Name: testCategoryName}

	mockRepo.On("Update", updateCategory).Return(nil)
	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New(findFailedError))

	result, err := uc.Update(updateCategory)

	// Assertion 267: Update should return error when repository find fails after update
	assert.Error(t, err)
	// Assertion 268: Update should return nil category when repository find fails after update
	assert.Nil(t, result)
	// Assertion 269: Update should return the exact find repository error
	assert.EqualError(t, err, findFailedError)

	mockRepo.AssertExpectations(t)
}

func TestCategoryUsecaseDeleteSuccess(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	existingCategory := &model.Category{
		ID:   1,
		Name: "Category to Delete",
	}

	mockRepo.On("FindByID", uint(1)).Return(existingCategory, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	err := uc.Delete(1)

	// Assertion 270: Delete should not return an error for valid category ID
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestCategoryUsecaseDeleteCategoryNotFound(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	mockRepo.On("FindByID", uint(999)).Return(nil, nil)

	err := uc.Delete(999)

	// Assertion 271: Delete should return gorm.ErrRecordNotFound for non-existent category
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Delete")
}

func TestCategoryUsecaseDeleteRepositoryFindError(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("find error"))

	err := uc.Delete(1)

	// Assertion 272: Delete should return error when repository find fails
	assert.Error(t, err)
	// Assertion 273: Delete should return the exact find repository error
	assert.EqualError(t, err, "find error")

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Delete")
}

func TestCategoryUsecaseDeleteRepositoryDeleteError(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	existingCategory := &model.Category{ID: 1, Name: testCategoryName}

	mockRepo.On("FindByID", uint(1)).Return(existingCategory, nil)
	mockRepo.On("Delete", uint(1)).Return(errors.New("delete failed"))

	err := uc.Delete(1)

	// Assertion 274: Delete should return error when repository delete fails
	assert.Error(t, err)
	// Assertion 275: Delete should return the exact delete repository error
	assert.EqualError(t, err, "delete failed")

	mockRepo.AssertExpectations(t)
}

// Helper function to create test category for integration tests
func createIntegrationTestCategory() (*model.Category, *model.Category) {
	iconURL := "https://example.com/integration.png"
	parentID := uint(1)

	newCategory := &model.Category{
		Name:     "Integration Category",
		IconURL:  &iconURL,
		ParentID: &parentID,
	}

	createdCategory := &model.Category{
		ID:       2,
		Name:     "Integration Category",
		IconURL:  &iconURL,
		ParentID: &parentID,
	}

	return newCategory, createdCategory
}

func TestCategoryUsecaseIntegrationCreate(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	newCategory, createdCategory := createIntegrationTestCategory()

	// Mock create flow
	mockRepo.On("Create", newCategory).Return(nil).Run(func(args mock.Arguments) {
		cat := args.Get(0).(*model.Category)
		cat.ID = 2
	})
	mockRepo.On("FindByID", uint(2)).Return(createdCategory, nil).Once()

	// Create the category
	created, err := uc.Create(newCategory)

	// Assertion 276: Integration test should successfully create category
	assert.NoError(t, err)
	// Assertion 277: Integration test should return created category with correct properties
	assert.Equal(t, createdCategory, created)

	mockRepo.AssertExpectations(t)
}

func TestCategoryUsecaseIntegrationUpdate(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	parentID := uint(1)
	updatedIconURL := updatedIntegrationURL
	updateCategory := &model.Category{
		ID:       2,
		Name:     updatedIntegrationName,
		IconURL:  &updatedIconURL,
		ParentID: &parentID,
	}

	updatedCategory := &model.Category{
		ID:       2,
		Name:     updatedIntegrationName,
		IconURL:  &updatedIconURL,
		ParentID: &parentID,
	}

	mockRepo.On("Update", updateCategory).Return(nil)
	mockRepo.On("FindByID", uint(2)).Return(updatedCategory, nil).Once()

	// Update the category
	updated, err := uc.Update(updateCategory)

	// Assertion 278: Integration test should successfully update category
	assert.NoError(t, err)
	// Assertion 279: Integration test should return updated category with new properties
	assert.Equal(t, updatedCategory, updated)
	// Assertion 280: Integration test should show updated name
	assert.Equal(t, updatedIntegrationName, updated.Name)

	mockRepo.AssertExpectations(t)
}

func TestCategoryUsecaseIntegrationRetrieve(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	parentID := uint(1)
	updatedIconURL := updatedIntegrationURL
	updatedCategory := &model.Category{
		ID:       2,
		Name:     updatedIntegrationName,
		IconURL:  &updatedIconURL,
		ParentID: &parentID,
	}

	// Get the category by ID
	mockRepo.On("FindByID", uint(2)).Return(updatedCategory, nil).Once()

	retrieved, err := uc.GetByID(2)

	// Assertion 281: Integration test should successfully retrieve category by ID
	assert.NoError(t, err)
	// Assertion 282: Integration test should return the same category as updated
	assert.Equal(t, updatedCategory, retrieved)

	mockRepo.AssertExpectations(t)
}

func TestCategoryUsecaseIntegrationListAndFilter(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	parentID := uint(1)
	updatedIconURL := updatedIntegrationURL
	updatedCategory := &model.Category{
		ID:       2,
		Name:     updatedIntegrationName,
		IconURL:  &updatedIconURL,
		ParentID: &parentID,
	}

	// List all categories
	allCategories := []model.Category{*updatedCategory}
	mockRepo.On("FindAll").Return(allCategories, nil)

	all, err := uc.GetAll()

	// Assertion 283: Integration test should successfully get all categories
	assert.NoError(t, err)
	// Assertion 284: Integration test should return correct number of categories
	assert.Len(t, all, 1)
	// Assertion 285: Integration test should return the updated category in the list
	assert.Equal(t, *updatedCategory, all[0])

	// Filter categories
	filters := map[string]string{"parent_id": "1"}
	filteredCategories := []model.Category{*updatedCategory}
	mockRepo.On("FindWithFilters", filters).Return(filteredCategories, nil)

	filtered, err := uc.GetWithFilters(filters)

	// Assertion 286: Integration test should successfully filter categories
	assert.NoError(t, err)
	// Assertion 287: Integration test should return correct filtered results
	assert.Len(t, filtered, 1)
	// Assertion 288: Integration test should return categories matching filter criteria
	assert.Equal(t, &parentID, filtered[0].ParentID)

	mockRepo.AssertExpectations(t)
}

func TestCategoryUsecaseIntegrationDelete(t *testing.T) {
	uc, mockRepo := setupCategoryUsecase()

	parentID := uint(1)
	updatedIconURL := updatedIntegrationURL
	updatedCategory := &model.Category{
		ID:       2,
		Name:     updatedIntegrationName,
		IconURL:  &updatedIconURL,
		ParentID: &parentID,
	}

	// Delete the category
	mockRepo.On("FindByID", uint(2)).Return(updatedCategory, nil).Once()
	mockRepo.On("Delete", uint(2)).Return(nil)

	err := uc.Delete(2)

	// Assertion 289: Integration test should successfully delete category
	assert.NoError(t, err)

	// Verify category is deleted
	mockRepo.On("FindByID", uint(2)).Return(nil, nil)

	deleted, err := uc.GetByID(2)

	// Assertion 290: Integration test should return error for deleted category
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	// Assertion 291: Integration test should return nil for deleted category
	assert.Nil(t, deleted)

	mockRepo.AssertExpectations(t)
}

package usecase

import (
	"errors"
	"testing"

	"go-ecommerce-api/internal/domain/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Test constants - User specific
const (
	addressCreationFailed   = "address creation failed"
	findFailed              = "find failed"
	deleteFailed            = "delete failed"
	findError               = "find error"
	invalidInput            = "invalid input"
	invalidUser             = "invalid user"
	invalidCredentials      = "invalid credentials"
	emailAlreadyInUse       = "email already in use"
	johnDoeEmail            = "john.doe@example.com"
	janeSmithEmail          = "jane.smith@example.com"
	bobWilsonEmail          = "bob.wilson@example.com"
	newUserEmail            = "new.user@example.com"
	testEmail               = "test@example.com"
	existingEmail           = "existing@example.com"
	userExampleEmail        = "user@example.com"
	updatedEmail            = "updated@example.com"
	integrationEmail        = "integration@example.com"
	updatedIntegrationEmail = "updated.integration@example.com"
	hashEmail               = "hash@example.com"
	nonexistentEmail        = "nonexistent@example.com"
	mainStreet              = "Main St"
	integrationStreet       = "Integration St"
	modelUser               = "*model.User"
	userTestCity            = "Anytown"
	integrationNumber       = "456"
	integrationCity         = "Test City"
	hashNumber              = "789"
	hashCity                = "Hash City"
	johnName                = "John"
	doeSurname              = "Doe"
	janeName                = "Jane"
	smithSurname            = "Smith"
	bobName                 = "Bob"
	wilsonSurname           = "Wilson"
	newName                 = "New"
	userSurname             = "User"
	testName                = "Test"
	existingName            = "Existing"
	updatedName             = "Updated"
	integrationName         = "Integration"
	testUserSurname         = "Test User"
	hashName                = "Hash"
	testSurname             = "Test"
	userRole                = "user"
	adminRole               = "admin"
	password123             = "password123"
	correctPassword         = "correctpassword"
	wrongPassword           = "wrongpassword"
	plainPassword           = "plainpassword"
)

func setupUserUsecase() (*userUsecase, *MockUserRepository, *MockAddressRepository) {
	mockUserRepo := new(MockUserRepository)
	mockAddrRepo := new(MockAddressRepository)

	uc := &userUsecase{
		userRepo: mockUserRepo,
		addrRepo: mockAddrRepo,
	}

	return uc, mockUserRepo, mockAddrRepo
}

func TestNewUserUsecase(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAddrRepo := new(MockAddressRepository)

	uc := NewUserUsecase(mockUserRepo, mockAddrRepo)

	// Assertion 292: NewUserUsecase should return a non-nil usecase instance
	assert.NotNil(t, uc)
	// Assertion 293: NewUserUsecase should return a UserUsecase interface
	assert.Implements(t, (*UserUsecase)(nil), uc)
}

func TestUserUsecaseGetByIDSuccess(t *testing.T) {
	uc, mockUserRepo, _ := setupUserUsecase()

	expectedUser := &model.User{
		ID:        1,
		Email:     johnDoeEmail,
		Name:      newName,
		Surname:   doeSurname,
		Role:      userRole,
		AddressID: 1,
	}

	mockUserRepo.On("FindByID", uint(1)).Return(expectedUser, nil)

	result, err := uc.GetByID(1)

	// Assertion 294: GetByID should not return an error for valid user ID
	assert.NoError(t, err)
	// Assertion 295: GetByID should return the expected user
	assert.Equal(t, expectedUser, result)
	// Assertion 296: GetByID should return a user with correct ID
	assert.Equal(t, uint(1), result.ID)
	// Assertion 297: GetByID should return a user with correct email
	assert.Equal(t, johnDoeEmail, result.Email)
	// Assertion 298: GetByID should return a user with correct name
	assert.Equal(t, newName, result.Name)
	// Assertion 299: GetByID should return a user with correct surname
	assert.Equal(t, doeSurname, result.Surname)
	// Assertion 300: GetByID should return a user with correct role
	assert.Equal(t, userRole, result.Role)

	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecaseGetByIDNotFound(t *testing.T) {
	uc, mockUserRepo, _ := setupUserUsecase()

	mockUserRepo.On("FindByID", uint(999)).Return(nil, nil)

	result, err := uc.GetByID(999)

	// Assertion 301: GetByID should return gorm.ErrRecordNotFound for non-existent user
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	// Assertion 302: GetByID should return nil user for non-existent user
	assert.Nil(t, result)

	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecaseGetByIDRepositoryError(t *testing.T) {
	uc, mockUserRepo, _ := setupUserUsecase()

	mockUserRepo.On("FindByID", uint(1)).Return(nil, errors.New(dbError))

	result, err := uc.GetByID(1)

	// Assertion 303: GetByID should return error when repository fails
	assert.Error(t, err)
	// Assertion 304: GetByID should return nil user when repository fails
	assert.Nil(t, result)
	// Assertion 305: GetByID should return the exact repository error
	assert.EqualError(t, err, dbError)

	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecaseGetAllSuccess(t *testing.T) {
	uc, mockUserRepo, _ := setupUserUsecase()

	expectedUsers := []model.User{
		{ID: 1, Email: johnDoeEmail, Name: newName, Surname: doeSurname, Role: userRole},
		{ID: 2, Email: janeSmithEmail, Name: janeName, Surname: smithSurname, Role: adminRole},
		{ID: 3, Email: bobWilsonEmail, Name: bobName, Surname: wilsonSurname, Role: userRole},
	}

	mockUserRepo.On("FindAll").Return(expectedUsers, nil)

	result, err := uc.GetAll()

	// Assertion 306: GetAll should not return an error when repository succeeds
	assert.NoError(t, err)
	// Assertion 307: GetAll should return the expected users slice
	assert.Equal(t, expectedUsers, result)
	// Assertion 308: GetAll should return correct number of users
	assert.Len(t, result, 3)
	// Assertion 309: GetAll should return users with correct emails
	assert.Equal(t, johnDoeEmail, result[0].Email)
	// Assertion 310: GetAll should return users with different roles
	assert.Equal(t, adminRole, result[1].Role)

	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecaseGetAllEmptyResult(t *testing.T) {
	uc, mockUserRepo, _ := setupUserUsecase()

	mockUserRepo.On("FindAll").Return([]model.User{}, nil)

	result, err := uc.GetAll()

	// Assertion 311: GetAll should not return an error when no users exist
	assert.NoError(t, err)
	// Assertion 312: GetAll should return empty slice when no users exist
	assert.Empty(t, result)
	// Assertion 313: GetAll should return slice with zero length when no users exist
	assert.Len(t, result, 0)

	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecaseGetAllRepositoryError(t *testing.T) {
	uc, mockUserRepo, _ := setupUserUsecase()

	mockUserRepo.On("FindAll").Return([]model.User{}, errors.New(dbConnectionFailed))

	result, err := uc.GetAll()

	// Assertion 314: GetAll should return error when repository fails
	assert.Error(t, err)
	// Assertion 315: GetAll should return empty slice when repository fails
	assert.Empty(t, result)
	// Assertion 316: GetAll should return the exact repository error
	assert.EqualError(t, err, dbConnectionFailed)

	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecaseGetWithFiltersSuccess(t *testing.T) {
	uc, mockUserRepo, _ := setupUserUsecase()

	filters := map[string]string{
		"role": "admin",
		"name": "Jane",
	}

	expectedUsers := []model.User{
		{ID: 2, Email: janeSmithEmail, Name: janeName, Surname: smithSurname, Role: adminRole},
	}

	mockUserRepo.On("FindWithFilters", filters).Return(expectedUsers, nil)

	result, err := uc.GetWithFilters(filters)

	// Assertion 317: GetWithFilters should not return an error when repository succeeds
	assert.NoError(t, err)
	// Assertion 318: GetWithFilters should return the expected filtered users
	assert.Equal(t, expectedUsers, result)
	// Assertion 319: GetWithFilters should return correct number of filtered users
	assert.Len(t, result, 1)
	// Assertion 320: GetWithFilters should return users matching the filter role
	assert.Equal(t, adminRole, result[0].Role)
	// Assertion 321: GetWithFilters should return users matching the filter name
	assert.Equal(t, janeName, result[0].Name)

	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecaseGetWithFiltersRepositoryError(t *testing.T) {
	uc, mockUserRepo, _ := setupUserUsecase()

	filters := map[string]string{"invalid": "filter"}

	mockUserRepo.On("FindWithFilters", filters).Return([]model.User{}, errors.New(invalidFilter))

	result, err := uc.GetWithFilters(filters)

	// Assertion 322: GetWithFilters should return error when repository fails
	assert.Error(t, err)
	// Assertion 323: GetWithFilters should return empty slice when repository fails
	assert.Empty(t, result)
	// Assertion 324: GetWithFilters should return the exact repository error
	assert.EqualError(t, err, invalidFilter)

	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecaseRegisterSuccess(t *testing.T) {
	uc, mockUserRepo, mockAddrRepo := setupUserUsecase()

	address := &model.Address{
		Street: mainStreet,
		Number: testNumber,
		City:   userTestCity,
	}

	user := &model.User{
		Email:   newUserEmail,
		Name:    newName,
		Surname: userSurname,
		Role:    userRole,
	}

	createdUser := &model.User{
		ID:        1,
		Email:     newUserEmail,
		Name:      newName,
		Surname:   userSurname,
		Role:      userRole,
		AddressID: 1,
	}

	mockUserRepo.On("FindByEmail", newUserEmail).Return(nil, nil)
	mockAddrRepo.On("Create", address).Return(nil).Run(func(args mock.Arguments) {
		addr := args.Get(0).(*model.Address)
		addr.ID = 1
	})
	mockUserRepo.On("Create", mock.AnythingOfType(modelUser)).Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(0).(*model.User)
		u.ID = 1
	})
	mockUserRepo.On("FindByID", uint(1)).Return(createdUser, nil)

	result, err := uc.Register(user, password123, address)

	// Assertion 325: Register should not return an error for valid registration
	assert.NoError(t, err)
	// Assertion 326: Register should return the created user
	assert.Equal(t, createdUser, result)
	// Assertion 327: Register should return a user with assigned ID
	assert.Equal(t, uint(1), result.ID)
	// Assertion 328: Register should return a user with correct email
	assert.Equal(t, newUserEmail, result.Email)
	// Assertion 329: Register should return a user with assigned address ID
	assert.Equal(t, uint(1), result.AddressID)

	mockUserRepo.AssertExpectations(t)
	mockAddrRepo.AssertExpectations(t)
}

func TestUserUsecaseRegisterNilUser(t *testing.T) {
	uc, mockUserRepo, mockAddrRepo := setupUserUsecase()

	address := &model.Address{Street: mainStreet, Number: testNumber, City: userTestCity}

	result, err := uc.Register(nil, password123, address)

	// Assertion 330: Register should return error for nil user
	assert.Error(t, err)
	// Assertion 331: Register should return nil user for nil input
	assert.Nil(t, result)
	// Assertion 332: Register should return appropriate error message for nil user
	assert.EqualError(t, err, invalidInput)

	mockUserRepo.AssertNotCalled(t, "FindByEmail")
	mockAddrRepo.AssertNotCalled(t, "Create")
}

func TestUserUsecaseRegisterNilAddress(t *testing.T) {
	uc, mockUserRepo, mockAddrRepo := setupUserUsecase()

	user := &model.User{Email: testEmail, Name: newName, Surname: userSurname}

	result, err := uc.Register(user, password123, nil)

	// Assertion 333: Register should return error for nil address
	assert.Error(t, err)
	// Assertion 334: Register should return nil user for nil address
	assert.Nil(t, result)
	// Assertion 335: Register should return appropriate error message for nil address
	assert.EqualError(t, err, invalidInput)

	mockUserRepo.AssertNotCalled(t, "FindByEmail")
	mockAddrRepo.AssertNotCalled(t, "Create")
}

func TestUserUsecaseRegisterEmailAlreadyExists(t *testing.T) {
	uc, mockUserRepo, mockAddrRepo := setupUserUsecase()

	existingUser := &model.User{
		ID:    1,
		Email: existingEmail,
		Name:  existingName,
	}

	user := &model.User{Email: existingEmail, Name: newName, Surname: userSurname}
	address := &model.Address{Street: mainStreet, Number: testNumber, City: userTestCity}

	mockUserRepo.On("FindByEmail", existingEmail).Return(existingUser, nil)

	result, err := uc.Register(user, password123, address)

	// Assertion 336: Register should return error for existing email
	assert.Error(t, err)
	// Assertion 337: Register should return nil user for existing email
	assert.Nil(t, result)
	// Assertion 338: Register should return appropriate error message for existing email
	assert.EqualError(t, err, emailAlreadyInUse)

	mockUserRepo.AssertExpectations(t)
	mockAddrRepo.AssertNotCalled(t, "Create")
}

func TestUserUsecaseRegisterAddressCreationError(t *testing.T) {
	uc, mockUserRepo, mockAddrRepo := setupUserUsecase()

	user := &model.User{Email: testEmail, Name: newName, Surname: userSurname}
	address := &model.Address{Street: mainStreet, Number: testNumber, City: userTestCity}

	mockUserRepo.On("FindByEmail", testEmail).Return(nil, nil)
	mockAddrRepo.On("Create", address).Return(errors.New(addressCreationFailed))

	result, err := uc.Register(user, password123, address)

	// Assertion 339: Register should return error when address creation fails
	assert.Error(t, err)
	// Assertion 340: Register should return nil user when address creation fails
	assert.Nil(t, result)
	// Assertion 341: Register should return the exact address creation error
	assert.EqualError(t, err, addressCreationFailed)

	mockUserRepo.AssertExpectations(t)
	mockAddrRepo.AssertExpectations(t)
}

func TestUserUsecaseLoginSuccess(t *testing.T) {
	uc, mockUserRepo, _ := setupUserUsecase()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password123), bcrypt.DefaultCost)
	existingUser := &model.User{
		ID:       1,
		Email:    userExampleEmail,
		Password: string(hashedPassword),
		Name:     testName,
		Surname:  userSurname,
	}

	mockUserRepo.On("FindByEmail", userExampleEmail).Return(existingUser, nil)

	result, err := uc.Login(userExampleEmail, password123)

	// Assertion 342: Login should not return an error for valid credentials
	assert.NoError(t, err)
	// Assertion 343: Login should return the authenticated user
	assert.Equal(t, existingUser, result)
	// Assertion 344: Login should return a user with correct email
	assert.Equal(t, userExampleEmail, result.Email)
	// Assertion 345: Login should return a user with correct ID
	assert.Equal(t, uint(1), result.ID)

	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecaseLoginUserNotFound(t *testing.T) {
	uc, mockUserRepo, _ := setupUserUsecase()

	mockUserRepo.On("FindByEmail", nonexistentEmail).Return(nil, nil)

	result, err := uc.Login(nonexistentEmail, password123)

	// Assertion 346: Login should return error for non-existent user
	assert.Error(t, err)
	// Assertion 347: Login should return nil user for non-existent user
	assert.Nil(t, result)
	// Assertion 348: Login should return appropriate error message for non-existent user
	assert.EqualError(t, err, invalidCredentials)

	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecaseLoginWrongPassword(t *testing.T) {
	uc, mockUserRepo, _ := setupUserUsecase()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)
	existingUser := &model.User{
		ID:       1,
		Email:    userExampleEmail,
		Password: string(hashedPassword),
		Name:     testName,
		Surname:  userSurname,
	}

	mockUserRepo.On("FindByEmail", userExampleEmail).Return(existingUser, nil)

	result, err := uc.Login(userExampleEmail, wrongPassword)

	// Assertion 349: Login should return error for wrong password
	assert.Error(t, err)
	// Assertion 350: Login should return nil user for wrong password
	assert.Nil(t, result)
	// Assertion 351: Login should return appropriate error message for wrong password
	assert.EqualError(t, err, invalidCredentials)

	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecaseLoginRepositoryError(t *testing.T) {
	uc, mockUserRepo, _ := setupUserUsecase()

	mockUserRepo.On("FindByEmail", userExampleEmail).Return(nil, errors.New(dbError))

	result, err := uc.Login(userExampleEmail, password123)

	// Assertion 352: Login should return error when repository fails
	assert.Error(t, err)
	// Assertion 353: Login should return nil user when repository fails
	assert.Nil(t, result)
	// Assertion 354: Login should return the exact repository error
	assert.EqualError(t, err, dbError)

	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecaseUpdateSuccess(t *testing.T) {
	uc, mockUserRepo, _ := setupUserUsecase()

	updateUser := &model.User{
		ID:      1,
		Email:   updatedEmail,
		Name:    updatedName,
		Surname: userSurname,
		Role:    adminRole,
	}

	updatedUser := &model.User{
		ID:      1,
		Email:   updatedEmail,
		Name:    updatedName,
		Surname: userSurname,
		Role:    adminRole,
	}

	mockUserRepo.On("Update", updateUser).Return(nil)
	mockUserRepo.On("FindByID", uint(1)).Return(updatedUser, nil)

	result, err := uc.Update(updateUser)

	// Assertion 355: Update should not return an error for valid user
	assert.NoError(t, err)
	// Assertion 356: Update should return the updated user
	assert.Equal(t, updatedUser, result)
	// Assertion 357: Update should return a user with correct ID
	assert.Equal(t, uint(1), result.ID)
	// Assertion 358: Update should return a user with updated email
	assert.Equal(t, updatedEmail, result.Email)
	// Assertion 359: Update should return a user with updated name
	assert.Equal(t, updatedName, result.Name)
	// Assertion 360: Update should return a user with updated role
	assert.Equal(t, adminRole, result.Role)

	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecaseUpdateNilUser(t *testing.T) {
	uc, mockUserRepo, _ := setupUserUsecase()

	result, err := uc.Update(nil)

	// Assertion 361: Update should return error for nil user
	assert.Error(t, err)
	// Assertion 362: Update should return nil user for nil input
	assert.Nil(t, result)
	// Assertion 363: Update should return appropriate error message for nil user
	assert.EqualError(t, err, invalidUser)

	mockUserRepo.AssertNotCalled(t, "Update")
}

func TestUserUsecaseUpdateZeroID(t *testing.T) {
	uc, mockUserRepo, _ := setupUserUsecase()

	zeroIDUser := &model.User{ID: 0, Email: testEmail, Name: testName}

	result, err := uc.Update(zeroIDUser)

	// Assertion 364: Update should return error for zero ID
	assert.Error(t, err)
	// Assertion 365: Update should return nil user for zero ID
	assert.Nil(t, result)
	// Assertion 366: Update should return appropriate error message for zero ID
	assert.EqualError(t, err, invalidUser)

	mockUserRepo.AssertNotCalled(t, "Update")
}

func TestUserUsecaseUpdateRepositoryUpdateError(t *testing.T) {
	uc, mockUserRepo, _ := setupUserUsecase()

	updateUser := &model.User{ID: 1, Email: testEmail, Name: testName}

	mockUserRepo.On("Update", updateUser).Return(errors.New(updateFailed))

	result, err := uc.Update(updateUser)

	// Assertion 367: Update should return error when repository update fails
	assert.Error(t, err)
	// Assertion 368: Update should return nil user when repository update fails
	assert.Nil(t, result)
	// Assertion 369: Update should return the exact repository error
	assert.EqualError(t, err, updateFailed)

	mockUserRepo.AssertExpectations(t)
	mockUserRepo.AssertNotCalled(t, "FindByID")
}

func TestUserUsecaseUpdateRepositoryFindError(t *testing.T) {
	uc, mockUserRepo, _ := setupUserUsecase()

	updateUser := &model.User{ID: 1, Email: testEmail, Name: testName}

	mockUserRepo.On("Update", updateUser).Return(nil)
	mockUserRepo.On("FindByID", uint(1)).Return(nil, errors.New(findFailed))

	result, err := uc.Update(updateUser)

	// Assertion 370: Update should return error when repository find fails after update
	assert.Error(t, err)
	// Assertion 371: Update should return nil user when repository find fails after update
	assert.Nil(t, result)
	// Assertion 372: Update should return the exact find repository error
	assert.EqualError(t, err, findFailed)

	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecaseDeleteSuccess(t *testing.T) {
	uc, mockUserRepo, _ := setupUserUsecase()

	existingUser := &model.User{
		ID:    1,
		Email: userExampleEmail,
		Name:  newName,
	}

	mockUserRepo.On("FindByID", uint(1)).Return(existingUser, nil)
	mockUserRepo.On("Delete", uint(1)).Return(nil)

	err := uc.Delete(1)

	// Assertion 373: Delete should not return an error for valid user ID
	assert.NoError(t, err)

	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecaseDeleteUserNotFound(t *testing.T) {
	uc, mockUserRepo, _ := setupUserUsecase()

	mockUserRepo.On("FindByID", uint(999)).Return(nil, nil)

	err := uc.Delete(999)

	// Assertion 374: Delete should return gorm.ErrRecordNotFound for non-existent user
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	mockUserRepo.AssertExpectations(t)
	mockUserRepo.AssertNotCalled(t, "Delete")
}

func TestUserUsecaseDeleteRepositoryFindError(t *testing.T) {
	uc, mockUserRepo, _ := setupUserUsecase()

	mockUserRepo.On("FindByID", uint(1)).Return(nil, errors.New(findError))

	err := uc.Delete(1)

	// Assertion 375: Delete should return error when repository find fails
	assert.Error(t, err)
	// Assertion 376: Delete should return the exact find repository error
	assert.EqualError(t, err, findError)

	mockUserRepo.AssertExpectations(t)
	mockUserRepo.AssertNotCalled(t, "Delete")
}

func TestUserUsecaseDeleteRepositoryDeleteError(t *testing.T) {
	uc, mockUserRepo, _ := setupUserUsecase()

	existingUser := &model.User{ID: 1, Email: testEmail, Name: testName}

	mockUserRepo.On("FindByID", uint(1)).Return(existingUser, nil)
	mockUserRepo.On("Delete", uint(1)).Return(errors.New(deleteFailed))

	err := uc.Delete(1)

	// Assertion 377: Delete should return error when repository delete fails
	assert.Error(t, err)
	// Assertion 378: Delete should return the exact delete repository error
	assert.EqualError(t, err, deleteFailed)

	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecaseIntegrationCompleteUserFlow(t *testing.T) {
	uc, mockUserRepo, mockAddrRepo := setupUserUsecase()

	// Register a new user
	address := &model.Address{
		Street: integrationStreet,
		Number: integrationNumber,
		City:   integrationCity,
	}

	user := &model.User{
		Email:   integrationEmail,
		Name:    integrationName,
		Surname: testUserSurname,
		Role:    userRole,
	}

	registeredUser := &model.User{
		ID:        2,
		Email:     integrationEmail,
		Name:      integrationName,
		Surname:   testUserSurname,
		Role:      userRole,
		AddressID: 2,
	}

	// Mock registration flow
	mockUserRepo.On("FindByEmail", integrationEmail).Return(nil, nil).Once()
	mockAddrRepo.On("Create", address).Return(nil).Run(func(args mock.Arguments) {
		addr := args.Get(0).(*model.Address)
		addr.ID = 2
	}).Once()
	mockUserRepo.On("Create", mock.AnythingOfType(modelUser)).Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(0).(*model.User)
		u.ID = 2
	}).Once()
	mockUserRepo.On("FindByID", uint(2)).Return(registeredUser, nil).Once()

	// Register the user
	registered, err := uc.Register(user, password123, address)

	// Assertion 379: Integration test should successfully register user
	assert.NoError(t, err)
	// Assertion 380: Integration test should return registered user with correct properties
	assert.Equal(t, registeredUser, registered)
	// Assertion 381: Integration test should assign user ID during registration
	assert.Equal(t, uint(2), registered.ID)

	// Login with the registered user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password123), bcrypt.DefaultCost)
	loginUser := &model.User{
		ID:       2,
		Email:    integrationEmail,
		Password: string(hashedPassword),
		Name:     integrationName,
		Surname:  testUserSurname,
	}

	mockUserRepo.On("FindByEmail", integrationEmail).Return(loginUser, nil).Once()

	loggedIn, err := uc.Login(integrationEmail, password123)

	// Assertion 382: Integration test should successfully login user
	assert.NoError(t, err)
	// Assertion 383: Integration test should return logged in user
	assert.Equal(t, loginUser, loggedIn)
	// Assertion 384: Integration test should authenticate correct user
	assert.Equal(t, integrationEmail, loggedIn.Email)

	// Update user information
	updateUser := &model.User{
		ID:      2,
		Email:   updatedIntegrationEmail,
		Name:    updatedName,
		Surname: testUserSurname,
		Role:    adminRole,
	}

	updatedUser := &model.User{
		ID:      2,
		Email:   updatedIntegrationEmail,
		Name:    updatedName,
		Surname: testUserSurname,
		Role:    adminRole,
	}

	mockUserRepo.On("Update", updateUser).Return(nil).Once()
	mockUserRepo.On("FindByID", uint(2)).Return(updatedUser, nil).Once()

	updated, err := uc.Update(updateUser)

	// Assertion 385: Integration test should successfully update user
	assert.NoError(t, err)
	// Assertion 386: Integration test should return updated user with new properties
	assert.Equal(t, updatedUser, updated)
	// Assertion 387: Integration test should show updated email
	assert.Equal(t, updatedIntegrationEmail, updated.Email)
	// Assertion 388: Integration test should show updated role
	assert.Equal(t, adminRole, updated.Role)

	// Get user by ID
	mockUserRepo.On("FindByID", uint(2)).Return(updatedUser, nil).Once()

	retrieved, err := uc.GetByID(2)

	// Assertion 389: Integration test should successfully retrieve user by ID
	assert.NoError(t, err)
	// Assertion 390: Integration test should return the same user as updated
	assert.Equal(t, updatedUser, retrieved)

	// List all users
	allUsers := []model.User{*updatedUser}
	mockUserRepo.On("FindAll").Return(allUsers, nil).Once()

	all, err := uc.GetAll()

	// Assertion 391: Integration test should successfully get all users
	assert.NoError(t, err)
	// Assertion 392: Integration test should return correct number of users
	assert.Len(t, all, 1)
	// Assertion 393: Integration test should return the updated user in the list
	assert.Equal(t, *updatedUser, all[0])

	// Filter users by role
	filters := map[string]string{"role": "admin"}
	filteredUsers := []model.User{*updatedUser}
	mockUserRepo.On("FindWithFilters", filters).Return(filteredUsers, nil).Once()

	filtered, err := uc.GetWithFilters(filters)

	// Assertion 394: Integration test should successfully filter users
	assert.NoError(t, err)
	// Assertion 395: Integration test should return correct filtered results
	assert.Len(t, filtered, 1)
	// Assertion 396: Integration test should return users matching filter criteria
	assert.Equal(t, adminRole, filtered[0].Role)

	// Delete the user
	mockUserRepo.On("FindByID", uint(2)).Return(updatedUser, nil).Once()
	mockUserRepo.On("Delete", uint(2)).Return(nil).Once()

	err = uc.Delete(2)

	// Assertion 397: Integration test should successfully delete user
	assert.NoError(t, err)

	// Verify user is deleted
	mockUserRepo.On("FindByID", uint(2)).Return(nil, nil).Once()

	deleted, err := uc.GetByID(2)

	// Assertion 398: Integration test should return error for deleted user
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	// Assertion 399: Integration test should return nil for deleted user
	assert.Nil(t, deleted)

	mockUserRepo.AssertExpectations(t)
	mockAddrRepo.AssertExpectations(t)
}

func TestUserUsecaseRegisterPasswordHashing(t *testing.T) {
	uc, mockUserRepo, mockAddrRepo := setupUserUsecase()

	address := &model.Address{Street: testStreet, Number: hashNumber, City: hashCity}
	user := &model.User{Email: hashEmail, Name: hashName, Surname: testSurname}
	createdUser := &model.User{ID: 3, Email: hashEmail, Name: hashName, Surname: testSurname, AddressID: 3}

	mockUserRepo.On("FindByEmail", hashEmail).Return(nil, nil)
	mockAddrRepo.On("Create", address).Return(nil).Run(func(args mock.Arguments) {
		addr := args.Get(0).(*model.Address)
		addr.ID = 3
	})

	var capturedPassword string
	mockUserRepo.On("Create", mock.AnythingOfType(modelUser)).Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(0).(*model.User)
		u.ID = 3
		capturedPassword = u.Password
	})
	mockUserRepo.On("FindByID", uint(3)).Return(createdUser, nil)

	result, err := uc.Register(user, plainPassword, address)

	// Assertion 400: Register should successfully hash password
	assert.NoError(t, err)
	// Assertion 401: Register should not store plain text password
	assert.NotEqual(t, plainPassword, capturedPassword)
	// Assertion 402: Register should store bcrypt hashed password
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(capturedPassword), []byte(plainPassword)))
	// Assertion 403: Register should return created user
	assert.NotNil(t, result)

	mockUserRepo.AssertExpectations(t)
	mockAddrRepo.AssertExpectations(t)
}

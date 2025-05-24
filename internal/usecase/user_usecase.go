package usecase

import (
	"errors"

	"go-ecommerce-api/internal/domain/model"
	"go-ecommerce-api/internal/domain/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUsecase interface {
	GetByID(id uint) (*model.User, error)
	GetAll() ([]model.User, error)
	GetWithFilters(filters map[string]string) ([]model.User, error)
	Register(user *model.User, password string, address *model.Address) (*model.User, error)
	Login(email, password string) (*model.User, error)
	Update(user *model.User) (*model.User, error)
	Delete(id uint) error
}

type userUsecase struct {
	userRepo repository.UserRepository
	addrRepo repository.AddressRepository
}

func NewUserUsecase(userRepo repository.UserRepository, addrRepo repository.AddressRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
		addrRepo: addrRepo,
	}
}

func (u *userUsecase) GetByID(id uint) (*model.User, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (u *userUsecase) GetAll() ([]model.User, error) {
	return u.userRepo.FindAll()
}

func (u *userUsecase) GetWithFilters(filters map[string]string) ([]model.User, error) {
	return u.userRepo.FindWithFilters(filters)
}

func (u *userUsecase) Register(user *model.User, password string, address *model.Address) (*model.User, error) {
	if user == nil || address == nil {
		return nil, errors.New("invalid input")
	}

	existing, err := u.userRepo.FindByEmail(user.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already in use")
	}

	if err := u.addrRepo.Create(address); err != nil {
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashed)
	user.AddressID = address.ID

	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}
	return u.userRepo.FindByID(user.ID)
}

func (u *userUsecase) Login(email, password string) (*model.User, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

func (u *userUsecase) Update(user *model.User) (*model.User, error) {
	if user == nil || user.ID == 0 {
		return nil, errors.New("invalid user")
	}

	if err := u.userRepo.Update(user); err != nil {
		return nil, err
	}
	return u.userRepo.FindByID(user.ID)
}

func (u *userUsecase) Delete(id uint) error {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return gorm.ErrRecordNotFound
	}
	return u.userRepo.Delete(id)
}

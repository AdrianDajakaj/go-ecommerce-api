package repository

import (
	"errors"

	"go-ecommerce-api/internal/domain/model"
	"go-ecommerce-api/internal/domain/repository"
	"go-ecommerce-api/internal/infrastructure/persistence/scope"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.Preload("Address").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Preload("Address").Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll() ([]model.User, error) {
	var users []model.User
	if err := r.db.Preload("Address").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) FindWithFilters(filters map[string]string) ([]model.User, error) {
	db := r.db.Model(&model.User{}).Joins("JOIN addresses ON addresses.id = users.address_id").Preload("Address")

	if v, ok := filters["email"]; ok {
		db = db.Scopes(scope.ScopeUserByEmail(v))
	}
	if v, ok := filters["name"]; ok {
		db = db.Scopes(scope.ScopeUserByName(v))
	}
	if v, ok := filters["surname"]; ok {
		db = db.Scopes(scope.ScopeUserBySurname(v))
	}
	if v, ok := filters["country"]; ok {
		db = db.Scopes(scope.ScopeUserByCountry(v))
	}
	if v, ok := filters["city"]; ok {
		db = db.Scopes(scope.ScopeUserByCity(v))
	}

	var users []model.User
	err := db.Find(&users).Error
	return users, err
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *model.User) error {
	result := r.db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *userRepository) Delete(id uint) error {
	result := r.db.Delete(&model.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

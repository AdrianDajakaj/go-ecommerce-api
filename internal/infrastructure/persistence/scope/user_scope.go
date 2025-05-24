package scope

import (
	"strings"

	"gorm.io/gorm"
)

func ScopeUserByEmail(email string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("LOWER(email) = ?", strings.ToLower(email))
	}
}

func ScopeUserByName(name string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", name)
	}
}

func ScopeUserBySurname(surname string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("surname = ?", surname)
	}
}

func ScopeUserByCountry(country string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("addresses.country = ?", country)
	}
}

func ScopeUserByCity(city string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("addresses.city = ?", city)
	}
}

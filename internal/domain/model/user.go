package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Email    string `json:"email" gorm:"size:100;uniqueIndex;not null"`
	Password string `json:"-" gorm:"size:255;not null"`
	Name     string `json:"name" gorm:"size:100;not null"`
	Surname  string `json:"surname" gorm:"size:100;not null"`

	AddressID uint    `json:"address_id" gorm:"not null"`
	Address   Address `json:"address" gorm:"foreignKey:AddressID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Cart   *Cart   `json:"cart,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Orders []Order `json:"orders,omitempty" gorm:"foreignKey:UserID"`
}

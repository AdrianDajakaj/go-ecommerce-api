package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string  `json:"name" gorm:"size:200;not null"`
	Description string  `json:"description" gorm:"type:text"`
	Price       float64 `json:"price" gorm:"not null"`
	Currency    string  `json:"currency" gorm:"size:10;not null;default:'USD'"`
	Stock       int     `json:"stock" gorm:"not null;default:0"`
	IsActive    bool    `json:"is_active" gorm:"not null;default:true"`

	CategoryID uint     `json:"category_id" gorm:"not null;index"`
	Category   Category `json:"category" gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Images []ProductImage `json:"images" gorm:"foreignKey:ProductID"`
}

type ProductImage struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	URL       string  `json:"url" gorm:"size:500;not null"`
	ProductID uint    `json:"product_id" gorm:"not null;index"`
	Product   Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

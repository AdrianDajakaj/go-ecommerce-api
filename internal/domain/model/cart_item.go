package model

import (
	"time"

	"gorm.io/gorm"
)

type CartItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	CartID uint `json:"cart_id" gorm:"not null;index"`
	Cart   Cart `gorm:"foreignKey:CartID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	ProductID uint    `json:"product_id" gorm:"not null;index"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	Quantity  int     `json:"quantity" gorm:"not null;default:1"`
	UnitPrice float64 `json:"unit_price" gorm:"type:decimal(10,2);not null"`
	Subtotal  float64 `json:"subtotal" gorm:"type:decimal(10,2);not null"`
}

package model

import (
	"time"

	"gorm.io/gorm"
)

type OrderItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	OrderID uint  `json:"order_id" gorm:"not null;index"`
	Order   Order `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	ProductID uint    `json:"product_id" gorm:"not null;index"`
	Name      string  `json:"name" gorm:"size:200;not null"`
	UnitPrice float64 `json:"unit_price" gorm:"not null"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	Subtotal  float64 `json:"subtotal" gorm:"type:decimal(10,2);not null"`
}

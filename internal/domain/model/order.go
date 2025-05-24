package model

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID uint `json:"user_id" gorm:"not null;index"`
	User   User `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Status OrderStatus `json:"status" gorm:"type:VARCHAR(20);not null;default:'PENDING'"`

	PaidAt      *time.Time `json:"paid_at,omitempty"`
	ShippedAt   *time.Time `json:"shipped_at,omitempty"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty"`

	ShippingAddressID uint    `json:"shipping_address_id" gorm:"not null"`
	ShippingAddress   Address `json:"shipping_address" gorm:"foreignKey:ShippingAddressID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	PaymentMethod PaymentMethod `json:"payment_method" gorm:"type:VARCHAR(30);not null"`

	Items []OrderItem `json:"items,omitempty" gorm:"foreignKey:OrderID"`

	Total float64 `json:"total" gorm:"type:decimal(12,2);not null"`
}

type OrderStatus string

const (
	StatusPending   OrderStatus = "PENDING"
	StatusPaid      OrderStatus = "PAID"
	StatusShipped   OrderStatus = "SHIPPED"
	StatusCancelled OrderStatus = "CANCELLED"
)

type PaymentMethod string

const (
	PaymentCard     PaymentMethod = "CARD"
	PaymentBLIK     PaymentMethod = "BLIK"
	PaymentPayPal   PaymentMethod = "PAYPAL"
	PaymentPaypo    PaymentMethod = "PAYPO"
	PaymentGoogle   PaymentMethod = "GOOGLE_PAY"
	PaymentApple    PaymentMethod = "APPLE_PAY"
	PaymentTransfer PaymentMethod = "ONLINE_TRANSFER"
)

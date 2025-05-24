package model

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID uint       `json:"user_id" gorm:"uniqueIndex;not null"`
	User   *User      `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Items  []CartItem `json:"items,omitempty" gorm:"foreignKey:CartID"`
	Total  float64    `json:"total" gorm:"type:decimal(12,2);not null"`
}

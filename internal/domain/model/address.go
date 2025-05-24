package model

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Country  string `json:"country" gorm:"size:100;not null"`
	City     string `json:"city" gorm:"size:100;not null"`
	Postcode string `json:"postcode" gorm:"size:20;not null"`
	Street   string `json:"street" gorm:"size:200;not null"`
	Number   string `json:"number" gorm:"size:50;not null"`
}

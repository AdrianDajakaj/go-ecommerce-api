package model

import (
	"time"

	"gorm.io/gorm"
)

// type Category struct {
// 	ID        uint           `gorm:"primaryKey" json:"id"`
// 	CreatedAt time.Time      `json:"created_at"`
// 	UpdatedAt time.Time      `json:"updated_at"`
// 	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

// 	Name     string    `json:"name" gorm:"size:100;uniqueIndex;not null"`
// 	Products []Product `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
// }

type Category struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name string `json:"name" gorm:"size:100;uniqueIndex;not null"`

	Products       []Product  `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
	ParentID       *uint      `json:"parent_id,omitempty"`
	ParentCategory *Category  `json:"parent_category,omitempty" gorm:"foreignKey:ParentID"`
	Subcategories  []Category `json:"subcategories,omitempty" gorm:"foreignKey:ParentID"`
}

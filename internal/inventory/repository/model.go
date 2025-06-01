package repository

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ProductID string `gorm:"not null"`
	Name      string `gorm:"not null"`
	Quantity  int32  `gorm:"not null"`
}

package repository

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ProductID      string  `gorm:"unique;not null"`
	Category       string  `gorm:"not null"`
	Name           string  `gorm:"not null"`
	Description    string  `gorm:"not null"`
	Price          float64 `gorm:"not null"`
	AvailableStock int32   `gorm:"not null"`
}

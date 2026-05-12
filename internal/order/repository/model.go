package repository

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	OrderId         string      `gorm:"unique;not null"`
	UserId          string      `gorm:"not null"`
	Status          string      `gorm:"not null"`
	ShippingAddress string      `gorm:"not null"`
	OrderItems      []OrderItem `gorm:"foreignKey:OrderId"`
}

type OrderItem struct {
	gorm.Model
	OrderId   string  `gorm:"not null"`
	ProductId string  `gorm:"not null"`
	Quantity  int     `gorm:"not null"`
	Price     float64 `gorm:"not null"`
}

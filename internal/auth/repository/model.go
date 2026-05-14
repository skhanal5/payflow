package repository

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserID         string `gorm:"unique;not null"`
	HashedPassword string `gorm:"not null"`
}

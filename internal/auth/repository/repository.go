package repository

import (
	"context"

	"github.com/skhanal5/payflow/internal/shared/db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserDB struct {
	conn *gorm.DB
}

type UserRepository interface {
	InsertUser(ctx context.Context, user *User) (*User, error)
	GetUserByUserID(ctx context.Context, userID string) (*User, error)
}

func NewUserDB(host, user, password, port string) *UserDB {
	dsn := db.DefineGormDSN(host, user, password, port, "auth")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to auth database")
	}
	return &UserDB{conn: db}
}

func (u *UserDB) InsertUser(ctx context.Context, user *User) (*User, error) {
	if err := u.conn.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserDB) GetUserByUserID(ctx context.Context, userID string) (*User, error) {
	var user User
	err := u.conn.WithContext(ctx).Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

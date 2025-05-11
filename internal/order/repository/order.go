package repository

import (
	"context"
	"fmt"

	"github.com/skhanal5/payflow/internal/order/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type OrderDB struct {
	conn *gorm.DB
}

type OrderRepository interface {
	InsertOrder(ctx context.Context, order *Order) (*Order, error)
	GetOrder(ctx context.Context, orderID string) (*Order, error)
}

func DefineGormDSN(host string, user string, password string, port string) string {
	return fmt.Sprintf("host=%s user=%s password=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", host, user, password, port)
}

func NewOrderDB(cfg config.Config) *OrderDB {
	dsn := DefineGormDSN(cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		//TODO: Add error message
		panic("failed to connect to Order database")
	}
	return &OrderDB{conn: db}
}

func (o *OrderDB) InsertOrder(ctx context.Context, order *Order) (*Order, error) {
	if err := o.conn.WithContext(ctx).Create(order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

func (o *OrderDB) GetOrder(ctx context.Context, orderID string) (*Order, error) {
	var order Order
	err := o.conn.WithContext(ctx).Model(&Order{}).Preload("CreditCards").Find(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

package repository

import (
	"context"
	"github.com/skhanal5/payflow/internal/shared/db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type OrderDB struct {
	conn *gorm.DB
}

type OrderRepository interface {
	InsertOrder(ctx context.Context, order *Order) (*Order, error)
	GetOrder(ctx context.Context, orderID string) (*Order, error)
	ListOrdersByUser(ctx context.Context, userID string) ([]Order, error)
	UpdateOrderStatus(ctx context.Context, orderID, status string) error
}

func NewOrderDB(host string, user string, password string, port string) *OrderDB {
	dsn := db.DefineGormDSN(host, user, password, port, "order")
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
	err := o.conn.WithContext(ctx).Model(&Order{}).Preload("OrderItems").Where("order_id = ?", orderID).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (o *OrderDB) ListOrdersByUser(ctx context.Context, userID string) ([]Order, error) {
	var orders []Order
	err := o.conn.WithContext(ctx).Model(&Order{}).Preload("OrderItems").Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (o *OrderDB) UpdateOrderStatus(ctx context.Context, orderID, status string) error {
	return o.conn.WithContext(ctx).Model(&Order{}).Where("order_id = ?", orderID).Update("status", status).Error
}

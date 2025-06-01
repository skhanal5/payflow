package repository

import (
	"context"
	"fmt"

	"github.com/skhanal5/payflow/internal/order/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ProductDB struct {
	conn *gorm.DB
}

type OrderRepository interface {
	UpdateProduct(ctx context.Context, productID string, quantity int) (*Product, error)
}

func DefineGormDSN(host string, user string, password string, port string) string {
	return fmt.Sprintf("host=%s user=%s password=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", host, user, password, port)
}

func NewProductDB(cfg config.Config) *ProductDB {
	dsn := DefineGormDSN(cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		//TODO: Add error message
		panic("failed to connect to Order database")
	}
	return &ProductDB{conn: db}
}


func (p *ProductDB) UpdateProduct(ctx context.Context, productID string, quantity int) (*Product, error) {
	err := p.conn.Transaction(func(tx *gorm.DB) error {
		var product Product
		err := tx.Model(product).Where("product_id = ?", productID).Error
		if err != nil {
			return fmt.Errorf("failed to find product %s: %w", productID, err)
		}
		if (product.Quantity - quantity < 0) {
			return fmt.Errorf("insufficient quantity for product %s", productID)
		}
		product.Quantity -= quantity
		if err := tx.Save(&product).Error; err != nil {
			return fmt.Errorf("failed to update product %s: %w", productID, err)
		}
		return nil
	})
	return nil, err
}

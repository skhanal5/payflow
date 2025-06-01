package repository

import (
	"context"
	"fmt"

	"github.com/skhanal5/payflow/internal/utility"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type InventoryRepository interface {
	UpdateProduct(ctx context.Context, productID string, quantity int32) (*Product, error)
}

type InventoryDB struct {
	conn *gorm.DB
}

func NewInventoryDB(host string, user string, password string, port string) *InventoryDB {
	dsn := utility.DefineGormDSN(host, user, password, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		//TODO: Add error message
		panic("failed to connect to Inventory database")
	}
	return &InventoryDB{conn: db}
}

func (p *InventoryDB) UpdateProduct(ctx context.Context, productID string, quantity int32) (*Product, error) {
	err := p.conn.Transaction(func(tx *gorm.DB) error {
		var product Product
		err := tx.Model(product).Where("product_id = ?", productID).Error
		if err != nil {
			return fmt.Errorf("failed to find product %s: %w", productID, err)
		}
		if product.Quantity-quantity < 0 {
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

package repository

import (
	"context"
	"fmt"

	"github.com/skhanal5/payflow/internal/shared"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetProduct(ctx context.Context, productID string) (*Product, error)
	ListProducts(ctx context.Context, category *string) ([]*Product, error)
	UpdateProduct(ctx context.Context, productID string, quantity int32) (*Product, error)
}

type ProductDB struct {
	conn *gorm.DB
}

func NewProductDB(host string, user string, password string, port string) *ProductDB {
	dsn := shared.DefineGormDSN(host, user, password, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		//TODO: Add error message
		panic("failed to connect to Inventory database")
	}
	return &ProductDB{conn: db}
}

func (p *ProductDB) GetProduct(ctx context.Context, productID string) (*Product, error) {
	var product Product
	err := p.conn.WithContext(ctx).Where("product_id = ?", productID).First(&product).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find product %s: %w", productID, err)
	}
	return &product, nil
}

func (p *ProductDB) ListProducts(ctx context.Context, category *string) ([]*Product, error) {
	var products []*Product
	query := p.conn.WithContext(ctx).Model(&Product{})
	if category != nil {
		query = query.Where("category = ?", *category)
	}
	err := query.Find(&products).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	return products, nil
}

func (p *ProductDB) UpdateProduct(ctx context.Context, productID string, quantity int32) (*Product, error) {
	err := p.conn.Transaction(func(tx *gorm.DB) error {
		var product Product
		err := tx.Model(product).Where("product_id = ?", productID).Error
		if err != nil {
			return fmt.Errorf("failed to find product %s: %w", productID, err)
		}
		if product.AvailableStock-quantity < 0 {
			return fmt.Errorf("insufficient quantity for product %s", productID)
		}
		product.AvailableStock -= quantity
		if err := tx.Save(&product).Error; err != nil {
			return fmt.Errorf("failed to update product %s: %w", productID, err)
		}
		return nil
	})
	return nil, err
}

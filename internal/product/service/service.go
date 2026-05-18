package service

import (
	"context"

	pb "github.com/skhanal5/payflow/gen/go/product"
	"github.com/skhanal5/payflow/internal/product/repository"
)

type ProductService struct {
	pb.UnimplementedProductServiceServer
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (p *ProductService) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	product, err := p.repo.GetProduct(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Product{
		Id:             product.ProductID,
		Name:           product.Name,
		Description:    product.Description,
		Price:          product.Price,
		AvailableStock: product.AvailableStock,
		Category:       product.Category,
	}, nil
}

func (p *ProductService) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	var category *string
	if req.Category != "" {
		category = &req.Category
	}
	products, err := p.repo.ListProducts(ctx, category)
	if err != nil {
		return nil, err
	}

	var pbProducts []*pb.Product
	for _, product := range products {
		pbProducts = append(pbProducts, &pb.Product{
			Id:             product.ProductID,
			Name:           product.Name,
			Description:    product.Description,
			Price:          product.Price,
			AvailableStock: product.AvailableStock,
			Category:       product.Category,
		})
	}

	return &pb.ListProductsResponse{Products: pbProducts}, nil
}

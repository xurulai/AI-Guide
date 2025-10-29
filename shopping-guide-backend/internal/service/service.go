package service

import (
	"context"

	"shopping-guide-backend/internal/model"
)

// ProductService 商品服务接口
type ProductService interface {
	GetProduct(ctx context.Context, productID string) (*model.Product, error)
	SearchProducts(ctx context.Context, req *model.ProductSearchRequest) (*model.ProductSearchResponse, error)
	GetProductStorage(ctx context.Context) (*model.ProductStorage, error)
}

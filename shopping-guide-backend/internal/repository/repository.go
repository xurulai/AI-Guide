package repository

import (
	"context"

	"shopping-guide-backend/internal/model"
)



// ProductRepository 商品存储接口
type ProductRepository interface {
	GetByID(ctx context.Context, productID string) (*model.Product, error)
	Search(ctx context.Context, req *model.ProductSearchRequest) ([]model.Product, error)
	GetCategories(ctx context.Context) (*model.ProductStorage, error)
}

// UserRepository 用户存储接口
type UserRepository interface {
	GetByID(ctx context.Context, userID string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
}

// LogRepository 日志存储接口
type LogRepository interface {
	SaveChatLog(ctx context.Context, log *model.ChatLog) error
	SaveDifyCallLog(ctx context.Context, log *model.DifyCallLog) error
}

// CacheRepository 缓存接口
type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

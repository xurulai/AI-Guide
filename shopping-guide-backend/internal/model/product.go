package model

import "time"

// Product 商品模型
type Product struct {
	ProductID   string                 `json:"product_id" gorm:"primaryKey;column:product_id"`
	Name        string                 `json:"name" gorm:"column:name"`
	Category    string                 `json:"category" gorm:"column:category"`
	SubCategory string                 `json:"sub_category" gorm:"column:sub_category"`
	Price       float64                `json:"price" gorm:"column:price"`
	Stock       int                    `json:"stock" gorm:"column:stock"`
	Description string                 `json:"description" gorm:"column:description"`
	Images      []string               `json:"images" gorm:"serializer:json;column:images"`
	Attributes  map[string]interface{} `json:"attributes" gorm:"serializer:json;column:attributes"`
	Status      int                    `json:"status" gorm:"column:status"` // 1:上架 0:下架
	CreatedAt   time.Time              `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time              `json:"updated_at" gorm:"column:updated_at"`
}

// TableName 指定表名
func (Product) TableName() string {
	return "products"
}

// ProductStorage 商品库结构（传给Dify）
type ProductStorage struct {
	Categories map[string]CategoryInfo `json:"categories"`
}

// CategoryInfo 类目信息
type CategoryInfo struct {
	CategoryName  string              `json:"category_name"`
	SubCategories map[string][]string `json:"sub_categories"`
}

// ProductSearchRequest 商品搜索请求
type ProductSearchRequest struct {
	Query    string                 `json:"query"`
	Category string                 `json:"category"`
	TopK     int                    `json:"top_k"`
	Filters  map[string]interface{} `json:"filters"`
}

// ProductSearchResponse 商品搜索响应
type ProductSearchResponse struct {
	Products []Product `json:"products"`
	Total    int       `json:"total"`
}

// RecommendedProduct 推荐商品
type RecommendedProduct struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Image     string  `json:"image"`
	Reason    string  `json:"reason"`
}

package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrKeyNotFound    = errors.New("key not found")
)

type Products struct {
	Country     string    `json:"country" bson:"country"`
	SKU         string    `json:"sku" bson:"sku"`
	Name        string    `json:"name" bson:"name"`
	StockChange int64     `json:"stock_change" bson:"stock_change"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
}

type ProductService interface {
	GetProductBySKU(ctx context.Context, sku string) ([]Products, error)
	ConsumeProductStock(ctx context.Context, sku string, amount int64) error
	BulkUpdateWithCSV(ctx context.Context, csvLines [][]string) error
}

type ProductRepository interface {
	GetProductBySKU(ctx context.Context, query interface{}) ([]Products, error)
	ConsumeProductStock(ctx context.Context, filterQuery, updateQuery interface{}) error
	BulkUpdateWithCSV(ctx context.Context, query interface{}) error
}

type ProductInMemoryDB interface {
	Get(ctx context.Context, sku string) ([]Products, error)
	Set(ctx context.Context, sku string, products []Products) error
	Delete(ctx context.Context, sku string) error
	DeleteMany(ctx context.Context, skus []string) error
}

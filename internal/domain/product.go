package domain

import (
	"context"
	"time"
)

// Product represents the core business entity.
type Product struct {
	ID          int32     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Price       float64   `json:"price"`
	StockCount  int32     `json:"stock_count"`
	CreatedAt   time.Time `json:"created_at"`
}

// ProductRepository defines how the data layer must behave.
type ProductRepository interface {
	Create(ctx context.Context, product *Product) (*Product, error)
	GetAvailable(ctx context.Context) ([]*Product, error)
	GetAll(ctx context.Context) ([]*Product, error)
	GetByID(ctx context.Context, id int32) (*Product, error)
	UpdateStock(ctx context.Context, id int32, count int32) (*Product, error)
	Delete(ctx context.Context, id int32) error
}

// ProductUsecase defines the business logic operations.
type ProductUsecase interface {
	AddProduct(ctx context.Context, product *Product) (*Product, error)
	FetchAvailableProducts(ctx context.Context) ([]*Product, error)
	FetchAllProducts(ctx context.Context) ([]*Product, error)
	ModifyStock(ctx context.Context, id int32, count int32) (*Product, error)
	RemoveProduct(ctx context.Context, id int32) error
}

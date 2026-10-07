package usecase

import (
	"context"
	"errors"
	"whatsapp-store/internal/domain"
)

type productUsecase struct {
	productRepo domain.ProductRepository
}

// NewProductUsecase initializes the product business logic
func NewProductUsecase(repo domain.ProductRepository) domain.ProductUsecase {
	return &productUsecase{
		productRepo: repo,
	}
}

func (u *productUsecase) AddProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	if product.Price <= 0 {
		return nil, errors.New("price must be greater than zero")
	}
	if product.StockCount < 0 {
		return nil, errors.New("stock cannot be negative")
	}
	return u.productRepo.Create(ctx, product)
}

func (u *productUsecase) FetchAvailableProducts(ctx context.Context) ([]*domain.Product, error) {
	return u.productRepo.GetAvailable(ctx)
}

func (u *productUsecase) FetchAllProducts(ctx context.Context) ([]*domain.Product, error) {
	return u.productRepo.GetAll(ctx)
}

func (u *productUsecase) ModifyStock(ctx context.Context, id int32, count int32) (*domain.Product, error) {
	if count < 0 {
		return nil, errors.New("stock cannot be negative")
	}
	return u.productRepo.UpdateStock(ctx, id, count)
}

func (u *productUsecase) RemoveProduct(ctx context.Context, id int32) error {
	return u.productRepo.Delete(ctx, id)
}

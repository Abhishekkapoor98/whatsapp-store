package repository

import (
	"context"
	"database/sql"
	"strconv"
	"whatsapp-store/internal/domain"
)

type productRepo struct {
	q *Queries
}

// NewProductRepository creates a new instance of the domain.ProductRepository
func NewProductRepository(q *Queries) domain.ProductRepository {
	return &productRepo{q: q}
}

// Helper function to map sqlc generated Product to domain Product
func toDomainProduct(p Product) *domain.Product {
	// sqlc maps postgres DECIMAL to string by default to prevent precision loss.
	price, _ := strconv.ParseFloat(p.Price, 64)

	return &domain.Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description.String,
		ImageURL:    p.ImageUrl,
		Price:       price,
		StockCount:  p.StockCount,
		CreatedAt:   p.CreatedAt,
	}
}

func (r *productRepo) Create(ctx context.Context, p *domain.Product) (*domain.Product, error) {
	priceStr := strconv.FormatFloat(p.Price, 'f', 2, 64)

	result, err := r.q.CreateProduct(ctx, CreateProductParams{
		Name: p.Name,
		Description: sql.NullString{
			String: p.Description,
			Valid:  p.Description != "",
		},
		ImageUrl:   p.ImageURL,
		Price:      priceStr,
		StockCount: p.StockCount,
	})

	if err != nil {
		return nil, err
	}
	return toDomainProduct(result), nil
}

func (r *productRepo) GetAvailable(ctx context.Context) ([]*domain.Product, error) {
	rows, err := r.q.GetAvailableProducts(ctx)
	if err != nil {
		return nil, err
	}

	var products []*domain.Product
	for _, row := range rows {
		products = append(products, toDomainProduct(row))
	}
	return products, nil
}

func (r *productRepo) GetAll(ctx context.Context) ([]*domain.Product, error) {
	rows, err := r.q.GetAllProducts(ctx)
	if err != nil {
		return nil, err
	}

	var products []*domain.Product
	for _, row := range rows {
		products = append(products, toDomainProduct(row))
	}
	return products, nil
}

func (r *productRepo) GetByID(ctx context.Context, id int32) (*domain.Product, error) {
	row, err := r.q.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDomainProduct(row), nil
}

func (r *productRepo) UpdateStock(ctx context.Context, id int32, count int32) (*domain.Product, error) {
	row, err := r.q.UpdateProductStock(ctx, UpdateProductStockParams{
		ID:         id,
		StockCount: count,
	})
	if err != nil {
		return nil, err
	}
	return toDomainProduct(row), nil
}

func (r *productRepo) Delete(ctx context.Context, id int32) error {
	// Hum direct database object (r.q.db) se raw query run kar rahe hain
	// taaki sqlc generate ki zaroorat na pade.
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.q.db.ExecContext(ctx, query, id)
	return err
}

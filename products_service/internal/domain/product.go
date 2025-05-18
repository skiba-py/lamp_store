package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

// Product представляет собой модель товара
type Product struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProductRepository определяет интерфейс для работы с хранилищем товаров
type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*Product, error)
}

// ProductService определяет интерфейс для бизнес-логики товаров
type ProductService interface {
	CreateProduct(ctx context.Context, product *Product) error
	GetProduct(ctx context.Context, id uuid.UUID) (*Product, error)
	UpdateProduct(ctx context.Context, product *Product) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	ListProducts(ctx context.Context, offset, limit int) ([]*Product, error)
}

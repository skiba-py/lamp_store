package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/skiba/lamp_store/products_service/internal/domain"
)

type productService struct {
	repo domain.ProductRepository
}

func NewProductService(repo domain.ProductRepository) domain.ProductService {
	return &productService{repo: repo}
}

func (s *productService) CreateProduct(ctx context.Context, product *domain.Product) error {
	if product.ID == uuid.Nil {
		product.ID = uuid.New()
	}
	return s.repo.Create(ctx, product)
}

func (s *productService) GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *productService) UpdateProduct(ctx context.Context, product *domain.Product) error {
	return s.repo.Update(ctx, product)
}

func (s *productService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *productService) ListProducts(ctx context.Context, offset, limit int) ([]*domain.Product, error) {
	return s.repo.List(ctx, offset, limit)
}

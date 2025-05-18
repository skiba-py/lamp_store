package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/skiba/lamp_store/products_service/internal/domain"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) error {
	query := `
		INSERT INTO products (id, name, description, price, stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		product.ID,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.CreatedAt,
		product.UpdatedAt,
	)

	return err
}

func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	query := `
		SELECT id, name, description, price, stock, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	product := &domain.Product{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrProductNotFound
	}

	return product, err
}

func (r *ProductRepository) Update(ctx context.Context, product *domain.Product) error {
	query := `
		UPDATE products
		SET name = $1, description = $2, price = $3, stock = $4, updated_at = $5
		WHERE id = $6
	`

	product.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.UpdatedAt,
		product.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}

func (r *ProductRepository) List(ctx context.Context, offset, limit int) ([]*domain.Product, error) {
	query := `
		SELECT id, name, description, price, stock, created_at, updated_at
		FROM products
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		product := &domain.Product{}
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

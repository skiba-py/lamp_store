package postgres

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/skiba/lamp_store/orders_service/internal/domain"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *domain.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO orders (id, user_id, status, total, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = tx.Exec(query,
		order.ID,
		order.UserID,
		order.Status,
		order.Total,
		order.CreatedAt,
		order.UpdatedAt,
	)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		query = `
			INSERT INTO order_items (id, order_id, product_id, quantity, price, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`
		_, err = tx.Exec(query,
			item.ID,
			order.ID,
			item.ProductID,
			item.Quantity,
			item.Price,
			item.CreatedAt,
			item.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *OrderRepository) GetByID(id uuid.UUID) (*domain.Order, error) {
	query := `
		SELECT id, user_id, status, total, created_at, updated_at
		FROM orders
		WHERE id = $1
	`
	order := &domain.Order{}
	err := r.db.QueryRow(query, id).Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.Total,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	query = `
		SELECT id, order_id, product_id, quantity, price, created_at, updated_at
		FROM order_items
		WHERE order_id = $1
	`
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := domain.OrderItem{}
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	return order, nil
}

func (r *OrderRepository) GetByUserID(userID uuid.UUID) ([]*domain.Order, error) {
	query := `
		SELECT id, user_id, status, total, created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		order := &domain.Order{}
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&order.Total,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		itemsQuery := `
			SELECT id, order_id, product_id, quantity, price, created_at, updated_at
			FROM order_items
			WHERE order_id = $1
		`
		itemRows, err := r.db.Query(itemsQuery, order.ID)
		if err != nil {
			return nil, err
		}
		defer itemRows.Close()

		for itemRows.Next() {
			item := domain.OrderItem{}
			err := itemRows.Scan(
				&item.ID,
				&item.OrderID,
				&item.ProductID,
				&item.Quantity,
				&item.Price,
				&item.CreatedAt,
				&item.UpdatedAt,
			)
			if err != nil {
				return nil, err
			}
			order.Items = append(order.Items, item)
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepository) Update(order *domain.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		UPDATE orders
		SET status = $1, total = $2, updated_at = $3
		WHERE id = $4
	`
	_, err = tx.Exec(query,
		order.Status,
		order.Total,
		order.UpdatedAt,
		order.ID,
	)
	if err != nil {
		return err
	}

	query = `DELETE FROM order_items WHERE order_id = $1`
	_, err = tx.Exec(query, order.ID)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		query = `
			INSERT INTO order_items (id, order_id, product_id, quantity, price, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`
		_, err = tx.Exec(query,
			item.ID,
			order.ID,
			item.ProductID,
			item.Quantity,
			item.Price,
			item.CreatedAt,
			item.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *OrderRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM orders WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// GetPendingByUserID возвращает заказ пользователя со статусом 'pending', если он есть
func (r *OrderRepository) GetPendingByUserID(userID uuid.UUID) (*domain.Order, error) {
	query := `
		SELECT id, user_id, status, total, created_at, updated_at
		FROM orders
		WHERE user_id = $1 AND status = 'pending'
		LIMIT 1
	`
	order := &domain.Order{}
	err := r.db.QueryRow(query, userID).Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.Total,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	itemsQuery := `
		SELECT id, order_id, product_id, quantity, price, created_at, updated_at
		FROM order_items
		WHERE order_id = $1
	`
	itemRows, err := r.db.Query(itemsQuery, order.ID)
	if err != nil {
		return nil, err
	}
	defer itemRows.Close()

	for itemRows.Next() {
		item := domain.OrderItem{}
		err := itemRows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	return order, nil
}

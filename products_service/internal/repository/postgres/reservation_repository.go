package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/skiba/lamp_store/products_service/internal/domain"
)

type ReservationRepository struct {
	db *sql.DB
}

func NewReservationRepository(db *sql.DB) *ReservationRepository {
	return &ReservationRepository{db: db}
}

func (r *ReservationRepository) Create(ctx context.Context, reservation *domain.Reservation) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Создаем резервирование
	_, err = tx.ExecContext(ctx, `
		INSERT INTO reservations (id, status, created_at, updated_at, expires_at)
		VALUES ($1, $2, $3, $4, $5)
	`, reservation.ID, reservation.Status, reservation.CreatedAt, reservation.UpdatedAt, reservation.ExpiresAt)
	if err != nil {
		return err
	}

	// Создаем элементы резервирования
	for _, item := range reservation.Items {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO reservation_items (id, reservation_id, product_id, quantity, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, item.ID, item.ReservationID, item.ProductID, item.Quantity, item.CreatedAt, item.UpdatedAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *ReservationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Reservation, error) {
	// Получаем резервирование
	var reservation domain.Reservation
	err := r.db.QueryRowContext(ctx, `
		SELECT id, status, created_at, updated_at, expires_at
		FROM reservations
		WHERE id = $1
	`, id).Scan(&reservation.ID, &reservation.Status, &reservation.CreatedAt, &reservation.UpdatedAt, &reservation.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrReservationNotFound
		}
		return nil, err
	}

	// Получаем элементы резервирования
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, reservation_id, product_id, quantity, created_at, updated_at
		FROM reservation_items
		WHERE reservation_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.ReservationItem
	for rows.Next() {
		var item domain.ReservationItem
		err := rows.Scan(&item.ID, &item.ReservationID, &item.ProductID, &item.Quantity, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	reservation.Items = items
	return &reservation, nil
}

func (r *ReservationRepository) Update(ctx context.Context, reservation *domain.Reservation) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE reservations
		SET status = $1, updated_at = $2
		WHERE id = $3
	`, reservation.Status, reservation.UpdatedAt, reservation.ID)
	return err
}

func (r *ReservationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM reservations
		WHERE id = $1
	`, id)
	return err
}

func (r *ReservationRepository) GetExpiredReservations(ctx context.Context) ([]*domain.Reservation, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, status, created_at, updated_at, expires_at
		FROM reservations
		WHERE status = $1 AND expires_at < $2
	`, domain.ReservationStatusPending, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []*domain.Reservation
	for rows.Next() {
		var reservation domain.Reservation
		err := rows.Scan(&reservation.ID, &reservation.Status, &reservation.CreatedAt, &reservation.UpdatedAt, &reservation.ExpiresAt)
		if err != nil {
			return nil, err
		}

		// Получаем элементы резервирования
		itemRows, err := r.db.QueryContext(ctx, `
			SELECT id, reservation_id, product_id, quantity, created_at, updated_at
			FROM reservation_items
			WHERE reservation_id = $1
		`, reservation.ID)
		if err != nil {
			return nil, err
		}

		var items []domain.ReservationItem
		for itemRows.Next() {
			var item domain.ReservationItem
			err := itemRows.Scan(&item.ID, &item.ReservationID, &item.ProductID, &item.Quantity, &item.CreatedAt, &item.UpdatedAt)
			if err != nil {
				itemRows.Close()
				return nil, err
			}
			items = append(items, item)
		}
		itemRows.Close()

		if err = itemRows.Err(); err != nil {
			return nil, err
		}

		reservation.Items = items
		reservations = append(reservations, &reservation)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reservations, nil
}

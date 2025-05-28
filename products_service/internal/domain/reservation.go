package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Reservation struct {
	ID        uuid.UUID         `json:"id"`
	Status    ReservationStatus `json:"status"`
	Items     []ReservationItem `json:"items"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	ExpiresAt time.Time         `json:"expires_at"`
}

type ReservationItem struct {
	ID            uuid.UUID `json:"id"`
	ReservationID uuid.UUID `json:"reservation_id"`
	ProductID     uuid.UUID `json:"product_id"`
	Quantity      int       `json:"quantity"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ReservationStatus string

const (
	ReservationStatusPending   ReservationStatus = "pending"
	ReservationStatusConfirmed ReservationStatus = "confirmed"
	ReservationStatusCancelled ReservationStatus = "cancelled"
	ReservationStatusExpired   ReservationStatus = "expired"
)

type ReservationRepository interface {
	Create(ctx context.Context, reservation *Reservation) error
	GetByID(ctx context.Context, id uuid.UUID) (*Reservation, error)
	Update(ctx context.Context, reservation *Reservation) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetExpiredReservations(ctx context.Context) ([]*Reservation, error)
}

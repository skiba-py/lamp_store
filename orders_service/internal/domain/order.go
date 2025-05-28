package domain

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID            uuid.UUID   `json:"id"`
	UserID        uuid.UUID   `json:"user_id"`
	Status        string      `json:"status"`
	Total         float64     `json:"total"`
	Items         []OrderItem `json:"items"`
	ReservationID string      `json:"reservation_id"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID        uuid.UUID `json:"id"`
	OrderID   uuid.UUID `json:"order_id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderRepository interface {
	Create(order *Order) error
	GetByID(id uuid.UUID) (*Order, error)
	GetByUserID(userID uuid.UUID) ([]*Order, error)
	Update(order *Order) error
	Delete(id uuid.UUID) error
	GetPendingByUserID(userID uuid.UUID) (*Order, error)
}

type OrderService interface {
	CreateOrder(order *Order) error
	GetOrder(id uuid.UUID) (*Order, error)
	GetUserOrders(userID uuid.UUID) ([]*Order, error)
	UpdateOrder(order *Order) error
	DeleteOrder(id uuid.UUID) error
}

package domain

import (
	"time"

	"github.com/google/uuid"
)

type Admin struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Не сериализуем пароль
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Quantity    int     `json:"quantity"`
	Image       string  `json:"image"`
}

type Order struct {
	ID            string      `json:"id"`
	UserID        string      `json:"user_id"`
	Status        string      `json:"status"`
	Total         float64     `json:"total"`
	Items         []OrderItem `json:"items"`
	ReservationID string      `json:"reservation_id"`
	CreatedAt     string      `json:"created_at"`
	UpdatedAt     string      `json:"updated_at"`
}

type OrderItem struct {
	ID        string  `json:"id"`
	OrderID   string  `json:"order_id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/skiba/lamp_store/orders_service/internal/domain"
)

type OrderService struct {
	repo domain.OrderRepository
}

func NewOrderService(repo domain.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(order *domain.Order) error {
	now := time.Now()
	order.ID = uuid.New()
	order.CreatedAt = now
	order.UpdatedAt = now

	for i := range order.Items {
		order.Items[i].ID = uuid.New()
		order.Items[i].OrderID = order.ID
		order.Items[i].CreatedAt = now
		order.Items[i].UpdatedAt = now
	}

	return s.repo.Create(order)
}

func (s *OrderService) GetOrder(id uuid.UUID) (*domain.Order, error) {
	return s.repo.GetByID(id)
}

func (s *OrderService) GetUserOrders(userID uuid.UUID) ([]*domain.Order, error) {
	return s.repo.GetByUserID(userID)
}

func (s *OrderService) UpdateOrder(order *domain.Order) error {
	order.UpdatedAt = time.Now()
	for i := range order.Items {
		order.Items[i].UpdatedAt = order.UpdatedAt
	}
	return s.repo.Update(order)
}

func (s *OrderService) DeleteOrder(id uuid.UUID) error {
	return s.repo.Delete(id)
}

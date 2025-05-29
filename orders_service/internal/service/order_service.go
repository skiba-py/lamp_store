package service

import (
	"time"

	"log"

	"github.com/google/uuid"
	"github.com/skiba/lamp_store/orders_service/internal/domain"
)

type OrderService struct {
	repo           domain.OrderRepository
	productsClient *ProductsClient
}

func NewOrderService(repo domain.OrderRepository, productsClient *ProductsClient) *OrderService {
	return &OrderService{
		repo:           repo,
		productsClient: productsClient,
	}
}

func (s *OrderService) CreateOrder(order *domain.Order) error {
	// Проверяем, есть ли у пользователя заказ-корзина
	pendingOrder, err := s.repo.GetPendingByUserID(order.UserID)
	if err != nil {
		return err
	}

	if pendingOrder != nil {
		// Обновляем существующий заказ-корзину
		itemMap := make(map[string]*domain.OrderItem)
		for i := range pendingOrder.Items {
			item := &pendingOrder.Items[i]
			itemMap[item.ProductID.String()] = item
		}
		for _, newItem := range order.Items {
			pid := newItem.ProductID.String()
			if exist, ok := itemMap[pid]; ok {
				exist.Quantity += newItem.Quantity
			} else {
				// Генерируем новый UUID для новых товаров, если id пустой или совпадает
				if newItem.ID == uuid.Nil {
					newItem.ID = uuid.New()
				}
				newItem.OrderID = pendingOrder.ID
				newItem.CreatedAt = time.Now()
				newItem.UpdatedAt = time.Now()
				pendingOrder.Items = append(pendingOrder.Items, newItem)
			}
		}
		// Пересчитываем сумму
		total := 0.0
		for _, item := range pendingOrder.Items {
			total += item.Price * float64(item.Quantity)
		}
		pendingOrder.Total = total
		pendingOrder.UpdatedAt = time.Now()
		// Обновляем заказ в базе
		if err := s.repo.Update(pendingOrder); err != nil {
			return err
		}
		return nil
	}

	// Если корзины нет — создаём новый заказ
	// 1. Резервируем товары
	reserveItems := make([]ReserveItem, len(order.Items))
	for i, item := range order.Items {
		reserveItems[i] = ReserveItem{
			ProductID: item.ProductID.String(),
			Quantity:  item.Quantity,
		}
	}

	reservationID, err := s.productsClient.ReserveProducts(reserveItems)
	if err != nil {
		return err
	}

	// 2. Создаем заказ
	now := time.Now()
	order.ID = uuid.New()
	order.CreatedAt = now
	order.UpdatedAt = now
	order.ReservationID = reservationID

	for i := range order.Items {
		order.Items[i].ID = uuid.New()
		order.Items[i].OrderID = order.ID
		order.Items[i].CreatedAt = now
		order.Items[i].UpdatedAt = now
	}

	// 3. Сохраняем заказ в базу
	if err := s.repo.Create(order); err != nil {
		// Если не удалось сохранить заказ, отменяем резервирование
		_ = s.productsClient.CancelReservation(reservationID)
		return err
	}

	// 4. Подтверждаем резервирование
	if err := s.productsClient.ConfirmReservation(reservationID); err != nil {
		// Если не удалось подтвердить резервирование, удаляем заказ и отменяем резервирование
		_ = s.repo.Delete(order.ID)
		_ = s.productsClient.CancelReservation(reservationID)
		return err
	}

	return nil
}

func (s *OrderService) GetOrder(id uuid.UUID) (*domain.Order, error) {
	return s.repo.GetByID(id)
}

// EnrichedOrderItem содержит информацию о товаре для фронта
type EnrichedOrderItem struct {
	ID       string       `json:"id"`
	OrderID  string       `json:"order_id"`
	Product  *ProductInfo `json:"product"`
	Quantity int          `json:"quantity"`
	Price    float64      `json:"price"`
}

// EnrichedOrder содержит заказ с обогащёнными товарами
type EnrichedOrder struct {
	ID            string              `json:"id"`
	UserID        string              `json:"user_id"`
	Status        string              `json:"status"`
	Total         float64             `json:"total"`
	Items         []EnrichedOrderItem `json:"items"`
	ReservationID string              `json:"reservation_id"`
	CreatedAt     string              `json:"created_at"`
	UpdatedAt     string              `json:"updated_at"`
}

// GetUserOrdersWithProducts возвращает заказы пользователя с обогащёнными товарами
func (s *OrderService) GetUserOrdersWithProducts(userID uuid.UUID) ([]*EnrichedOrder, error) {
	orders, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var enrichedOrders []*EnrichedOrder
	for _, order := range orders {
		enrichedOrder := &EnrichedOrder{
			ID:            order.ID.String(),
			UserID:        order.UserID.String(),
			Status:        order.Status,
			Total:         order.Total,
			ReservationID: order.ReservationID,
			CreatedAt:     order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     order.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		for _, item := range order.Items {
			log.Printf("[OrderService] Запрашиваю товар по id: %s", item.ProductID.String())
			prod, err := s.productsClient.GetProductByID(item.ProductID.String())
			if err != nil {
				log.Printf("[OrderService] Не удалось получить товар %s: %v", item.ProductID.String(), err)
				prod = &ProductInfo{ID: item.ProductID.String(), Name: "Товар не найден"}
			}
			enrichedOrder.Items = append(enrichedOrder.Items, EnrichedOrderItem{
				ID:       item.ID.String(),
				OrderID:  item.OrderID.String(),
				Product:  prod,
				Quantity: item.Quantity,
				Price:    item.Price,
			})
		}
		enrichedOrders = append(enrichedOrders, enrichedOrder)
	}
	return enrichedOrders, nil
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

// GetUserCartWithProducts возвращает только заказ-корзину (pending) с товарами
func (s *OrderService) GetUserCartWithProducts(userID uuid.UUID) (*EnrichedOrder, error) {
	order, err := s.repo.GetPendingByUserID(userID)
	if err != nil || order == nil {
		return nil, err
	}
	enrichedOrder := &EnrichedOrder{
		ID:            order.ID.String(),
		UserID:        order.UserID.String(),
		Status:        order.Status,
		Total:         order.Total,
		ReservationID: order.ReservationID,
		CreatedAt:     order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     order.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	for _, item := range order.Items {
		log.Printf("[OrderService] Запрашиваю товар по id: %s", item.ProductID.String())
		prod, err := s.productsClient.GetProductByID(item.ProductID.String())
		if err != nil {
			log.Printf("[OrderService] Не удалось получить товар %s: %v", item.ProductID.String(), err)
			prod = &ProductInfo{ID: item.ProductID.String(), Name: "Товар не найден"}
		}
		enrichedOrder.Items = append(enrichedOrder.Items, EnrichedOrderItem{
			ID:       item.ID.String(),
			OrderID:  item.OrderID.String(),
			Product:  prod,
			Quantity: item.Quantity,
			Price:    item.Price,
		})
	}
	return enrichedOrder, nil
}

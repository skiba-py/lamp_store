package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/skiba/lamp_store/products_service/internal/domain"
)

type ReservationService struct {
	repo        domain.ReservationRepository
	productRepo domain.ProductRepository
}

func NewReservationService(repo domain.ReservationRepository, productRepo domain.ProductRepository) *ReservationService {
	return &ReservationService{
		repo:        repo,
		productRepo: productRepo,
	}
}

type ReserveRequest struct {
	Items []ReserveItem `json:"items"`
}

type ReserveItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func (s *ReservationService) ReserveProducts(ctx context.Context, req ReserveRequest) (*domain.Reservation, error) {
	// Создаем резервирование
	reservation := &domain.Reservation{
		ID:        uuid.New(),
		Status:    domain.ReservationStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(15 * time.Minute), // Резервирование на 15 минут
	}

	// Проверяем доступность товаров и создаем элементы резервирования
	reservation.Items = make([]domain.ReservationItem, len(req.Items))
	for i, item := range req.Items {
		productID, err := uuid.Parse(item.ProductID)
		if err != nil {
			return nil, err
		}

		product, err := s.productRepo.GetByID(ctx, productID)
		if err != nil {
			return nil, domain.ErrProductNotFound
		}

		if product.Stock < item.Quantity {
			return nil, domain.ErrInsufficientQuantity
		}

		// Уменьшаем количество товара
		product.Stock -= item.Quantity
		if err := s.productRepo.Update(ctx, product); err != nil {
			return nil, err
		}

		// Создаем элемент резервирования
		reservation.Items[i] = domain.ReservationItem{
			ID:            uuid.New(),
			ReservationID: reservation.ID,
			ProductID:     productID,
			Quantity:      item.Quantity,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
	}

	// Сохраняем резервирование
	if err := s.repo.Create(ctx, reservation); err != nil {
		// В случае ошибки возвращаем товары
		for _, item := range reservation.Items {
			product, _ := s.productRepo.GetByID(ctx, item.ProductID)
			if product != nil {
				product.Stock += item.Quantity
				_ = s.productRepo.Update(ctx, product)
			}
		}
		return nil, err
	}

	return reservation, nil
}

func (s *ReservationService) ConfirmReservation(ctx context.Context, id uuid.UUID) error {
	reservation, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.ErrReservationNotFound
	}

	if reservation.Status != domain.ReservationStatusPending {
		return domain.ErrInvalidStatus
	}

	reservation.Status = domain.ReservationStatusConfirmed
	reservation.UpdatedAt = time.Now()

	return s.repo.Update(ctx, reservation)
}

func (s *ReservationService) CancelReservation(ctx context.Context, id uuid.UUID) error {
	reservation, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.ErrReservationNotFound
	}

	if reservation.Status != domain.ReservationStatusPending {
		return domain.ErrInvalidStatus
	}

	// Возвращаем товары
	for _, item := range reservation.Items {
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			continue
		}
		product.Stock += item.Quantity
		_ = s.productRepo.Update(ctx, product)
	}

	reservation.Status = domain.ReservationStatusCancelled
	reservation.UpdatedAt = time.Now()

	return s.repo.Update(ctx, reservation)
}

func (s *ReservationService) CleanupExpiredReservations(ctx context.Context) error {
	expired, err := s.repo.GetExpiredReservations(ctx)
	if err != nil {
		return err
	}

	for _, reservation := range expired {
		if reservation.Status == domain.ReservationStatusPending {
			// Возвращаем товары
			for _, item := range reservation.Items {
				product, err := s.productRepo.GetByID(ctx, item.ProductID)
				if err != nil {
					continue
				}
				product.Stock += item.Quantity
				_ = s.productRepo.Update(ctx, product)
			}

			reservation.Status = domain.ReservationStatusExpired
			reservation.UpdatedAt = time.Now()
			_ = s.repo.Update(ctx, reservation)
		}
	}

	return nil
}

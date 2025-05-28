package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/skiba/lamp_store/products_service/internal/domain"
	"github.com/skiba/lamp_store/products_service/internal/service"
)

type ReservationHandler struct {
	service *service.ReservationService
}

func NewReservationHandler(service *service.ReservationService) *ReservationHandler {
	return &ReservationHandler{service: service}
}

func (h *ReservationHandler) Register(r chi.Router) {
	r.Route("/api/products/reserve", func(r chi.Router) {
		r.Post("/", h.ReserveProducts)
		r.Post("/{id}/confirm", h.ConfirmReservation)
		r.Post("/{id}/cancel", h.CancelReservation)
	})
}

func (h *ReservationHandler) ReserveProducts(w http.ResponseWriter, r *http.Request) {
	var req service.ReserveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	reservation, err := h.service.ReserveProducts(r.Context(), req)
	if err != nil {
		switch err {
		case domain.ErrProductNotFound:
			http.Error(w, "товар не найден", http.StatusNotFound)
		case domain.ErrInsufficientQuantity:
			http.Error(w, "недостаточное количество товара", http.StatusBadRequest)
		default:
			http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"reservation_id": reservation.ID.String()})
}

func (h *ReservationHandler) ConfirmReservation(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "неверный формат ID", http.StatusBadRequest)
		return
	}

	err = h.service.ConfirmReservation(r.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrReservationNotFound:
			http.Error(w, "резервирование не найдено", http.StatusNotFound)
		case domain.ErrInvalidStatus:
			http.Error(w, "недопустимый статус резервирования", http.StatusBadRequest)
		default:
			http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ReservationHandler) CancelReservation(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "неверный формат ID", http.StatusBadRequest)
		return
	}

	err = h.service.CancelReservation(r.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrReservationNotFound:
			http.Error(w, "резервирование не найдено", http.StatusNotFound)
		case domain.ErrInvalidStatus:
			http.Error(w, "недопустимый статус резервирования", http.StatusBadRequest)
		default:
			http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

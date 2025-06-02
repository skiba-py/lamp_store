package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/skiba/lamp_store/orders_service/internal/domain"
	"github.com/skiba/lamp_store/orders_service/internal/service"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) Register(r chi.Router) {
	r.Route("/api/orders", func(r chi.Router) {
		r.Post("/", h.CreateOrder)
		r.Get("/", h.GetAllOrders)
		r.Get("/{id}", h.GetOrder)
		r.Get("/user/{userID}", h.GetUserOrders)
		r.Put("/{id}", h.UpdateOrder)
		r.Delete("/{id}", h.DeleteOrder)
	})
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	type orderItemRaw struct {
		ID        string  `json:"id"`
		OrderID   string  `json:"order_id"`
		ProductID string  `json:"product_id"`
		Quantity  int     `json:"quantity"`
		Price     float64 `json:"price"`
	}
	type orderRaw struct {
		ID            string         `json:"id"`
		UserID        string         `json:"user_id"`
		Status        string         `json:"status"`
		Total         float64        `json:"total"`
		Items         []orderItemRaw `json:"items"`
		ReservationID string         `json:"reservation_id"`
		CreatedAt     string         `json:"created_at"`
		UpdatedAt     string         `json:"updated_at"`
	}

	var raw orderRaw
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("RAW ORDER: %+v", raw)

	order := domain.Order{
		Status:        raw.Status,
		Total:         raw.Total,
		ReservationID: raw.ReservationID,
	}
	if raw.ID != "" {
		if id, err := uuid.Parse(raw.ID); err == nil {
			order.ID = id
		}
	}
	if raw.UserID != "" {
		if uid, err := uuid.Parse(raw.UserID); err == nil {
			order.UserID = uid
		}
	}
	// created_at, updated_at можно добавить при необходимости

	for _, item := range raw.Items {
		var itemUUID, orderUUID, productUUID uuid.UUID
		var err error
		if item.ID != "" {
			itemUUID, _ = uuid.Parse(item.ID)
		}
		if item.OrderID != "" {
			orderUUID, _ = uuid.Parse(item.OrderID)
		}
		if item.ProductID != "" {
			productUUID, err = uuid.Parse(item.ProductID)
			if err != nil {
				http.Error(w, "Некорректный product_id: "+item.ProductID, http.StatusBadRequest)
				return
			}
		} else {
			http.Error(w, "product_id обязателен", http.StatusBadRequest)
			return
		}
		order.Items = append(order.Items, domain.OrderItem{
			ID:        itemUUID,
			OrderID:   orderUUID,
			ProductID: productUUID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	log.Printf("ORDER ITEMS BEFORE SAVE: %+v", order.Items)
	if err := h.orderService.CreateOrder(&order); err != nil {
		log.Printf("CreateOrder error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	order, err := h.orderService.GetOrder(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if order == nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	cart, err := h.orderService.GetUserCartWithProducts(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if cart == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(nil)
		return
	}

	json.NewEncoder(w).Encode(cart)
}

func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var order domain.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order.ID = id
	if err := h.orderService.UpdateOrder(&order); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	if err := h.orderService.DeleteOrder(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.orderService.GetAllOrders()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
}

package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/skiba/lamp_store/products_service/internal/domain"
)

type ProductHandler struct {
	service domain.ProductService
}

func NewProductHandler(service domain.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product domain.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateProduct(r.Context(), &product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := h.service.GetProduct(r.Context(), id)
	if err != nil {
		if err == domain.ErrProductNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var product domain.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	product.ID = id
	if err := h.service.UpdateProduct(r.Context(), &product); err != nil {
		if err == domain.ErrProductNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteProduct(r.Context(), id); err != nil {
		if err == domain.ErrProductNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	offset := 0
	limit := 10
	products, err := h.service.ListProducts(r.Context(), offset, limit)
	if err != nil {
		log.Printf("Ошибка при получении списка продуктов: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	type reqBody struct {
		Quantity int `json:"quantity"`
	}
	var req reqBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	product, err := h.service.GetProduct(r.Context(), id)
	if err != nil {
		if err == domain.ErrProductNotFound {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if product.Stock < req.Quantity {
		http.Error(w, "Insufficient stock", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "available"})
}

func (h *ProductHandler) Register(r chi.Router) {
	r.Route("/api/products", func(r chi.Router) {
		r.Post("/", h.CreateProduct)
		r.Post("/{id}/image", h.UploadImage)
		r.Get("/{id}", h.GetProduct)
		r.Put("/{id}", h.UpdateProduct)
		r.Delete("/{id}", h.DeleteProduct)
		r.Get("/", h.ListProducts)
	})
}

func (h *ProductHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Не удалось получить файл", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Создаём директорию для картинок, если не существует
	imgDir := "./images"
	if err := os.MkdirAll(imgDir, 0755); err != nil {
		http.Error(w, "Не удалось создать директорию для картинок", http.StatusInternalServerError)
		return
	}

	filename := id.String() + filepath.Ext(header.Filename)
	imgPath := filepath.Join(imgDir, filename)
	out, err := os.Create(imgPath)
	if err != nil {
		http.Error(w, "Не удалось сохранить файл", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		http.Error(w, "Ошибка при сохранении файла", http.StatusInternalServerError)
		return
	}

	// Сохраняем путь к картинке в БД
	product, err := h.service.GetProduct(r.Context(), id)
	if err != nil {
		http.Error(w, "Товар не найден", http.StatusNotFound)
		return
	}
	product.Image = "/images/" + filename
	if err := h.service.UpdateProduct(r.Context(), product); err != nil {
		http.Error(w, "Не удалось обновить товар", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"image": product.Image})
}

// UploadProductImage обрабатывает загрузку картинки для товара
func (h *ProductHandler) UploadProductImage(w http.ResponseWriter, r *http.Request) {
	// Ограничиваем размер файла до 5MB
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		log.Printf("Ошибка при разборе multipart form: %v", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		log.Printf("Ошибка при получении файла: %v", err)
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Создаем директорию, если её нет
	uploadDir := "/app/static/images"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("Ошибка при создании директории: %v", err)
		http.Error(w, "Error creating upload directory", http.StatusInternalServerError)
		return
	}

	// Создаем файл в директории
	dst, err := os.Create(filepath.Join(uploadDir, handler.Filename))
	if err != nil {
		log.Printf("Ошибка при создании файла: %v", err)
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Копируем содержимое загруженного файла
	if _, err := io.Copy(dst, file); err != nil {
		log.Printf("Ошибка при копировании файла: %v", err)
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "File uploaded successfully",
		"path":    filepath.Join("/static/images", handler.Filename),
	})
}

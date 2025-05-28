package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ProductsClient struct {
	baseURL string
	client  *http.Client
}

func NewProductsClient(baseURL string) *ProductsClient {
	return &ProductsClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

type ReserveProductsRequest struct {
	Items []ReserveItem `json:"items"`
}

type ReserveItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type ReserveProductsResponse struct {
	ReservationID string `json:"reservation_id"`
}

// ProductInfo описывает минимальную информацию о товаре для корзины/заказа
type ProductInfo struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Image       string  `json:"image,omitempty"`
	Description string  `json:"description,omitempty"`
}

func (c *ProductsClient) ReserveProducts(items []ReserveItem) (string, error) {
	reqBody := ReserveProductsRequest{Items: items}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("ошибка при маршалинге запроса: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/products/reserve", c.baseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ошибка при резервировании товаров: %d", resp.StatusCode)
	}

	var response ReserveProductsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return response.ReservationID, nil
}

func (c *ProductsClient) ConfirmReservation(reservationID string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/products/reserve/%s/confirm", c.baseURL, reservationID), nil)
	if err != nil {
		return fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка при подтверждении резервирования: %d", resp.StatusCode)
	}

	return nil
}

func (c *ProductsClient) CancelReservation(reservationID string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/products/reserve/%s/cancel", c.baseURL, reservationID), nil)
	if err != nil {
		return fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка при отмене резервирования: %d", resp.StatusCode)
	}

	return nil
}

// GetProductByID получает информацию о товаре по id
func (c *ProductsClient) GetProductByID(productID string) (*ProductInfo, error) {
	url := fmt.Sprintf("%s/api/products/%s", c.baseURL, productID)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе товара: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("товар не найден: %d", resp.StatusCode)
	}

	var product ProductInfo
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}
	return &product, nil
}

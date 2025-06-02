package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/skiba/lamp_store/admin_service/internal/domain"
)

type OrdersClient struct {
	baseURL string
	client  *http.Client
}

func NewOrdersClient(baseURL string) *OrdersClient {
	return &OrdersClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *OrdersClient) GetOrders() ([]domain.Order, error) {
	resp, err := c.client.Get(fmt.Sprintf("%s/api/orders", c.baseURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var orders []domain.Order
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (c *OrdersClient) UpdateOrder(id string, order *domain.Order) error {
	body, err := json.Marshal(order)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/api/orders/%s", c.baseURL, id),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

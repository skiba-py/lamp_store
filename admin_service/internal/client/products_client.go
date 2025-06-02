package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/skiba/lamp_store/admin_service/internal/domain"
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

func (c *ProductsClient) GetProducts() ([]domain.Product, error) {
	resp, err := c.client.Get(fmt.Sprintf("%s/api/products", c.baseURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var products []domain.Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, err
	}

	return products, nil
}

func (c *ProductsClient) CreateProduct(product *domain.Product) error {
	// Генерируем UUID для нового товара
	if product.ID == "" {
		product.ID = uuid.New().String()
	}

	body, err := json.Marshal(product)
	if err != nil {
		return err
	}

	resp, err := c.client.Post(
		fmt.Sprintf("%s/api/products", c.baseURL),
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *ProductsClient) UpdateProduct(id string, product *domain.Product) error {
	// Синхронизируем поля Stock и Quantity
	product.Quantity = product.Stock

	body, err := json.Marshal(product)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/api/products/%s", c.baseURL, id),
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *ProductsClient) DeleteProduct(id string) error {
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/api/products/%s", c.baseURL, id),
		nil,
	)
	if err != nil {
		return err
	}

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

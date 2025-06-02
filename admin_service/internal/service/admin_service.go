package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/skiba/lamp_store/admin_service/internal/client"
	"github.com/skiba/lamp_store/admin_service/internal/domain"
)

type AdminService struct {
	productsClient *client.ProductsClient
	ordersClient   *client.OrdersClient
}

func NewAdminService(productsClient *client.ProductsClient, ordersClient *client.OrdersClient) *AdminService {
	return &AdminService{
		productsClient: productsClient,
		ordersClient:   ordersClient,
	}
}

func (s *AdminService) Login(req *domain.LoginRequest) (*domain.LoginResponse, error) {
	// В продакшене проверять пароль из БД
	if req.Username != "admin" || req.Password != "admin" {
		return nil, errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "admin",
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{Token: tokenString}, nil
}

func (s *AdminService) GetProducts() ([]domain.Product, error) {
	return s.productsClient.GetProducts()
}

func (s *AdminService) CreateProduct(product *domain.Product) error {
	return s.productsClient.CreateProduct(product)
}

func (s *AdminService) UpdateProduct(id string, product *domain.Product) error {
	return s.productsClient.UpdateProduct(id, product)
}

func (s *AdminService) DeleteProduct(id string) error {
	return s.productsClient.DeleteProduct(id)
}

func (s *AdminService) GetOrders() ([]domain.Order, error) {
	return s.ordersClient.GetOrders()
}

func (s *AdminService) UpdateOrder(id string, order *domain.Order) error {
	return s.ordersClient.UpdateOrder(id, order)
}

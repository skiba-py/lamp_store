package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/skiba/lamp_store/admin_service/internal/client"
	"github.com/skiba/lamp_store/admin_service/internal/handler"
	"github.com/skiba/lamp_store/admin_service/internal/middleware"
	"github.com/skiba/lamp_store/admin_service/internal/service"
)

func main() {
	// Инициализация клиентов
	productsClient := client.NewProductsClient(os.Getenv("PRODUCTS_SERVICE_URL"))
	ordersClient := client.NewOrdersClient(os.Getenv("ORDERS_SERVICE_URL"))

	// Инициализация сервиса
	adminService := service.NewAdminService(productsClient, ordersClient)

	// Инициализация обработчиков
	adminHandler := handler.NewAdminHandler(adminService)

	// Инициализация роутера
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "1728000")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Публичные эндпоинты
	r.POST("/api/admin/login", adminHandler.Login)

	// Защищённые эндпоинты
	admin := r.Group("/api/admin")
	admin.Use(middleware.JWTAuthMiddleware())
	{
		// Продукты
		admin.GET("/products", adminHandler.GetProducts)
		admin.POST("/products", adminHandler.CreateProduct)
		admin.PUT("/products/:id", adminHandler.UpdateProduct)
		admin.DELETE("/products/:id", adminHandler.DeleteProduct)

		// Заказы
		admin.GET("/orders", adminHandler.GetOrders)
		admin.PUT("/orders/:id", adminHandler.UpdateOrder)
	}

	// Запуск сервера
	r.Run(":8003")
}

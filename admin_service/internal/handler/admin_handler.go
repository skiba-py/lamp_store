package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/skiba/lamp_store/admin_service/internal/domain"
	"github.com/skiba/lamp_store/admin_service/internal/service"
)

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

func (h *AdminHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.adminService.Login(&req)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, resp)
}

func (h *AdminHandler) GetProducts(c *gin.Context) {
	products, err := h.adminService.GetProducts()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, products)
}

func (h *AdminHandler) CreateProduct(c *gin.Context) {
	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.adminService.CreateProduct(&product); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, product)
}

func (h *AdminHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.adminService.UpdateProduct(id, &product); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, product)
}

func (h *AdminHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := h.adminService.DeleteProduct(id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Product deleted"})
}

func (h *AdminHandler) GetOrders(c *gin.Context) {
	orders, err := h.adminService.GetOrders()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, orders)
}

func (h *AdminHandler) UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	var order domain.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.adminService.UpdateOrder(id, &order); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, order)
}

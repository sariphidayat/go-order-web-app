package handlers

import (
	"fmt"
	"net/http"
	"time"

	"customer-order-app/backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderHandler struct {
	db *gorm.DB
}

type orderRequest struct {
	OrderNumber string  `json:"order_number"`
	CustomerID  uint    `json:"customer_id" binding:"required"`
	Status      string  `json:"status" binding:"required"`
	TotalAmount float64 `json:"total_amount"`
	Notes       string  `json:"notes"`
}

func NewOrderHandler(db *gorm.DB) OrderHandler {
	return OrderHandler{db: db}
}

func (h OrderHandler) List(c *gin.Context) {
	var orders []models.Order
	if err := h.db.Preload("Customer").Preload("Tracks").Order("created_at desc").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h OrderHandler) Get(c *gin.Context) {
	var order models.Order
	if err := h.db.Preload("Customer").Preload("Tracks", func(db *gorm.DB) *gorm.DB {
		return db.Order("tracked_at desc")
	}).First(&order, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h OrderHandler) Create(c *gin.Context) {
	var input orderRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderNumber := input.OrderNumber
	if orderNumber == "" {
		orderNumber = fmt.Sprintf("ORD-%d", time.Now().UnixNano())
	}

	order := models.Order{
		OrderNumber: orderNumber,
		CustomerID:  input.CustomerID,
		Status:      input.Status,
		TotalAmount: input.TotalAmount,
		Notes:       input.Notes,
	}
	if err := h.db.Create(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.db.Preload("Customer").First(&order, order.ID)
	c.JSON(http.StatusCreated, order)
}

func (h OrderHandler) Update(c *gin.Context) {
	var order models.Order
	if err := h.db.First(&order, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	var input orderRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.OrderNumber != "" {
		order.OrderNumber = input.OrderNumber
	}
	order.CustomerID = input.CustomerID
	order.Status = input.Status
	order.TotalAmount = input.TotalAmount
	order.Notes = input.Notes

	if err := h.db.Save(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.db.Preload("Customer").First(&order, order.ID)
	c.JSON(http.StatusOK, order)
}

func (h OrderHandler) Delete(c *gin.Context) {
	if err := h.db.Delete(&models.Order{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

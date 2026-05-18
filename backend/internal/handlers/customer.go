package handlers

import (
	"net/http"

	"customer-order-app/backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CustomerHandler struct {
	db *gorm.DB
}

type customerRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

func NewCustomerHandler(db *gorm.DB) CustomerHandler {
	return CustomerHandler{db: db}
}

func (h CustomerHandler) List(c *gin.Context) {
	var customers []models.Customer
	if err := h.db.Order("created_at desc").Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load customers"})
		return
	}
	c.JSON(http.StatusOK, customers)
}

func (h CustomerHandler) Get(c *gin.Context) {
	var customer models.Customer
	if err := h.db.Preload("Orders").First(&customer, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}
	c.JSON(http.StatusOK, customer)
}

func (h CustomerHandler) Create(c *gin.Context) {
	var input customerRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer := models.Customer{
		Name:    input.Name,
		Email:   input.Email,
		Phone:   input.Phone,
		Address: input.Address,
	}
	if err := h.db.Create(&customer).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, customer)
}

func (h CustomerHandler) Update(c *gin.Context) {
	var customer models.Customer
	if err := h.db.First(&customer, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}

	var input customerRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer.Name = input.Name
	customer.Email = input.Email
	customer.Phone = input.Phone
	customer.Address = input.Address
	if err := h.db.Save(&customer).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customer)
}

func (h CustomerHandler) Delete(c *gin.Context) {
	if err := h.db.Delete(&models.Customer{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

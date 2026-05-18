package handlers

import (
	"net/http"
	"time"

	"customer-order-app/backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderTrackHandler struct {
	db *gorm.DB
}

type orderTrackRequest struct {
	OrderID     uint      `json:"order_id" binding:"required"`
	Status      string    `json:"status" binding:"required"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
	TrackedAt   time.Time `json:"tracked_at"`
}

func NewOrderTrackHandler(db *gorm.DB) OrderTrackHandler {
	return OrderTrackHandler{db: db}
}

func (h OrderTrackHandler) List(c *gin.Context) {
	var tracks []models.OrderTrack
	query := h.db.Preload("Order").Order("tracked_at desc")
	if orderID := c.Query("order_id"); orderID != "" {
		query = query.Where("order_id = ?", orderID)
	}
	if err := query.Find(&tracks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load order tracks"})
		return
	}
	c.JSON(http.StatusOK, tracks)
}

func (h OrderTrackHandler) Get(c *gin.Context) {
	var track models.OrderTrack
	if err := h.db.Preload("Order").First(&track, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order track not found"})
		return
	}
	c.JSON(http.StatusOK, track)
}

func (h OrderTrackHandler) Create(c *gin.Context) {
	var input orderTrackRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trackedAt := input.TrackedAt
	if trackedAt.IsZero() {
		trackedAt = time.Now()
	}

	track := models.OrderTrack{
		OrderID:     input.OrderID,
		Status:      input.Status,
		Location:    input.Location,
		Description: input.Description,
		TrackedAt:   trackedAt,
	}
	if err := h.db.Create(&track).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.db.Model(&models.Order{}).Where("id = ?", input.OrderID).Update("status", input.Status)
	c.JSON(http.StatusCreated, track)
}

func (h OrderTrackHandler) Update(c *gin.Context) {
	var track models.OrderTrack
	if err := h.db.First(&track, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order track not found"})
		return
	}

	var input orderTrackRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trackedAt := input.TrackedAt
	if trackedAt.IsZero() {
		trackedAt = track.TrackedAt
	}

	track.OrderID = input.OrderID
	track.Status = input.Status
	track.Location = input.Location
	track.Description = input.Description
	track.TrackedAt = trackedAt
	if err := h.db.Save(&track).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, track)
}

func (h OrderTrackHandler) Delete(c *gin.Context) {
	if err := h.db.Delete(&models.OrderTrack{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

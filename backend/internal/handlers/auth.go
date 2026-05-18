package handlers

import (
	"fmt"
	"net/http"
	"time"

	"customer-order-app/backend/internal/config"
	"customer-order-app/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db  *gorm.DB
	cfg config.Config
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	OTP      string `json:"otp" binding:"required,len=6"`
}

func NewAuthHandler(db *gorm.DB, cfg config.Config) AuthHandler {
	return AuthHandler{db: db, cfg: cfg}
}

func (h AuthHandler) Login(c *gin.Context) {
	var input loginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if !totp.Validate(input.OTP, user.TOTPSecret) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authenticator code"})
		return
	}

	token, expiresAt, err := h.createToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      token,
		"expires_at": expiresAt,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func (h AuthHandler) Me(c *gin.Context) {
	userID := fmt.Sprintf("%v", c.GetString("user_id"))
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": user.ID, "name": user.Name, "email": user.Email})
}

func (h AuthHandler) createToken(user models.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(h.cfg.JWTExpiresMinutes) * time.Minute)
	claims := jwt.MapClaims{
		"sub":   fmt.Sprintf("%d", user.ID),
		"email": user.Email,
		"exp":   expiresAt.Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(h.cfg.JWTSecret))
	return signed, expiresAt, err
}

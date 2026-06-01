package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"customer-order-app/backend/internal/config"
	"customer-order-app/backend/internal/models"
	"customer-order-app/backend/internal/web"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const authCookieName = "app_token"

type WebHandler struct {
	db       *gorm.DB
	cfg      config.Config
	renderer *web.Renderer
}

type PageData struct {
	Title     string
	Active    string
	Error     string
	User      models.User
	Stats     map[string]int64
	Customers []models.Customer
	Orders    []models.Order
	Tracks    []models.OrderTrack
	Customer  models.Customer
	Order     models.Order
	Track     models.OrderTrack
	Statuses  []string
}

func NewWebHandler(db *gorm.DB, cfg config.Config, renderer *web.Renderer) WebHandler {
	return WebHandler{db: db, cfg: cfg, renderer: renderer}
}

func (h WebHandler) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie(authCookieName)
		if err != nil || tokenString == "" {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		claims, err := h.parseToken(tokenString)
		if err != nil {
			h.clearAuthCookie(c)
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		c.Set("user_id", fmt.Sprint(claims["sub"]))
		c.Next()
	}
}

func (h WebHandler) Home(c *gin.Context) {
	if _, err := c.Cookie(authCookieName); err == nil {
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}
	c.Redirect(http.StatusSeeOther, "/login")
}

func (h WebHandler) LoginPage(c *gin.Context) {
	h.renderer.Render(c.Writer, http.StatusOK, "login", PageData{Title: "Login"})
}

func (h WebHandler) Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	otp := c.PostForm("otp")

	var user models.User
	if err := h.db.Where("email = ?", email).First(&user).Error; err != nil {
		h.renderLoginError(c, "Email, password, atau kode authenticator tidak valid.")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		h.renderLoginError(c, "Email, password, atau kode authenticator tidak valid.")
		return
	}
	if !totp.Validate(otp, user.TOTPSecret) {
		h.renderLoginError(c, "Kode authenticator tidak valid.")
		return
	}

	token, expiresAt, err := h.createToken(user)
	if err != nil {
		h.renderLoginError(c, "Token login gagal dibuat.")
		return
	}

	maxAge := int(time.Until(expiresAt).Seconds())
	c.SetCookie(authCookieName, token, maxAge, "/", "", false, true)
	c.Header("HX-Redirect", "/dashboard")
	c.Redirect(http.StatusSeeOther, "/dashboard")
}

func (h WebHandler) Logout(c *gin.Context) {
	h.clearAuthCookie(c)
	c.Redirect(http.StatusSeeOther, "/login")
}

func (h WebHandler) Dashboard(c *gin.Context) {
	var customers int64
	var orders int64
	var tracks int64
	h.db.Model(&models.Customer{}).Count(&customers)
	h.db.Model(&models.Order{}).Count(&orders)
	h.db.Model(&models.OrderTrack{}).Count(&tracks)

	var latest []models.Order
	h.db.Preload("Customer").Order("created_at desc").Limit(8).Find(&latest)

	h.renderer.Render(c.Writer, http.StatusOK, "dashboard", h.pageData(c, "Dashboard", "dashboard", PageData{
		Stats: map[string]int64{
			"customers": customers,
			"orders":    orders,
			"tracks":    tracks,
		},
		Orders: latest,
	}))
}

func (h WebHandler) Customers(c *gin.Context) {
	data := h.customerData(c, PageData{Title: "Customers", Active: "customers"})
	h.renderer.Render(c.Writer, http.StatusOK, "customers", data)
}

func (h WebHandler) CreateCustomer(c *gin.Context) {
	customer := models.Customer{
		Name:    c.PostForm("name"),
		Email:   c.PostForm("email"),
		Phone:   c.PostForm("phone"),
		Address: c.PostForm("address"),
	}
	if err := h.db.Create(&customer).Error; err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	h.renderCustomersTable(c)
}

func (h WebHandler) EditCustomer(c *gin.Context) {
	var customer models.Customer
	if err := h.db.First(&customer, c.Param("id")).Error; err != nil {
		c.String(http.StatusNotFound, "customer not found")
		return
	}
	h.renderer.RenderPartial(c.Writer, http.StatusOK, "customers", "customer-edit-form", PageData{Customer: customer})
}

func (h WebHandler) UpdateCustomer(c *gin.Context) {
	var customer models.Customer
	if err := h.db.First(&customer, c.Param("id")).Error; err != nil {
		c.String(http.StatusNotFound, "customer not found")
		return
	}
	customer.Name = c.PostForm("name")
	customer.Email = c.PostForm("email")
	customer.Phone = c.PostForm("phone")
	customer.Address = c.PostForm("address")
	if err := h.db.Save(&customer).Error; err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	h.renderCustomersTable(c)
}

func (h WebHandler) DeleteCustomer(c *gin.Context) {
	if err := h.db.Delete(&models.Customer{}, c.Param("id")).Error; err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	h.renderCustomersTable(c)
}

func (h WebHandler) Orders(c *gin.Context) {
	data := h.orderData(c, PageData{Title: "Orders", Active: "orders"})
	h.renderer.Render(c.Writer, http.StatusOK, "orders", data)
}

func (h WebHandler) CreateOrder(c *gin.Context) {
	orderNumber := c.PostForm("order_number")
	if orderNumber == "" {
		orderNumber = fmt.Sprintf("ORD-%d", time.Now().UnixNano())
	}
	order := models.Order{
		OrderNumber: orderNumber,
		CustomerID:  parseUint(c.PostForm("customer_id")),
		Status:      c.PostForm("status"),
		TotalAmount: parseFloat(c.PostForm("total_amount")),
		Notes:       c.PostForm("notes"),
	}
	if err := h.db.Create(&order).Error; err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	h.renderOrdersTable(c)
}

func (h WebHandler) EditOrder(c *gin.Context) {
	var order models.Order
	if err := h.db.First(&order, c.Param("id")).Error; err != nil {
		c.String(http.StatusNotFound, "order not found")
		return
	}
	h.renderer.RenderPartial(c.Writer, http.StatusOK, "orders", "order-edit-form", h.orderData(c, PageData{Order: order}))
}

func (h WebHandler) UpdateOrder(c *gin.Context) {
	var order models.Order
	if err := h.db.First(&order, c.Param("id")).Error; err != nil {
		c.String(http.StatusNotFound, "order not found")
		return
	}
	if number := c.PostForm("order_number"); number != "" {
		order.OrderNumber = number
	}
	order.CustomerID = parseUint(c.PostForm("customer_id"))
	order.Status = c.PostForm("status")
	order.TotalAmount = parseFloat(c.PostForm("total_amount"))
	order.Notes = c.PostForm("notes")
	if err := h.db.Save(&order).Error; err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	h.renderOrdersTable(c)
}

func (h WebHandler) DeleteOrder(c *gin.Context) {
	if err := h.db.Delete(&models.Order{}, c.Param("id")).Error; err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	h.renderOrdersTable(c)
}

func (h WebHandler) Tracks(c *gin.Context) {
	data := h.trackData(c, PageData{Title: "Tracking", Active: "tracks"})
	h.renderer.Render(c.Writer, http.StatusOK, "tracks", data)
}

func (h WebHandler) CreateTrack(c *gin.Context) {
	track := models.OrderTrack{
		OrderID:     parseUint(c.PostForm("order_id")),
		Status:      c.PostForm("status"),
		Location:    c.PostForm("location"),
		Description: c.PostForm("description"),
		TrackedAt:   parseDateTime(c.PostForm("tracked_at")),
	}
	if track.TrackedAt.IsZero() {
		track.TrackedAt = time.Now()
	}
	if err := h.db.Create(&track).Error; err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	h.db.Model(&models.Order{}).Where("id = ?", track.OrderID).Update("status", track.Status)
	h.renderTracksTable(c)
}

func (h WebHandler) EditTrack(c *gin.Context) {
	var track models.OrderTrack
	if err := h.db.First(&track, c.Param("id")).Error; err != nil {
		c.String(http.StatusNotFound, "tracking record not found")
		return
	}
	h.renderer.RenderPartial(c.Writer, http.StatusOK, "tracks", "track-edit-form", h.trackData(c, PageData{Track: track}))
}

func (h WebHandler) UpdateTrack(c *gin.Context) {
	var track models.OrderTrack
	if err := h.db.First(&track, c.Param("id")).Error; err != nil {
		c.String(http.StatusNotFound, "tracking record not found")
		return
	}
	track.OrderID = parseUint(c.PostForm("order_id"))
	track.Status = c.PostForm("status")
	track.Location = c.PostForm("location")
	track.Description = c.PostForm("description")
	track.TrackedAt = parseDateTime(c.PostForm("tracked_at"))
	if track.TrackedAt.IsZero() {
		track.TrackedAt = time.Now()
	}
	if err := h.db.Save(&track).Error; err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	h.renderTracksTable(c)
}

func (h WebHandler) DeleteTrack(c *gin.Context) {
	if err := h.db.Delete(&models.OrderTrack{}, c.Param("id")).Error; err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	h.renderTracksTable(c)
}

func (h WebHandler) renderLoginError(c *gin.Context, message string) {
	h.renderer.Render(c.Writer, http.StatusUnauthorized, "login", PageData{Title: "Login", Error: message})
}

func (h WebHandler) renderCustomersTable(c *gin.Context) {
	h.renderer.RenderPartial(c.Writer, http.StatusOK, "customers", "customers-table", h.customerData(c, PageData{}))
}

func (h WebHandler) renderOrdersTable(c *gin.Context) {
	h.renderer.RenderPartial(c.Writer, http.StatusOK, "orders", "orders-table", h.orderData(c, PageData{}))
}

func (h WebHandler) renderTracksTable(c *gin.Context) {
	h.renderer.RenderPartial(c.Writer, http.StatusOK, "tracks", "tracks-table", h.trackData(c, PageData{}))
}

func (h WebHandler) customerData(c *gin.Context, data PageData) PageData {
	h.db.Order("created_at desc").Find(&data.Customers)
	return h.pageData(c, data.Title, data.Active, data)
}

func (h WebHandler) orderData(c *gin.Context, data PageData) PageData {
	h.db.Order("name asc").Find(&data.Customers)
	h.db.Preload("Customer").Order("created_at desc").Find(&data.Orders)
	data.Statuses = []string{"pending", "processing", "shipped", "delivered", "cancelled"}
	return h.pageData(c, data.Title, data.Active, data)
}

func (h WebHandler) trackData(c *gin.Context, data PageData) PageData {
	h.db.Order("created_at desc").Find(&data.Orders)
	h.db.Preload("Order").Order("tracked_at desc").Find(&data.Tracks)
	data.Statuses = []string{"pending", "processing", "packed", "shipped", "in_transit", "delivered", "cancelled"}
	return h.pageData(c, data.Title, data.Active, data)
}

func (h WebHandler) pageData(c *gin.Context, title string, active string, data PageData) PageData {
	if title != "" {
		data.Title = title
	}
	if active != "" {
		data.Active = active
	}
	userID := fmt.Sprint(c.GetString("user_id"))
	if userID != "" {
		h.db.First(&data.User, userID)
	}
	return data
}

func (h WebHandler) createToken(user models.User) (string, time.Time, error) {
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

func (h WebHandler) parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.cfg.JWTSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}

func (h WebHandler) clearAuthCookie(c *gin.Context) {
	c.SetCookie(authCookieName, "", -1, "/", "", false, true)
}

func parseUint(value string) uint {
	parsed, _ := strconv.ParseUint(value, 10, 64)
	return uint(parsed)
}

func parseFloat(value string) float64 {
	parsed, _ := strconv.ParseFloat(value, 64)
	return parsed
}

func parseDateTime(value string) time.Time {
	parsed, _ := time.ParseInLocation("2006-01-02T15:04", value, time.Local)
	return parsed
}

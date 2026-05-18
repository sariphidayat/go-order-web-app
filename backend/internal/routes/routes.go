package routes

import (
	"net/http"
	"strings"
	"time"

	"customer-order-app/backend/internal/config"
	"customer-order-app/backend/internal/handlers"
	"customer-order-app/backend/internal/middleware"
	"customer-order-app/backend/internal/web"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, cfg config.Config) *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontendOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	authHandler := handlers.NewAuthHandler(db, cfg)
	customerHandler := handlers.NewCustomerHandler(db)
	orderHandler := handlers.NewOrderHandler(db)
	trackHandler := handlers.NewOrderTrackHandler(db)

	api := router.Group("/api")
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	api.POST("/auth/login", authHandler.Login)

	protected := api.Group("")
	protected.Use(middleware.AuthRequired(cfg))
	protected.GET("/auth/me", authHandler.Me)

	protected.GET("/customers", customerHandler.List)
	protected.POST("/customers", customerHandler.Create)
	protected.GET("/customers/:id", customerHandler.Get)
	protected.PUT("/customers/:id", customerHandler.Update)
	protected.DELETE("/customers/:id", customerHandler.Delete)

	protected.GET("/orders", orderHandler.List)
	protected.POST("/orders", orderHandler.Create)
	protected.GET("/orders/:id", orderHandler.Get)
	protected.PUT("/orders/:id", orderHandler.Update)
	protected.DELETE("/orders/:id", orderHandler.Delete)

	protected.GET("/order-tracks", trackHandler.List)
	protected.POST("/order-tracks", trackHandler.Create)
	protected.GET("/order-tracks/:id", trackHandler.Get)
	protected.PUT("/order-tracks/:id", trackHandler.Update)
	protected.DELETE("/order-tracks/:id", trackHandler.Delete)

	registerFrontend(router)

	return router
}

func registerFrontend(router *gin.Engine) {
	dist, err := web.Dist()
	if err != nil {
		return
	}

	fileServer := http.FileServer(http.FS(dist))
	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
			return
		}

		path := strings.TrimPrefix(c.Request.URL.Path, "/")
		if path != "" {
			if file, err := dist.Open(path); err == nil {
				_ = file.Close()
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}
		}

		c.Request.URL.Path = "/"
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}

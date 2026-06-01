package routes

import (
	"log"
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
	renderer, err := web.NewRenderer()
	if err != nil {
		log.Fatalf("template initialization failed: %v", err)
	}
	webHandler := handlers.NewWebHandler(db, cfg, renderer)

	router.GET("/", webHandler.Home)
	router.GET("/login", webHandler.LoginPage)
	router.POST("/login", webHandler.Login)
	router.GET("/logout", webHandler.Logout)

	pages := router.Group("")
	pages.Use(webHandler.RequireAuth())
	pages.GET("/dashboard", webHandler.Dashboard)
	pages.GET("/customers", webHandler.Customers)
	pages.POST("/customers", webHandler.CreateCustomer)
	pages.GET("/customers/:id/edit", webHandler.EditCustomer)
	pages.POST("/customers/:id/update", webHandler.UpdateCustomer)
	pages.DELETE("/customers/:id", webHandler.DeleteCustomer)
	pages.GET("/orders", webHandler.Orders)
	pages.POST("/orders", webHandler.CreateOrder)
	pages.GET("/orders/:id/edit", webHandler.EditOrder)
	pages.POST("/orders/:id/update", webHandler.UpdateOrder)
	pages.DELETE("/orders/:id", webHandler.DeleteOrder)
	pages.GET("/tracks", webHandler.Tracks)
	pages.POST("/tracks", webHandler.CreateTrack)
	pages.GET("/tracks/:id/edit", webHandler.EditTrack)
	pages.POST("/tracks/:id/update", webHandler.UpdateTrack)
	pages.DELETE("/tracks/:id", webHandler.DeleteTrack)

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

	return router
}

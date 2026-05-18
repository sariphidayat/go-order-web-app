package database

import (
	"customer-order-app/backend/internal/config"
	"customer-order-app/backend/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(dsn string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{}, &models.Customer{}, &models.Order{}, &models.OrderTrack{})
}

func SeedAdmin(db *gorm.DB, cfg config.Config) error {
	var count int64
	if err := db.Model(&models.User{}).Where("email = ?", cfg.AdminEmail).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return db.Create(&models.User{
		Name:       "Administrator",
		Email:      cfg.AdminEmail,
		Password:   string(passwordHash),
		TOTPSecret: cfg.AdminTOTPSecret,
	}).Error
}

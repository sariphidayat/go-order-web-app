package models

import "time"

type User struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name" gorm:"size:120;not null"`
	Email      string    `json:"email" gorm:"size:160;uniqueIndex;not null"`
	Password   string    `json:"-" gorm:"not null"`
	TOTPSecret string    `json:"-" gorm:"size:64;not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Customer struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:140;not null"`
	Email     string    `json:"email" gorm:"size:160;uniqueIndex;not null"`
	Phone     string    `json:"phone" gorm:"size:40"`
	Address   string    `json:"address" gorm:"size:255"`
	Orders    []Order   `json:"orders,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Order struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	OrderNumber string       `json:"order_number" gorm:"size:60;uniqueIndex;not null"`
	CustomerID  uint         `json:"customer_id" gorm:"not null"`
	Customer    Customer     `json:"customer,omitempty"`
	Status      string       `json:"status" gorm:"size:40;not null;default:pending"`
	TotalAmount float64      `json:"total_amount" gorm:"not null;default:0"`
	Notes       string       `json:"notes" gorm:"size:255"`
	Tracks      []OrderTrack `json:"tracks,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type OrderTrack struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	OrderID     uint      `json:"order_id" gorm:"not null"`
	Order       Order     `json:"order,omitempty"`
	Status      string    `json:"status" gorm:"size:60;not null"`
	Location    string    `json:"location" gorm:"size:120"`
	Description string    `json:"description" gorm:"size:255"`
	TrackedAt   time.Time `json:"tracked_at" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

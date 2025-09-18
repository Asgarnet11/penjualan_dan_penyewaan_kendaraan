package model

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	VehicleID    uuid.UUID `json:"vehicle_id"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	TotalPrice   float64   `json:"total_price"`
	Status       string    `json:"status"`
	PaymentToken string    `json:"payment_token,omitempty"`
	PaymentURL   string    `json:"payment_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateBookingInput struct {
	VehicleID string `json:"vehicle_id" binding:"required"`
	StartDate string `json:"start_date" binding:"required"` // Format: "YYYY-MM-DD"
	EndDate   string `json:"end_date" binding:"required"`   // Format: "YYYY-MM-DD"
}

package model

import (
	"time"

	"github.com/google/uuid"
)

type SalesTransaction struct {
	ID           uuid.UUID `json:"id"`
	VehicleID    uuid.UUID `json:"vehicle_id"`
	SellerID     uuid.UUID `json:"seller_id"`
	BuyerID      uuid.UUID `json:"buyer_id"`
	AgreedPrice  float64   `json:"agreed_price"`
	Status       string    `json:"status"`
	PaymentToken string    `json:"payment_token,omitempty"`
	PaymentURL   string    `json:"payment_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

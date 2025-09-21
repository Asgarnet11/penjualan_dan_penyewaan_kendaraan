package model

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID         uuid.UUID `json:"id"`
	CustomerID uuid.UUID `json:"customer_id"`
	VendorID   uuid.UUID `json:"vendor_id"`
	VehicleID  uuid.UUID `json:"vehicle_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

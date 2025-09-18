package model

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID        uuid.UUID `json:"id"`
	BookingID uuid.UUID `json:"booking_id"`
	UserID    uuid.UUID `json:"user_id"`
	VehicleID uuid.UUID `json:"vehicle_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateReviewInput struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment"`
}

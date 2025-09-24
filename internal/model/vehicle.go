package model

import (
	"time"

	"github.com/google/uuid"
)

type Vehicle struct {
	ID                 uuid.UUID `json:"id"`
	OwnerID            uuid.UUID `json:"owner_id"`
	Brand              string    `json:"brand"`
	Model              string    `json:"model"`
	Year               int       `json:"year"`
	PlateNumber        string    `json:"plate_number"`
	Color              string    `json:"color"`
	VehicleType        string    `json:"vehicle_type"`
	Transmission       string    `json:"transmission"`
	Fuel               string    `json:"fuel"`
	Status             string    `json:"status"`
	Description        string    `json:"description"`
	IsForSale          bool      `json:"is_for_sale"`
	SalePrice          float64   `json:"sale_price"`
	IsForRent          bool      `json:"is_for_rent"`
	RentalPriceDaily   float64   `json:"rental_price_daily"`
	RentalPriceWeekly  float64   `json:"rental_price_weekly"`
	RentalPriceMonthly float64   `json:"rental_price_monthly"`
	Location           *string   `json:"location,omitempty"`
	Features           []string  `json:"features,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type CreateVehicleInput struct {
	Brand              string  `json:"brand" binding:"required"`
	Model              string  `json:"model" binding:"required"`
	Year               int     `json:"year" binding:"required"`
	PlateNumber        string  `json:"plate_number" binding:"required"`
	Color              string  `json:"color"`
	VehicleType        string  `json:"vehicle_type" binding:"required,oneof=mobil motor"`
	Transmission       string  `json:"transmission" binding:"required,oneof=matic manual"`
	Fuel               string  `json:"fuel" binding:"required,oneof=bensin diesel listrik"`
	Description        string  `json:"description"`
	IsForSale          bool    `json:"is_for_sale"`
	SalePrice          float64 `json:"sale_price"`
	IsForRent          bool    `json:"is_for_rent"`
	RentalPriceDaily   float64 `json:"rental_price_daily"`
	RentalPriceWeekly  float64 `json:"rental_price_weekly"`
	RentalPriceMonthly float64 `json:"rental_price_monthly"`
}

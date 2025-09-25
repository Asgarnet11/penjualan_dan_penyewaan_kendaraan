package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type VehicleImage struct {
	ID        uuid.UUID `json:"id"`
	ImageURL  string    `json:"image_url"`
	IsPrimary bool      `json:"is_primary"`
}
type Vehicle struct {
	ID                 uuid.UUID     `json:"id"`
	OwnerID            uuid.UUID     `json:"owner_id"`
	Brand              string        `json:"brand"`
	Model              string        `json:"model"`
	Year               int           `json:"year"`
	PlateNumber        string        `json:"plate_number"`
	Color              *string       `json:"color,omitempty"` // <-- Pointer
	VehicleType        string        `json:"vehicle_type"`
	Transmission       string        `json:"transmission"`
	Fuel               string        `json:"fuel"`
	Status             string        `json:"status"`
	Description        *string       `json:"description,omitempty"` // <-- Pointer
	IsForSale          bool          `json:"is_for_sale"`
	SalePrice          *float64      `json:"sale_price,omitempty"` // <-- Pointer
	IsForRent          bool          `json:"is_for_rent"`
	RentalPriceDaily   *float64      `json:"rental_price_daily,omitempty"`   // <-- Pointer
	RentalPriceWeekly  *float64      `json:"rental_price_weekly,omitempty"`  // <-- Pointer
	RentalPriceMonthly *float64      `json:"rental_price_monthly,omitempty"` // <-- Pointer
	Location           *string       `json:"location,omitempty"`
	Features           []string      `json:"features,omitempty"`
	Images             VehicleImages `json:"images"`
	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`
}

type CreateVehicleInput struct {
	Brand              string   `json:"brand" binding:"required"`
	Model              string   `json:"model" binding:"required"`
	Year               int      `json:"year" binding:"required"`
	PlateNumber        string   `json:"plate_number" binding:"required"`
	Color              string   `json:"color"`
	VehicleType        string   `json:"vehicle_type" binding:"required,oneof=mobil motor"`
	Transmission       string   `json:"transmission" binding:"required,oneof=matic manual"`
	Fuel               string   `json:"fuel" binding:"required,oneof=bensin diesel listrik"`
	Description        string   `json:"description"`
	IsForSale          bool     `json:"is_for_sale"`
	SalePrice          float64  `json:"sale_price"`
	IsForRent          bool     `json:"is_for_rent"`
	RentalPriceDaily   float64  `json:"rental_price_daily"`
	RentalPriceWeekly  float64  `json:"rental_price_weekly"`
	RentalPriceMonthly float64  `json:"rental_price_monthly"`
	Location           string   `json:"location"`
	Features           []string `json:"features"`
}

type VehicleImages []VehicleImage

func (vi *VehicleImages) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &vi)
}
func (vi VehicleImages) Value() (driver.Value, error) {
	return json.Marshal(vi)
}

package model

type VehicleFilter struct {
	Type         string `form:"type"`
	Brand        string `form:"brand"`
	Model        string `form:"model"`
	Transmission string `form:"transmission"`
	MinYear      int    `form:"min_year"`
	MaxPrice     int    `form:"max_price"`
	Search       string `form:"search"`
	Sort         string `form:"sort"`
	IsForSale    bool   `form:"is_for_sale"`
	IsForRent    bool   `form:"is_for_rent"`
}

package repository

import (
	"context"
	"fmt"
	"strings"
	"sultra-otomotif-api/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// VehicleRepository adalah interface yang akan digunakan oleh service
type VehicleRepository interface {
	Create(ctx context.Context, vehicle model.Vehicle) (model.Vehicle, error)
	FindAll(ctx context.Context, filter model.VehicleFilter) ([]model.Vehicle, error)
	FindByID(ctx context.Context, id uuid.UUID) (model.Vehicle, error)
	Update(ctx context.Context, vehicle model.Vehicle) (model.Vehicle, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// vehicleRepository adalah implementasi dari interface di atas
type vehicleRepository struct {
	db *pgxpool.Pool
}

// NewVehicleRepository adalah constructor untuk vehicleRepository
func NewVehicleRepository(db *pgxpool.Pool) VehicleRepository {
	return &vehicleRepository{db: db}
}

func (r *vehicleRepository) Create(ctx context.Context, v model.Vehicle) (model.Vehicle, error) {
	query := `INSERT INTO vehicles (id, owner_id, brand, model, year, plate_number, color, vehicle_type, transmission, fuel, status, description, is_for_sale, sale_price, is_for_rent, rental_price_daily, rental_price_weekly, rental_price_monthly)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
              RETURNING created_at, updated_at`

	err := r.db.QueryRow(ctx, query, v.ID, v.OwnerID, v.Brand, v.Model, v.Year, v.PlateNumber, v.Color, v.VehicleType, v.Transmission, v.Fuel, v.Status, v.Description, v.IsForSale, v.SalePrice, v.IsForRent, v.RentalPriceDaily, v.RentalPriceWeekly, v.RentalPriceMonthly).Scan(&v.CreatedAt, &v.UpdatedAt)

	if err != nil {
		return model.Vehicle{}, err
	}
	return v, nil
}

// FindAll telah dimodifikasi untuk membangun query secara dinamis berdasarkan filter
func (r *vehicleRepository) FindAll(ctx context.Context, filter model.VehicleFilter) ([]model.Vehicle, error) {
	var vehicles []model.Vehicle
	baseQuery := `SELECT id, owner_id, brand, model, year, plate_number, color, vehicle_type, transmission, fuel, status, description, is_for_sale, sale_price, is_for_rent, rental_price_daily, rental_price_weekly, rental_price_monthly, created_at, updated_at FROM vehicles WHERE status = 'available'`

	conditions := []string{}
	args := []interface{}{}
	argID := 1

	if filter.Type != "" {
		conditions = append(conditions, fmt.Sprintf("vehicle_type = $%d", argID))
		args = append(args, filter.Type)
		argID++
	}
	if filter.Brand != "" {
		conditions = append(conditions, fmt.Sprintf("brand ILIKE $%d", argID))
		args = append(args, "%"+filter.Brand+"%")
		argID++
	}
	if filter.Transmission != "" {
		conditions = append(conditions, fmt.Sprintf("transmission = $%d", argID))
		args = append(args, filter.Transmission)
		argID++
	}
	if filter.MinYear > 0 {
		conditions = append(conditions, fmt.Sprintf("year >= $%d", argID))
		args = append(args, filter.MinYear)
		argID++
	}
	if filter.MaxPrice > 0 {
		conditions = append(conditions, fmt.Sprintf("sale_price <= $%d", argID))
		args = append(args, filter.MaxPrice)
		argID++
	}
	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(brand ILIKE $%d OR model ILIKE $%d)", argID, argID))
		args = append(args, "%"+filter.Search+"%")
		argID++
	}

	// Gabungkan semua kondisi filter ke query utama
	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	// Tambahkan sorting
	orderBy := " ORDER BY created_at DESC" // Urutan default
	switch filter.Sort {
	case "price_asc":
		orderBy = " ORDER BY sale_price ASC"
	case "price_desc":
		orderBy = " ORDER BY sale_price DESC"
	case "year_desc":
		orderBy = " ORDER BY year DESC"
	}
	baseQuery += orderBy

	// Eksekusi query yang sudah dibangun secara dinamis
	rows, err := r.db.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var v model.Vehicle
		err := rows.Scan(&v.ID, &v.OwnerID, &v.Brand, &v.Model, &v.Year, &v.PlateNumber, &v.Color, &v.VehicleType, &v.Transmission, &v.Fuel, &v.Status, &v.Description, &v.IsForSale, &v.SalePrice, &v.IsForRent, &v.RentalPriceDaily, &v.RentalPriceWeekly, &v.RentalPriceMonthly, &v.CreatedAt, &v.UpdatedAt)
		if err != nil {
			return nil, err
		}
		vehicles = append(vehicles, v)
	}

	return vehicles, nil
}

func (r *vehicleRepository) FindByID(ctx context.Context, id uuid.UUID) (model.Vehicle, error) {
	var v model.Vehicle
	query := `SELECT id, owner_id, brand, model, year, plate_number, color, vehicle_type, transmission, fuel, status, description, is_for_sale, sale_price, is_for_rent, rental_price_daily, rental_price_weekly, rental_price_monthly, created_at, updated_at FROM vehicles WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).Scan(&v.ID, &v.OwnerID, &v.Brand, &v.Model, &v.Year, &v.PlateNumber, &v.Color, &v.VehicleType, &v.Transmission, &v.Fuel, &v.Status, &v.Description, &v.IsForSale, &v.SalePrice, &v.IsForRent, &v.RentalPriceDaily, &v.RentalPriceWeekly, &v.RentalPriceMonthly, &v.CreatedAt, &v.UpdatedAt)

	if err != nil {
		return model.Vehicle{}, err
	}
	return v, nil
}

func (r *vehicleRepository) Update(ctx context.Context, v model.Vehicle) (model.Vehicle, error) {
	query := `UPDATE vehicles SET brand=$1, model=$2, year=$3, plate_number=$4, color=$5, vehicle_type=$6, transmission=$7, fuel=$8, status=$9, description=$10, is_for_sale=$11, sale_price=$12, is_for_rent=$13, rental_price_daily=$14, rental_price_weekly=$15, rental_price_monthly=$16, updated_at=NOW()
              WHERE id=$17 RETURNING updated_at`

	err := r.db.QueryRow(ctx, query, v.Brand, v.Model, v.Year, v.PlateNumber, v.Color, v.VehicleType, v.Transmission, v.Fuel, v.Status, v.Description, v.IsForSale, v.SalePrice, v.IsForRent, v.RentalPriceDaily, v.RentalPriceWeekly, v.RentalPriceMonthly, v.ID).Scan(&v.UpdatedAt)

	if err != nil {
		return model.Vehicle{}, err
	}
	return v, nil
}

func (r *vehicleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM vehicles WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

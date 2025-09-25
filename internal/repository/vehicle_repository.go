package repository

import (
	"context"
	"fmt"
	"strings"
	"sultra-otomotif-api/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const vehicleWithImagesQuery = `
	SELECT 
		v.id, v.owner_id, v.brand, v.model, v.year, v.plate_number, v.color, 
		v.vehicle_type, v.transmission, v.fuel, v.status, v.description, 
		v.is_for_sale, v.sale_price, v.is_for_rent, v.rental_price_daily, 
		v.rental_price_weekly, v.rental_price_monthly, v.location, v.features,
		v.created_at, v.updated_at,
		COALESCE(
			(SELECT json_agg(json_build_object('id', vi.id, 'image_url', vi.image_url, 'is_primary', vi.is_primary))
			 FROM vehicle_images vi WHERE vi.vehicle_id = v.id),
			'[]'::json
		) AS images
	FROM vehicles v
`

func scanVehicle(row pgx.Row, v *model.Vehicle) error {
	return row.Scan(
		&v.ID, &v.OwnerID, &v.Brand, &v.Model, &v.Year, &v.PlateNumber, &v.Color,
		&v.VehicleType, &v.Transmission, &v.Fuel, &v.Status, &v.Description,
		&v.IsForSale, &v.SalePrice, &v.IsForRent, &v.RentalPriceDaily,
		&v.RentalPriceWeekly, &v.RentalPriceMonthly, &v.Location, &v.Features,
		&v.CreatedAt, &v.UpdatedAt, &v.Images,
	)
}

type VehicleRepository interface {
	Create(ctx context.Context, vehicle model.Vehicle) (model.Vehicle, error)
	FindAll(ctx context.Context, filter model.VehicleFilter) ([]model.Vehicle, error)
	FindByID(ctx context.Context, id uuid.UUID) (model.Vehicle, error)
	Update(ctx context.Context, vehicle model.Vehicle) (model.Vehicle, error)
	Delete(ctx context.Context, id uuid.UUID) error
	FindAllAdmin(ctx context.Context) ([]model.Vehicle, error)
	FindAllByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]model.Vehicle, error)
}

type vehicleRepository struct {
	db *pgxpool.Pool
}

func NewVehicleRepository(db *pgxpool.Pool) VehicleRepository {
	return &vehicleRepository{db: db}
}

func (r *vehicleRepository) Create(ctx context.Context, v model.Vehicle) (model.Vehicle, error) {
	query := `INSERT INTO vehicles (id, owner_id, brand, model, year, plate_number, color, vehicle_type, transmission, fuel, status, description, is_for_sale, sale_price, is_for_rent, rental_price_daily, rental_price_weekly, rental_price_monthly, location, features)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
              RETURNING created_at, updated_at`

	err := r.db.QueryRow(ctx, query, v.ID, v.OwnerID, v.Brand, v.Model, v.Year, v.PlateNumber, v.Color, v.VehicleType, v.Transmission, v.Fuel, v.Status, v.Description, v.IsForSale, v.SalePrice, v.IsForRent, v.RentalPriceDaily, v.RentalPriceWeekly, v.RentalPriceMonthly, v.Location, v.Features).Scan(&v.CreatedAt, &v.UpdatedAt)

	if err != nil {
		return model.Vehicle{}, err
	}
	return v, nil
}

func (r *vehicleRepository) FindAll(ctx context.Context, filter model.VehicleFilter) ([]model.Vehicle, error) {
	var vehicles []model.Vehicle

	conditions := []string{"v.status = 'available'"}
	args := []interface{}{}
	argID := 1

	if filter.Type != "" {
		conditions = append(conditions, fmt.Sprintf("v.vehicle_type = $%d", argID))
		args = append(args, filter.Type)
		argID++
	}
	if filter.Brand != "" {
		conditions = append(conditions, fmt.Sprintf("v.brand ILIKE $%d", argID))
		args = append(args, "%"+filter.Brand+"%")
		argID++
	}
	if filter.Transmission != "" {
		conditions = append(conditions, fmt.Sprintf("v.transmission = $%d", argID))
		args = append(args, filter.Transmission)
		argID++
	}
	// ... tambahkan logika filter lain di sini jika perlu ...

	finalQuery := vehicleWithImagesQuery
	if len(conditions) > 0 {
		finalQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	orderBy := " ORDER BY v.created_at DESC"
	switch filter.Sort {
	case "price_asc":
		orderBy = " ORDER BY v.sale_price ASC NULLS LAST"
	case "price_desc":
		orderBy = " ORDER BY v.sale_price DESC NULLS LAST"
	case "year_desc":
		orderBy = " ORDER BY v.year DESC"
	}
	finalQuery += orderBy

	rows, err := r.db.Query(ctx, finalQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var v model.Vehicle
		if err := scanVehicle(rows, &v); err != nil {
			return nil, err
		}
		vehicles = append(vehicles, v)
	}

	return vehicles, nil
}

func (r *vehicleRepository) FindByID(ctx context.Context, id uuid.UUID) (model.Vehicle, error) {
	var v model.Vehicle
	query := vehicleWithImagesQuery + " WHERE v.id = $1"

	row := r.db.QueryRow(ctx, query, id)
	if err := scanVehicle(row, &v); err != nil {
		return model.Vehicle{}, err
	}
	return v, nil
}

func (r *vehicleRepository) Update(ctx context.Context, v model.Vehicle) (model.Vehicle, error) {
	query := `UPDATE vehicles SET brand=$1, model=$2, year=$3, plate_number=$4, color=$5, vehicle_type=$6, transmission=$7, fuel=$8, status=$9, description=$10, is_for_sale=$11, sale_price=$12, is_for_rent=$13, rental_price_daily=$14, rental_price_weekly=$15, rental_price_monthly=$16, location=$17, features=$18, updated_at=NOW()
              WHERE id=$19 RETURNING updated_at`

	err := r.db.QueryRow(ctx, query, v.Brand, v.Model, v.Year, v.PlateNumber, v.Color, v.VehicleType, v.Transmission, v.Fuel, v.Status, v.Description, v.IsForSale, v.SalePrice, v.IsForRent, v.RentalPriceDaily, v.RentalPriceWeekly, v.RentalPriceMonthly, v.Location, v.Features, v.ID).Scan(&v.UpdatedAt)

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

func (r *vehicleRepository) FindAllAdmin(ctx context.Context) ([]model.Vehicle, error) {
	var vehicles []model.Vehicle
	query := vehicleWithImagesQuery + " ORDER BY v.created_at DESC"

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var v model.Vehicle
		if err := scanVehicle(rows, &v); err != nil {
			return nil, err
		}
		vehicles = append(vehicles, v)
	}
	return vehicles, nil
}

func (r *vehicleRepository) FindAllByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]model.Vehicle, error) {
	var vehicles []model.Vehicle
	query := vehicleWithImagesQuery + " WHERE v.owner_id = $1 ORDER BY v.created_at DESC"

	rows, err := r.db.Query(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var v model.Vehicle
		if err := scanVehicle(rows, &v); err != nil {
			return nil, err
		}
		vehicles = append(vehicles, v)
	}

	return vehicles, nil
}

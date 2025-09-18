package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ImageRepository interface {
	SaveVehicleImage(ctx context.Context, vehicleID uuid.UUID, imageURL string) error
}

type imageRepository struct {
	db *pgxpool.Pool
}

func NewImageRepository(db *pgxpool.Pool) ImageRepository {
	return &imageRepository{db: db}
}

func (r *imageRepository) SaveVehicleImage(ctx context.Context, vehicleID uuid.UUID, imageURL string) error {
	query := `INSERT INTO vehicle_images (id, vehicle_id, image_url, is_primary) VALUES (uuid_generate_v4(), $1, $2, false)` // is_primary akan kita handle nanti
	_, err := r.db.Exec(ctx, query, vehicleID, imageURL)
	return err
}

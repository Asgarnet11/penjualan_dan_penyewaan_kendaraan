package repository

import (
	"context"
	"sultra-otomotif-api/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReviewRepository interface {
	Create(ctx context.Context, review model.Review) (model.Review, error)
	FindByVehicleID(ctx context.Context, vehicleID uuid.UUID) ([]model.Review, error)
}

type reviewRepository struct{ db *pgxpool.Pool }

func NewReviewRepository(db *pgxpool.Pool) ReviewRepository { return &reviewRepository{db: db} }

func (r *reviewRepository) Create(ctx context.Context, rev model.Review) (model.Review, error) {
	query := `INSERT INTO reviews (id, booking_id, user_id, vehicle_id, rating, comment)
              VALUES ($1, $2, $3, $4, $5, $6)
              RETURNING created_at`
	err := r.db.QueryRow(ctx, query, rev.ID, rev.BookingID, rev.UserID, rev.VehicleID, rev.Rating, rev.Comment).Scan(&rev.CreatedAt)
	return rev, err
}

func (r *reviewRepository) FindByVehicleID(ctx context.Context, vehicleID uuid.UUID) ([]model.Review, error) {
	var reviews []model.Review
	query := `SELECT id, booking_id, user_id, vehicle_id, rating, comment, created_at FROM reviews WHERE vehicle_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, vehicleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var rev model.Review
		if err := rows.Scan(&rev.ID, &rev.BookingID, &rev.UserID, &rev.VehicleID, &rev.Rating, &rev.Comment, &rev.CreatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, rev)
	}
	return reviews, nil
}

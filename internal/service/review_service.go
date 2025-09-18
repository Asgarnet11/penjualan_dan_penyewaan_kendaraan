package service

import (
	"context"
	"errors"
	"sultra-otomotif-api/internal/model"
	"sultra-otomotif-api/internal/repository"

	"github.com/google/uuid"
)

type ReviewService interface {
	CreateReview(ctx context.Context, input model.CreateReviewInput, bookingID, userID uuid.UUID) (model.Review, error)
	GetReviewsByVehicleID(ctx context.Context, vehicleID uuid.UUID) ([]model.Review, error)
}

type reviewService struct {
	reviewRepo  repository.ReviewRepository
	bookingRepo repository.BookingRepository
}

func NewReviewService(reviewRepo repository.ReviewRepository, bookingRepo repository.BookingRepository) ReviewService {
	return &reviewService{reviewRepo: reviewRepo, bookingRepo: bookingRepo}
}

func (s *reviewService) CreateReview(ctx context.Context, input model.CreateReviewInput, bookingID, userID uuid.UUID) (model.Review, error) {
	// 1. Ambil data booking untuk divalidasi
	booking, err := s.bookingRepo.FindBookingByID(ctx, bookingID)
	if err != nil {
		return model.Review{}, errors.New("booking not found")
	}

	// 2. Validasi: Apakah user yang login adalah customer yang membuat booking?
	if booking.UserID != userID {
		return model.Review{}, errors.New("forbidden: you can only review your own bookings")
	}

	// 3. Validasi: Apakah status booking sudah 'completed'?
	if booking.Status != "completed" {
		return model.Review{}, errors.New("you can only review a completed booking")
	}

	// Karena ada UNIQUE constraint di DB, error akan otomatis muncul jika review sudah ada.
	// Kita bisa menangani error spesifik dari Postgres (kode 23505) untuk pesan yang lebih baik.

	newReview := model.Review{
		ID:        uuid.New(),
		BookingID: bookingID,
		UserID:    userID,
		VehicleID: booking.VehicleID,
		Rating:    input.Rating,
		Comment:   input.Comment,
	}

	return s.reviewRepo.Create(ctx, newReview)
}

func (s *reviewService) GetReviewsByVehicleID(ctx context.Context, vehicleID uuid.UUID) ([]model.Review, error) {
	return s.reviewRepo.FindByVehicleID(ctx, vehicleID)
}

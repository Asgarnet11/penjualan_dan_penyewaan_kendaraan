package service

import (
	"context"
	"errors"
	"sultra-otomotif-api/internal/model"
	"sultra-otomotif-api/internal/repository"
	"time"

	"github.com/google/uuid"
)

type BookingService interface {
	CreateBooking(ctx context.Context, input model.CreateBookingInput, userID uuid.UUID) (model.Booking, error)
	ConfirmPayment(ctx context.Context, bookingID uuid.UUID) error
	GetBookingsByUserID(ctx context.Context, userID uuid.UUID) ([]model.Booking, error)
	GetBookingsByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]model.Booking, error)
	GetBookingByID(ctx context.Context, bookingID uuid.UUID, currentUserID uuid.UUID) (model.Booking, error)
	UpdateBookingStatus(ctx context.Context, bookingID, currentUserID uuid.UUID, newStatus string) (model.Booking, error)
}

type bookingService struct {
	bookingRepo repository.BookingRepository
	vehicleRepo repository.VehicleRepository
}

func NewBookingService(bookingRepo repository.BookingRepository, vehicleRepo repository.VehicleRepository) BookingService {
	return &bookingService{bookingRepo: bookingRepo, vehicleRepo: vehicleRepo}
}

func (s *bookingService) CreateBooking(ctx context.Context, input model.CreateBookingInput, userID uuid.UUID) (model.Booking, error) {
	vehicleID, err := uuid.Parse(input.VehicleID)
	if err != nil {
		return model.Booking{}, errors.New("invalid vehicle id format")
	}

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, input.StartDate)
	if err != nil {
		return model.Booking{}, errors.New("invalid start_date format, use YYYY-MM-DD")
	}
	endDate, err := time.Parse(layout, input.EndDate)
	if err != nil {
		return model.Booking{}, errors.New("invalid end_date format, use YYYY-MM-DD")
	}
	if endDate.Before(startDate) {
		return model.Booking{}, errors.New("end_date cannot be before start_date")
	}

	available, err := s.bookingRepo.IsVehicleAvailable(ctx, vehicleID, startDate, endDate)
	if err != nil {
		return model.Booking{}, err
	}
	if !available {
		return model.Booking{}, errors.New("vehicle is not available for the selected dates")
	}

	vehicle, err := s.vehicleRepo.FindByID(ctx, vehicleID)
	if err != nil {
		return model.Booking{}, errors.New("vehicle not found")
	}

	durationDays := endDate.Sub(startDate).Hours()/24 + 1
	if durationDays < 1 {
		durationDays = 1
	}

	// PERBAIKAN: Cek apakah harga sewa tidak NULL sebelum digunakan
	if vehicle.RentalPriceDaily == nil {
		return model.Booking{}, errors.New("rental price for this vehicle is not set")
	}
	// Ambil nilai dari pointer
	dailyRate := *vehicle.RentalPriceDaily
	totalPrice := durationDays * dailyRate

	newBooking := model.Booking{
		ID:         uuid.New(),
		UserID:     userID,
		VehicleID:  vehicleID,
		StartDate:  startDate,
		EndDate:    endDate,
		TotalPrice: totalPrice,
		Status:     "pending_payment",
	}

	createdBooking, err := s.bookingRepo.Create(ctx, newBooking)
	if err != nil {
		return model.Booking{}, err
	}

	return createdBooking, nil
}

func (s *bookingService) ConfirmPayment(ctx context.Context, bookingID uuid.UUID) error {
	return s.bookingRepo.UpdateStatus(ctx, bookingID, "confirmed")
}

func (s *bookingService) GetBookingsByUserID(ctx context.Context, userID uuid.UUID) ([]model.Booking, error) {
	return s.bookingRepo.FindBookingsByUserID(ctx, userID)
}

func (s *bookingService) GetBookingsByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]model.Booking, error) {
	return s.bookingRepo.FindBookingsByOwnerID(ctx, ownerID)
}

func (s *bookingService) GetBookingByID(ctx context.Context, bookingID uuid.UUID, currentUserID uuid.UUID) (model.Booking, error) {
	booking, err := s.bookingRepo.FindBookingByID(ctx, bookingID)
	if err != nil {
		return model.Booking{}, errors.New("booking not found")
	}

	vehicle, err := s.vehicleRepo.FindByID(ctx, booking.VehicleID)
	if err != nil {
		return model.Booking{}, errors.New("associated vehicle not found")
	}

	if booking.UserID != currentUserID && vehicle.OwnerID != currentUserID {
		return model.Booking{}, errors.New("forbidden: you are not authorized to view this booking")
	}

	return booking, nil
}

func (s *bookingService) UpdateBookingStatus(ctx context.Context, bookingID, currentUserID uuid.UUID, newStatus string) (model.Booking, error) {
	booking, err := s.bookingRepo.FindBookingByID(ctx, bookingID)
	if err != nil {
		return model.Booking{}, errors.New("booking not found")
	}

	vehicle, err := s.vehicleRepo.FindByID(ctx, booking.VehicleID)
	if err != nil {
		return model.Booking{}, errors.New("associated vehicle not found")
	}
	if vehicle.OwnerID != currentUserID {
		return model.Booking{}, errors.New("forbidden: you are not the owner of this vehicle's booking")
	}

	currentStatus := booking.Status
	isValidTransition := false
	switch currentStatus {
	case "confirmed":
		if newStatus == "rented_out" || newStatus == "cancelled" {
			isValidTransition = true
		}
	case "rented_out":
		if newStatus == "completed" {
			isValidTransition = true
		}
	case "completed", "cancelled":
		return model.Booking{}, errors.New("cannot change status of a completed or cancelled booking")
	}

	if !isValidTransition {
		return model.Booking{}, errors.New("invalid status transition from '" + currentStatus + "' to '" + newStatus + "'")
	}

	err = s.bookingRepo.UpdateStatus(ctx, bookingID, newStatus)
	if err != nil {
		return model.Booking{}, err
	}

	updatedBooking, err := s.bookingRepo.FindBookingByID(ctx, bookingID)
	if err != nil {
		return model.Booking{}, err
	}

	return updatedBooking, nil
}

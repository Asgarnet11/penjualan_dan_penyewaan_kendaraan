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
	vehicleRepo repository.VehicleRepository // Kita butuh ini untuk mengambil harga
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

	// 1. Cek Ketersediaan
	available, err := s.bookingRepo.IsVehicleAvailable(ctx, vehicleID, startDate, endDate)
	if err != nil {
		return model.Booking{}, err
	}
	if !available {
		return model.Booking{}, errors.New("vehicle is not available for the selected dates")
	}

	// 2. Ambil data kendaraan untuk menghitung harga
	vehicle, err := s.vehicleRepo.FindByID(ctx, vehicleID)
	if err != nil {
		return model.Booking{}, errors.New("vehicle not found")
	}

	// 3. Hitung harga total
	durationDays := endDate.Sub(startDate).Hours() / 24
	if durationDays < 1 {
		durationDays = 1
	}
	totalPrice := durationDays * vehicle.RentalPriceDaily

	// 4. Buat objek booking baru
	newBooking := model.Booking{
		ID:         uuid.New(),
		UserID:     userID,
		VehicleID:  vehicleID,
		StartDate:  startDate,
		EndDate:    endDate,
		TotalPrice: totalPrice,
		Status:     "pending_payment",
	}

	// 5. Simpan ke database
	createdBooking, err := s.bookingRepo.Create(ctx, newBooking)
	if err != nil {
		return model.Booking{}, err
	}

	return createdBooking, nil
}

func (s *bookingService) ConfirmPayment(ctx context.Context, bookingID uuid.UUID) error {
	// Di aplikasi nyata, kita akan verifikasi dulu notifikasi dari payment gateway
	// Untuk simulasi, kita langsung update statusnya
	return s.bookingRepo.UpdateStatus(ctx, bookingID, "confirmed")
}

func (s *bookingService) GetBookingsByUserID(ctx context.Context, userID uuid.UUID) ([]model.Booking, error) {
	return s.bookingRepo.FindBookingsByUserID(ctx, userID)
}

// FUNGSI BARU: Mengambil pesanan masuk untuk vendor
func (s *bookingService) GetBookingsByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]model.Booking, error) {
	return s.bookingRepo.FindBookingsByOwnerID(ctx, ownerID)
}

// FUNGSI BARU: Mengambil detail satu booking dengan validasi keamanan
func (s *bookingService) GetBookingByID(ctx context.Context, bookingID uuid.UUID, currentUserID uuid.UUID) (model.Booking, error) {
	// 1. Ambil data booking
	booking, err := s.bookingRepo.FindBookingByID(ctx, bookingID)
	if err != nil {
		return model.Booking{}, errors.New("booking not found")
	}

	// 2. Ambil data kendaraan untuk mendapatkan owner_id
	vehicle, err := s.vehicleRepo.FindByID(ctx, booking.VehicleID)
	if err != nil {
		return model.Booking{}, errors.New("associated vehicle not found")
	}

	// 3. Cek Otorisasi: Apakah user yang login adalah customer yang membuat booking ATAU vendor pemilik kendaraan
	if booking.UserID != currentUserID && vehicle.OwnerID != currentUserID {
		return model.Booking{}, errors.New("forbidden: you are not authorized to view this booking")
	}

	return booking, nil
}

func (s *bookingService) UpdateBookingStatus(ctx context.Context, bookingID, currentUserID uuid.UUID, newStatus string) (model.Booking, error) {
	// 1. Ambil data booking yang ada
	booking, err := s.bookingRepo.FindBookingByID(ctx, bookingID)
	if err != nil {
		return model.Booking{}, errors.New("booking not found")
	}

	// 2. Validasi kepemilikan vendor
	vehicle, err := s.vehicleRepo.FindByID(ctx, booking.VehicleID)
	if err != nil {
		return model.Booking{}, errors.New("associated vehicle not found")
	}
	if vehicle.OwnerID != currentUserID {
		return model.Booking{}, errors.New("forbidden: you are not the owner of this vehicle's booking")
	}

	// 3. Logika State Machine: Terapkan aturan transisi status
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
	// Booking yang sudah selesai atau batal tidak bisa diubah lagi statusnya
	case "completed", "cancelled":
		return model.Booking{}, errors.New("cannot change status of a completed or cancelled booking")
	}

	if !isValidTransition {
		return model.Booking{}, errors.New("invalid status transition from '" + currentStatus + "' to '" + newStatus + "'")
	}

	// 4. Jika valid, update status di repository
	err = s.bookingRepo.UpdateStatus(ctx, bookingID, newStatus)
	if err != nil {
		return model.Booking{}, err
	}

	// 5. Ambil data booking terbaru untuk dikembalikan
	updatedBooking, err := s.bookingRepo.FindBookingByID(ctx, bookingID)
	if err != nil {
		return model.Booking{}, err
	}

	return updatedBooking, nil
}

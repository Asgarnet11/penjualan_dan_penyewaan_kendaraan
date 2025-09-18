package repository

import (
	"context"
	"sultra-otomotif-api/internal/model"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// BookingRepository adalah interface yang akan digunakan oleh service
type BookingRepository interface {
	IsVehicleAvailable(ctx context.Context, vehicleID uuid.UUID, startDate, endDate time.Time) (bool, error)
	Create(ctx context.Context, booking model.Booking) (model.Booking, error)
	FindBookingsByUserID(ctx context.Context, userID uuid.UUID) ([]model.Booking, error)
	FindBookingByID(ctx context.Context, bookingID uuid.UUID) (model.Booking, error)
	UpdateStatus(ctx context.Context, bookingID uuid.UUID, status string) error
	FindBookingsByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]model.Booking, error)
}

// bookingRepository adalah implementasi dari interface di atas
type bookingRepository struct {
	db *pgxpool.Pool
}

// NewBookingRepository adalah constructor untuk bookingRepository
func NewBookingRepository(db *pgxpool.Pool) BookingRepository {
	return &bookingRepository{db: db}
}

// IsVehicleAvailable mengecek apakah ada booking lain yang tumpang tindih pada rentang tanggal tertentu
func (r *bookingRepository) IsVehicleAvailable(ctx context.Context, vehicleID uuid.UUID, startDate, endDate time.Time) (bool, error) {
	var count int
	query := `SELECT count(*) FROM bookings
              WHERE vehicle_id = $1
              AND status IN ('confirmed', 'rented_out')
              AND (start_date, end_date) OVERLAPS ($2, $3)`

	err := r.db.QueryRow(ctx, query, vehicleID, startDate, endDate).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// Create menyimpan data booking baru ke database
func (r *bookingRepository) Create(ctx context.Context, b model.Booking) (model.Booking, error) {
	query := `INSERT INTO bookings (id, user_id, vehicle_id, start_date, end_date, total_price, status, payment_token, payment_url)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
              RETURNING created_at, updated_at`

	// Simulasi pembuatan token & url pembayaran
	b.PaymentToken = "DUMMY-TOKEN-" + b.ID.String()
	b.PaymentURL = "https://ui-avatars.com/api/?name=Bayar+Disini&background=random&size=256&dummy-url=" + b.PaymentToken

	err := r.db.QueryRow(ctx, query, b.ID, b.UserID, b.VehicleID, b.StartDate, b.EndDate, b.TotalPrice, b.Status, b.PaymentToken, b.PaymentURL).Scan(&b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return model.Booking{}, err
	}
	return b, nil
}

// FindBookingsByUserID mengambil semua data booking milik seorang user
func (r *bookingRepository) FindBookingsByUserID(ctx context.Context, userID uuid.UUID) ([]model.Booking, error) {
	var bookings []model.Booking
	query := `SELECT id, user_id, vehicle_id, start_date, end_date, total_price, status, created_at, updated_at FROM bookings WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var b model.Booking
		err := rows.Scan(&b.ID, &b.UserID, &b.VehicleID, &b.StartDate, &b.EndDate, &b.TotalPrice, &b.Status, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}

// FindBookingByID mengambil satu data booking berdasarkan ID-nya
func (r *bookingRepository) FindBookingByID(ctx context.Context, bookingID uuid.UUID) (model.Booking, error) {
	var b model.Booking
	query := `SELECT id, user_id, vehicle_id, start_date, end_date, total_price, status, created_at, updated_at FROM bookings WHERE id = $1`

	err := r.db.QueryRow(ctx, query, bookingID).Scan(&b.ID, &b.UserID, &b.VehicleID, &b.StartDate, &b.EndDate, &b.TotalPrice, &b.Status, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return model.Booking{}, err
	}
	return b, nil
}

// UpdateStatus mengubah status sebuah booking
func (r *bookingRepository) UpdateStatus(ctx context.Context, bookingID uuid.UUID, status string) error {
	query := `UPDATE bookings SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, status, bookingID)
	return err
}

func (r *bookingRepository) FindBookingsByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]model.Booking, error) {
	var bookings []model.Booking
	// Query ini menggunakan JOIN untuk menghubungkan tabel bookings dan vehicles
	query := `SELECT b.id, b.user_id, b.vehicle_id, b.start_date, b.end_date, b.total_price, b.status, b.created_at, b.updated_at
              FROM bookings b
              JOIN vehicles v ON b.vehicle_id = v.id
              WHERE v.owner_id = $1
              ORDER BY b.created_at DESC`

	rows, err := r.db.Query(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var b model.Booking
		err := rows.Scan(&b.ID, &b.UserID, &b.VehicleID, &b.StartDate, &b.EndDate, &b.TotalPrice, &b.Status, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}

package handler

import (
	"net/http"
	"strings"
	"sultra-otomotif-api/internal/helper"
	"sultra-otomotif-api/internal/model"
	"sultra-otomotif-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookingHandler struct {
	bookingService service.BookingService
}

type StatusInput struct {
	Status string `json:"status" binding:"required,oneof=rented_out completed cancelled"`
}

func NewBookingHandler(bookingService service.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}

// CreateBooking menangani permintaan untuk membuat booking baru
func (h *BookingHandler) CreateBooking(ctx *gin.Context) {
	var input model.CreateBookingInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(ctx, "Invalid input data", http.StatusBadRequest, err)
		return
	}

	// Ambil userID dari context yang sudah di-set oleh middleware
	currentUserID := ctx.MustGet("currentUserID").(uuid.UUID)

	booking, err := h.bookingService.CreateBooking(ctx, input, currentUserID)
	if err != nil {
		// Beri status code yang lebih sesuai jika kendaraan tidak tersedia
		if err.Error() == "vehicle is not available for the selected dates" {
			helper.ErrorResponse(ctx, err.Error(), http.StatusConflict, err) // 409 Conflict
		} else {
			helper.ErrorResponse(ctx, err.Error(), http.StatusBadRequest, err)
		}
		return
	}
	helper.APIResponse(ctx, "Booking created successfully, waiting for payment", http.StatusCreated, booking)
}

// PaymentCallbackInput adalah struct untuk menampung data dari webhook payment gateway
type PaymentCallbackInput struct {
	BookingID string `json:"booking_id" binding:"required"`
	Status    string `json:"status" binding:"required"`
}

// PaymentCallback menangani notifikasi (webhook) dari payment gateway
func (h *BookingHandler) PaymentCallback(ctx *gin.Context) {
	var input PaymentCallbackInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(ctx, "Invalid callback data", http.StatusBadRequest, err)
		return
	}

	bookingID, err := uuid.Parse(input.BookingID)
	if err != nil {
		helper.ErrorResponse(ctx, "Invalid booking ID format", http.StatusBadRequest, err)
		return
	}

	// Di aplikasi nyata, akan ada verifikasi token/signature dari payment gateway
	// Untuk simulasi, kita cek status "success"
	if input.Status == "success" {
		err := h.bookingService.ConfirmPayment(ctx, bookingID)
		if err != nil {
			helper.ErrorResponse(ctx, "Failed to confirm payment", http.StatusInternalServerError, err)
			return
		}
	} else {
		// Anda bisa menambahkan logika untuk handle pembayaran gagal (misal: ubah status jadi 'cancelled')
	}

	helper.APIResponse(ctx, "Payment callback processed successfully", http.StatusOK, nil)
}

// GetMyBookings menangani permintaan untuk melihat riwayat booking milik user yang sedang login
func (h *BookingHandler) GetMyBookings(ctx *gin.Context) {
	currentUserID := ctx.MustGet("currentUserID").(uuid.UUID)

	bookings, err := h.bookingService.GetBookingsByUserID(ctx, currentUserID)
	if err != nil {
		helper.ErrorResponse(ctx, "Failed to fetch bookings", http.StatusInternalServerError, err)
		return
	}
	helper.APIResponse(ctx, "Successfully fetched user bookings", http.StatusOK, bookings)
}

func (h *BookingHandler) GetVendorBookings(ctx *gin.Context) {
	currentUserID := ctx.MustGet("currentUserID").(uuid.UUID)

	bookings, err := h.bookingService.GetBookingsByOwnerID(ctx, currentUserID)
	if err != nil {
		helper.ErrorResponse(ctx, "Failed to fetch vendor bookings", http.StatusInternalServerError, err)
		return
	}
	helper.APIResponse(ctx, "Successfully fetched vendor bookings", http.StatusOK, bookings)
}

func (h *BookingHandler) GetBookingByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helper.ErrorResponse(ctx, "Invalid booking ID", http.StatusBadRequest, err)
		return
	}

	currentUserID := ctx.MustGet("currentUserID").(uuid.UUID)

	booking, err := h.bookingService.GetBookingByID(ctx, id, currentUserID)
	if err != nil {
		if err.Error() == "forbidden: you are not authorized to view this booking" {
			helper.ErrorResponse(ctx, err.Error(), http.StatusForbidden, err)
		} else {
			helper.ErrorResponse(ctx, err.Error(), http.StatusNotFound, err)
		}
		return
	}
	helper.APIResponse(ctx, "Successfully fetched booking detail", http.StatusOK, booking)
}

func (h *BookingHandler) UpdateBookingStatus(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helper.ErrorResponse(ctx, "Invalid booking ID", http.StatusBadRequest, err)
		return
	}

	var input StatusInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(ctx, "Invalid input data. Status must be one of: rented_out, completed, cancelled", http.StatusBadRequest, err)
		return
	}

	currentUserID := ctx.MustGet("currentUserID").(uuid.UUID)

	updatedBooking, err := h.bookingService.UpdateBookingStatus(ctx, id, currentUserID, input.Status)
	if err != nil {
		if strings.HasPrefix(err.Error(), "forbidden") {
			helper.ErrorResponse(ctx, err.Error(), http.StatusForbidden, err)
		} else if strings.HasPrefix(err.Error(), "invalid status transition") {
			helper.ErrorResponse(ctx, err.Error(), http.StatusConflict, err) // 409 Conflict
		} else {
			helper.ErrorResponse(ctx, err.Error(), http.StatusNotFound, err)
		}
		return
	}

	helper.APIResponse(ctx, "Booking status updated successfully", http.StatusOK, updatedBooking)
}

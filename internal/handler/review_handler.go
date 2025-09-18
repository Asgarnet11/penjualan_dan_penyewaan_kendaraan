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

// ReviewHandler adalah struct untuk review handler
type ReviewHandler struct {
	reviewService service.ReviewService
}

// NewReviewHandler adalah constructor untuk ReviewHandler
func NewReviewHandler(reviewService service.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

// CreateReview menangani permintaan untuk membuat ulasan baru
func (h *ReviewHandler) CreateReview(ctx *gin.Context) {
	bookingID, err := uuid.Parse(ctx.Param("booking_id"))
	if err != nil {
		helper.ErrorResponse(ctx, "Invalid booking ID", http.StatusBadRequest, err)
		return
	}

	var input model.CreateReviewInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(ctx, "Invalid input data", http.StatusBadRequest, err)
		return
	}

	currentUserID := ctx.MustGet("currentUserID").(uuid.UUID)

	review, err := h.reviewService.CreateReview(ctx, input, bookingID, currentUserID)
	if err != nil {
		// Memberikan status code yang lebih spesifik berdasarkan jenis error dari service
		if strings.HasPrefix(err.Error(), "forbidden") {
			helper.ErrorResponse(ctx, err.Error(), http.StatusForbidden, err)
		} else if err.Error() == "booking not found" {
			helper.ErrorResponse(ctx, err.Error(), http.StatusNotFound, err)
		} else if err.Error() == "you can only review a completed booking" {
			helper.ErrorResponse(ctx, err.Error(), http.StatusConflict, err)
		} else if strings.Contains(err.Error(), "duplicate key") { // Error dari unique constraint di DB
			helper.ErrorResponse(ctx, "a review for this booking already exists", http.StatusConflict, err)
		} else {
			helper.ErrorResponse(ctx, err.Error(), http.StatusInternalServerError, err)
		}
		return
	}

	helper.APIResponse(ctx, "Review created successfully", http.StatusCreated, review)
}

// GetVehicleReviews menangani permintaan untuk mengambil semua ulasan dari sebuah kendaraan
func (h *ReviewHandler) GetVehicleReviews(ctx *gin.Context) {
	vehicleID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helper.ErrorResponse(ctx, "Invalid vehicle ID", http.StatusBadRequest, err)
		return
	}

	reviews, err := h.reviewService.GetReviewsByVehicleID(ctx, vehicleID)
	if err != nil {
		helper.ErrorResponse(ctx, "Failed to fetch reviews", http.StatusInternalServerError, err)
		return
	}

	helper.APIResponse(ctx, "Successfully fetched vehicle reviews", http.StatusOK, reviews)
}

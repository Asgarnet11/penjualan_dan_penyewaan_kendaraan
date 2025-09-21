package handler

import (
	"net/http"
	"sultra-otomotif-api/internal/helper"
	"sultra-otomotif-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SalesHandler struct {
	salesService service.SalesService
}

func NewSalesHandler(salesService service.SalesService) *SalesHandler {
	return &SalesHandler{salesService: salesService}
}

func (h *SalesHandler) InitiatePurchase(ctx *gin.Context) {
	vehicleID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helper.ErrorResponse(ctx, "Invalid vehicle ID", http.StatusBadRequest, err)
		return
	}

	buyerID := ctx.MustGet("currentUserID").(uuid.UUID)

	transaction, err := h.salesService.InitiatePurchase(ctx, vehicleID, buyerID)
	if err != nil {
		helper.ErrorResponse(ctx, err.Error(), http.StatusConflict, err)
		return
	}

	helper.APIResponse(ctx, "Purchase initiated, waiting for payment", http.StatusCreated, transaction)
}

type SalesCallbackInput struct {
	TransactionID string `json:"transaction_id" binding:"required"`
	Status        string `json:"status" binding:"required"`
}

func (h *SalesHandler) PaymentCallback(ctx *gin.Context) {
	var input SalesCallbackInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(ctx, "Invalid callback data", http.StatusBadRequest, err)
		return
	}

	transactionID, err := uuid.Parse(input.TransactionID)
	if err != nil {
		helper.ErrorResponse(ctx, "Invalid transaction ID format", http.StatusBadRequest, err)
		return
	}

	if input.Status == "success" {
		err := h.salesService.ConfirmSale(ctx, transactionID)
		if err != nil {
			helper.ErrorResponse(ctx, "Failed to confirm sale", http.StatusInternalServerError, err)
			return
		}
	}

	helper.APIResponse(ctx, "Sales callback processed successfully", http.StatusOK, nil)
}

func (h *SalesHandler) GetMyPurchases(ctx *gin.Context) {
	buyerID := ctx.MustGet("currentUserID").(uuid.UUID)
	transactions, err := h.salesService.GetPurchasesByBuyerID(ctx, buyerID)
	if err != nil {
		helper.ErrorResponse(ctx, "Failed to fetch purchase history", http.StatusInternalServerError, err)
		return
	}
	helper.APIResponse(ctx, "Successfully fetched purchase history", http.StatusOK, transactions)
}

func (h *SalesHandler) GetMySales(ctx *gin.Context) {
	sellerID := ctx.MustGet("currentUserID").(uuid.UUID)
	transactions, err := h.salesService.GetSalesBySellerID(ctx, sellerID)
	if err != nil {
		helper.ErrorResponse(ctx, "Failed to fetch sales history", http.StatusInternalServerError, err)
		return
	}
	helper.APIResponse(ctx, "Successfully fetched sales history", http.StatusOK, transactions)
}

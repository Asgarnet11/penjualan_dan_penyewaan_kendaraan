package handler

import (
	"net/http"
	"strings"
	"sultra-otomotif-api/internal/helper"
	"sultra-otomotif-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminHandler struct {
	adminService service.AdminService
}

func NewAdminHandler(adminService service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

func (h *AdminHandler) GetVendors(ctx *gin.Context) {
	vendors, err := h.adminService.GetVendors(ctx)
	if err != nil {
		helper.ErrorResponse(ctx, "Failed to fetch vendors", http.StatusInternalServerError, err)
		return
	}
	helper.APIResponse(ctx, "Successfully fetched all vendors", http.StatusOK, vendors)
}

func (h *AdminHandler) VerifyVendor(ctx *gin.Context) {
	vendorID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helper.ErrorResponse(ctx, "Invalid vendor ID", http.StatusBadRequest, err)
		return
	}

	updatedVendor, err := h.adminService.VerifyVendor(ctx, vendorID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			helper.ErrorResponse(ctx, err.Error(), http.StatusNotFound, err)
		} else {
			helper.ErrorResponse(ctx, err.Error(), http.StatusBadRequest, err)
		}
		return
	}

	helper.APIResponse(ctx, "Vendor verified successfully", http.StatusOK, updatedVendor)
}

func (h *AdminHandler) GetAllUsers(ctx *gin.Context) {
	users, err := h.adminService.GetAllUsers(ctx)
	if err != nil {
		helper.ErrorResponse(ctx, "Failed to fetch users", http.StatusInternalServerError, err)
		return
	}
	helper.APIResponse(ctx, "Successfully fetched all users", http.StatusOK, users)
}

// FUNGSI BARU:
func (h *AdminHandler) DeleteUser(ctx *gin.Context) {
	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helper.ErrorResponse(ctx, "Invalid user ID", http.StatusBadRequest, err)
		return
	}

	err = h.adminService.DeleteUser(ctx, userID)
	if err != nil {
		helper.ErrorResponse(ctx, "Failed to delete user", http.StatusInternalServerError, err)
		return
	}
	helper.APIResponse(ctx, "User deleted successfully", http.StatusOK, nil)
}

func (h *AdminHandler) GetAllVehicles(ctx *gin.Context) {
	vehicles, err := h.adminService.GetAllVehicles(ctx)
	if err != nil {
		helper.ErrorResponse(ctx, "Failed to fetch vehicles", http.StatusInternalServerError, err)
		return
	}
	helper.APIResponse(ctx, "Successfully fetched all vehicles", http.StatusOK, vehicles)
}

// FUNGSI BARU:
func (h *AdminHandler) DeleteVehicle(ctx *gin.Context) {
	vehicleID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helper.ErrorResponse(ctx, "Invalid vehicle ID", http.StatusBadRequest, err)
		return
	}

	err = h.adminService.DeleteVehicle(ctx, vehicleID)
	if err != nil {
		helper.ErrorResponse(ctx, "Failed to delete vehicle", http.StatusInternalServerError, err)
		return
	}
	helper.APIResponse(ctx, "Vehicle deleted successfully", http.StatusOK, nil)
}

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

type VehicleHandler struct {
	vehicleService service.VehicleService
}

// INI FUNGSI YANG HILANG: NewVehicleHandler
func NewVehicleHandler(vehicleService service.VehicleService) *VehicleHandler {
	return &VehicleHandler{vehicleService: vehicleService}
}

func (h *VehicleHandler) CreateVehicle(ctx *gin.Context) {
	var input model.CreateVehicleInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(ctx, "Invalid input data", http.StatusBadRequest, err)
		return
	}

	// Ambil userID dari context yang sudah di-set oleh middleware
	currentUserID := ctx.MustGet("currentUserID").(uuid.UUID)

	vehicle, err := h.vehicleService.CreateVehicle(ctx, input, currentUserID)
	if err != nil {
		helper.ErrorResponse(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}
	helper.APIResponse(ctx, "Vehicle created successfully", http.StatusCreated, vehicle)
}

func (h *VehicleHandler) GetAllVehicles(ctx *gin.Context) {
	vehicles, err := h.vehicleService.GetAllVehicles(ctx)
	if err != nil {
		helper.ErrorResponse(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}
	helper.APIResponse(ctx, "Successfully fetched all vehicles", http.StatusOK, vehicles)
}

func (h *VehicleHandler) GetVehicleByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helper.ErrorResponse(ctx, "Invalid vehicle ID", http.StatusBadRequest, err)
		return
	}

	vehicle, err := h.vehicleService.GetVehicleByID(ctx, id)
	if err != nil {
		helper.ErrorResponse(ctx, "Vehicle not found", http.StatusNotFound, err)
		return
	}
	helper.APIResponse(ctx, "Successfully fetched vehicle", http.StatusOK, vehicle)
}

func (h *VehicleHandler) UpdateVehicle(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helper.ErrorResponse(ctx, "Invalid vehicle ID", http.StatusBadRequest, err)
		return
	}

	var input model.CreateVehicleInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(ctx, "Invalid input data", http.StatusBadRequest, err)
		return
	}

	currentUserID := ctx.MustGet("currentUserID").(uuid.UUID)

	vehicle, err := h.vehicleService.UpdateVehicle(ctx, id, currentUserID, input)
	if err != nil {
		// Cek jenis error untuk memberikan status code yang sesuai
		if err.Error() == "forbidden: you are not the owner of this vehicle" {
			helper.ErrorResponse(ctx, err.Error(), http.StatusForbidden, err)
		} else if err.Error() == "vehicle not found" {
			helper.ErrorResponse(ctx, err.Error(), http.StatusNotFound, err)
		} else {
			helper.ErrorResponse(ctx, err.Error(), http.StatusInternalServerError, err)
		}
		return
	}
	helper.APIResponse(ctx, "Vehicle updated successfully", http.StatusOK, vehicle)
}

func (h *VehicleHandler) DeleteVehicle(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helper.ErrorResponse(ctx, "Invalid vehicle ID", http.StatusBadRequest, err)
		return
	}

	currentUserID := ctx.MustGet("currentUserID").(uuid.UUID)

	err = h.vehicleService.DeleteVehicle(ctx, id, currentUserID)
	if err != nil {
		if err.Error() == "forbidden: you are not the owner of this vehicle" {
			helper.ErrorResponse(ctx, err.Error(), http.StatusForbidden, err)
		} else if err.Error() == "vehicle not found" {
			helper.ErrorResponse(ctx, err.Error(), http.StatusNotFound, err)
		} else {
			helper.ErrorResponse(ctx, err.Error(), http.StatusInternalServerError, err)
		}
		return
	}
	helper.APIResponse(ctx, "Vehicle deleted successfully", http.StatusOK, nil)
}

func (h *VehicleHandler) UploadVehicleImage(ctx *gin.Context) {
	vehicleID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helper.ErrorResponse(ctx, "Invalid vehicle ID", http.StatusBadRequest, err)
		return
	}

	currentUserID := ctx.MustGet("currentUserID").(uuid.UUID)

	// Ambil file dari form request
	file, err := ctx.FormFile("image")
	if err != nil {
		helper.ErrorResponse(ctx, "Image file is required", http.StatusBadRequest, err)
		return
	}

	// (Opsional) Validasi tipe file dan ukuran
	// ...
	filename := uuid.New().String() + "-" + file.Filename

	// Panggil service untuk memproses upload
	imageURL, err := h.vehicleService.UploadImage(ctx, vehicleID, currentUserID, filename)
	if err != nil {
		if strings.HasPrefix(err.Error(), "forbidden") {
			helper.ErrorResponse(ctx, err.Error(), http.StatusForbidden, err)
		} else {
			helper.ErrorResponse(ctx, err.Error(), http.StatusNotFound, err)
		}
		return
	}

	response := gin.H{"image_url": imageURL}
	helper.APIResponse(ctx, "Image uploaded successfully", http.StatusOK, response)
}

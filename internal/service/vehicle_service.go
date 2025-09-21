package service

import (
	"context"
	"errors"
	"sultra-otomotif-api/internal/model"
	"sultra-otomotif-api/internal/repository"

	"github.com/google/uuid"
)

type VehicleService interface {
	CreateVehicle(ctx context.Context, input model.CreateVehicleInput, ownerID uuid.UUID) (model.Vehicle, error)
	GetAllVehicles(ctx context.Context, filter model.VehicleFilter) ([]model.Vehicle, error)
	GetVehicleByID(ctx context.Context, id uuid.UUID) (model.Vehicle, error)
	UpdateVehicle(ctx context.Context, id uuid.UUID, currentUserID uuid.UUID, input model.CreateVehicleInput) (model.Vehicle, error)
	DeleteVehicle(ctx context.Context, id uuid.UUID, currentUserID uuid.UUID) error
	UploadImage(ctx context.Context, vehicleID, currentUserID uuid.UUID, filename string) (string, error)
}

type vehicleService struct {
	repo      repository.VehicleRepository
	imageRepo repository.ImageRepository
	userRepo  repository.UserRepository
}

func NewVehicleService(repo repository.VehicleRepository, imageRepo repository.ImageRepository, userRepo repository.UserRepository) VehicleService {
	return &vehicleService{repo: repo, imageRepo: imageRepo, userRepo: userRepo}
}

func (s *vehicleService) CreateVehicle(ctx context.Context, input model.CreateVehicleInput, ownerID uuid.UUID) (model.Vehicle, error) {

	owner, err := s.userRepo.FindByID(ctx, ownerID) // Kita perlu fungsi FindByID di userRepo
	if err != nil {
		return model.Vehicle{}, errors.New("owner not found")
	}
	if !owner.IsVerified {
		return model.Vehicle{}, errors.New("forbidden: vendor account is not verified")
	}

	newVehicle := model.Vehicle{
		ID:                 uuid.New(),
		OwnerID:            ownerID, // Diambil dari token JWT
		Brand:              input.Brand,
		Model:              input.Model,
		Year:               input.Year,
		PlateNumber:        input.PlateNumber,
		Color:              input.Color,
		VehicleType:        input.VehicleType,
		Transmission:       input.Transmission,
		Fuel:               input.Fuel,
		Status:             "available", // Status default saat dibuat
		Description:        input.Description,
		IsForSale:          input.IsForSale,
		SalePrice:          input.SalePrice,
		IsForRent:          input.IsForRent,
		RentalPriceDaily:   input.RentalPriceDaily,
		RentalPriceWeekly:  input.RentalPriceWeekly,
		RentalPriceMonthly: input.RentalPriceMonthly,
	}

	createdVehicle, err := s.repo.Create(ctx, newVehicle)
	if err != nil {
		return model.Vehicle{}, err
	}

	return createdVehicle, nil
}

func (s *vehicleService) GetAllVehicles(ctx context.Context, filter model.VehicleFilter) ([]model.Vehicle, error) { // <-- Gunakan model.VehicleFilter
	return s.repo.FindAll(ctx, filter)
}

func (s *vehicleService) GetVehicleByID(ctx context.Context, id uuid.UUID) (model.Vehicle, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *vehicleService) UpdateVehicle(ctx context.Context, id uuid.UUID, currentUserID uuid.UUID, input model.CreateVehicleInput) (model.Vehicle, error) {
	// PENTING: Cek dulu apakah kendaraan ini ada
	vehicleToUpdate, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return model.Vehicle{}, errors.New("vehicle not found")
	}

	// PENTING: Validasi kepemilikan
	if vehicleToUpdate.OwnerID != currentUserID {
		return model.Vehicle{}, errors.New("forbidden: you are not the owner of this vehicle")
	}

	// Update field
	vehicleToUpdate.Brand = input.Brand
	vehicleToUpdate.Model = input.Model
	vehicleToUpdate.Year = input.Year
	vehicleToUpdate.PlateNumber = input.PlateNumber
	vehicleToUpdate.Color = input.Color
	vehicleToUpdate.VehicleType = input.VehicleType
	vehicleToUpdate.Transmission = input.Transmission
	vehicleToUpdate.Fuel = input.Fuel
	vehicleToUpdate.Description = input.Description
	vehicleToUpdate.IsForSale = input.IsForSale
	vehicleToUpdate.SalePrice = input.SalePrice
	vehicleToUpdate.IsForRent = input.IsForRent
	vehicleToUpdate.RentalPriceDaily = input.RentalPriceDaily
	vehicleToUpdate.RentalPriceWeekly = input.RentalPriceWeekly
	vehicleToUpdate.RentalPriceMonthly = input.RentalPriceMonthly

	updatedVehicle, err := s.repo.Update(ctx, vehicleToUpdate)
	if err != nil {
		return model.Vehicle{}, err
	}

	return updatedVehicle, nil
}

func (s *vehicleService) DeleteVehicle(ctx context.Context, id uuid.UUID, currentUserID uuid.UUID) error {
	// PENTING: Cek dulu apakah kendaraan ini ada
	vehicleToDelete, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("vehicle not found")
	}

	// PENTING: Validasi kepemilikan
	if vehicleToDelete.OwnerID != currentUserID {
		return errors.New("forbidden: you are not the owner of this vehicle")
	}

	return s.repo.Delete(ctx, id)
}

func (s *vehicleService) UploadImage(ctx context.Context, vehicleID, currentUserID uuid.UUID, filename string) (string, error) {
	// 1. Validasi kepemilikan kendaraan
	vehicle, err := s.repo.FindByID(ctx, vehicleID)
	if err != nil {
		return "", errors.New("vehicle not found")
	}
	if vehicle.OwnerID != currentUserID {
		return "", errors.New("forbidden: you are not the owner of this vehicle")
	}

	// 2. (SIMULASI) Proses upload ke cloud
	// Di aplikasi nyata, di sinilah kode untuk upload ke AWS S3, dll.
	// Untuk sekarang, kita buat URL placeholder.
	imageURL := "https://storage.googleapis.com/sultra-otomotif-bucket/images/" + filename

	// 3. Simpan URL ke database
	err = s.imageRepo.SaveVehicleImage(ctx, vehicleID, imageURL)
	if err != nil {
		return "", err
	}

	return imageURL, nil
}

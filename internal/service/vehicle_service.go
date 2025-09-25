package service

import (
	"context"
	"errors"
	"mime/multipart"
	"sultra-otomotif-api/internal/config"
	"sultra-otomotif-api/internal/model"
	"sultra-otomotif-api/internal/repository"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

// Helper function untuk membuat pointer dari string, mengembalikan nil jika string kosong
func stringToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Helper function untuk membuat pointer dari float64, mengembalikan nil jika float64 adalah 0
func float64ToPtr(f float64) *float64 {
	if f == 0 {
		return nil
	}
	return &f
}

type VehicleService interface {
	CreateVehicle(ctx context.Context, input model.CreateVehicleInput, ownerID uuid.UUID) (model.Vehicle, error)
	GetAllVehicles(ctx context.Context, filter model.VehicleFilter) ([]model.Vehicle, error)
	GetVehicleByID(ctx context.Context, id uuid.UUID) (model.Vehicle, error)
	UpdateVehicle(ctx context.Context, id uuid.UUID, currentUserID uuid.UUID, input model.CreateVehicleInput) (model.Vehicle, error)
	DeleteVehicle(ctx context.Context, id uuid.UUID, currentUserID uuid.UUID) error
	UploadImage(ctx context.Context, vehicleID, currentUserID uuid.UUID, file multipart.File) (string, error)
	GetVehiclesByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]model.Vehicle, error)
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
	owner, err := s.userRepo.FindByID(ctx, ownerID)
	if err != nil {
		return model.Vehicle{}, errors.New("owner not found")
	}
	if !owner.IsVerified {
		return model.Vehicle{}, errors.New("forbidden: vendor account is not verified")
	}

	newVehicle := model.Vehicle{
		ID:           uuid.New(),
		OwnerID:      ownerID,
		Brand:        input.Brand,
		Model:        input.Model,
		Year:         input.Year,
		PlateNumber:  input.PlateNumber,
		VehicleType:  input.VehicleType,
		Transmission: input.Transmission,
		Fuel:         input.Fuel,
		Status:       "available",
		IsForSale:    input.IsForSale,
		IsForRent:    input.IsForRent,
		Features:     input.Features,
		// PERBAIKAN: Gunakan helper untuk mengisi field pointer
		Color:              stringToPtr(input.Color),
		Description:        stringToPtr(input.Description),
		Location:           stringToPtr(input.Location),
		SalePrice:          float64ToPtr(input.SalePrice),
		RentalPriceDaily:   float64ToPtr(input.RentalPriceDaily),
		RentalPriceWeekly:  float64ToPtr(input.RentalPriceWeekly),
		RentalPriceMonthly: float64ToPtr(input.RentalPriceMonthly),
	}

	createdVehicle, err := s.repo.Create(ctx, newVehicle)
	if err != nil {
		return model.Vehicle{}, err
	}
	return createdVehicle, nil
}

func (s *vehicleService) GetAllVehicles(ctx context.Context, filter model.VehicleFilter) ([]model.Vehicle, error) {
	return s.repo.FindAll(ctx, filter)
}

func (s *vehicleService) GetVehicleByID(ctx context.Context, id uuid.UUID) (model.Vehicle, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *vehicleService) UpdateVehicle(ctx context.Context, id uuid.UUID, currentUserID uuid.UUID, input model.CreateVehicleInput) (model.Vehicle, error) {
	vehicleToUpdate, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return model.Vehicle{}, errors.New("vehicle not found")
	}

	if vehicleToUpdate.OwnerID != currentUserID {
		return model.Vehicle{}, errors.New("forbidden: you are not the owner of this vehicle")
	}

	vehicleToUpdate.Brand = input.Brand
	vehicleToUpdate.Model = input.Model
	vehicleToUpdate.Year = input.Year
	vehicleToUpdate.PlateNumber = input.PlateNumber
	vehicleToUpdate.VehicleType = input.VehicleType
	vehicleToUpdate.Transmission = input.Transmission
	vehicleToUpdate.Fuel = input.Fuel
	vehicleToUpdate.IsForSale = input.IsForSale
	vehicleToUpdate.IsForRent = input.IsForRent
	vehicleToUpdate.Features = input.Features

	// PERBAIKAN: Gunakan helper untuk memperbarui field pointer
	vehicleToUpdate.Color = stringToPtr(input.Color)
	vehicleToUpdate.Description = stringToPtr(input.Description)
	vehicleToUpdate.Location = stringToPtr(input.Location)
	vehicleToUpdate.SalePrice = float64ToPtr(input.SalePrice)
	vehicleToUpdate.RentalPriceDaily = float64ToPtr(input.RentalPriceDaily)
	vehicleToUpdate.RentalPriceWeekly = float64ToPtr(input.RentalPriceWeekly)
	vehicleToUpdate.RentalPriceMonthly = float64ToPtr(input.RentalPriceMonthly)

	updatedVehicle, err := s.repo.Update(ctx, vehicleToUpdate)
	if err != nil {
		return model.Vehicle{}, err
	}
	return updatedVehicle, nil
}

func (s *vehicleService) DeleteVehicle(ctx context.Context, id uuid.UUID, currentUserID uuid.UUID) error {
	vehicleToDelete, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("vehicle not found")
	}

	if vehicleToDelete.OwnerID != currentUserID {
		return errors.New("forbidden: you are not the owner of this vehicle")
	}

	return s.repo.Delete(ctx, id)
}

func (s *vehicleService) UploadImage(ctx context.Context, vehicleID, currentUserID uuid.UUID, file multipart.File) (string, error) {
	vehicle, err := s.repo.FindByID(ctx, vehicleID)
	if err != nil {
		return "", errors.New("vehicle not found")
	}
	if vehicle.OwnerID != currentUserID {
		return "", errors.New("forbidden: you are not the owner of this vehicle")
	}

	cfg := config.LoadConfig()
	cld, err := cloudinary.NewFromURL(cfg.CloudinaryURL)
	if err != nil {
		return "", errors.New("failed to connect to cloudinary")
	}

	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: "sultra-otomotif",
	})
	if err != nil {
		return "", errors.New("failed to upload image to cloudinary")
	}

	imageURL := uploadResult.SecureURL
	err = s.imageRepo.SaveVehicleImage(ctx, vehicleID, imageURL)
	if err != nil {
		return "", err
	}

	return imageURL, nil
}

func (s *vehicleService) GetVehiclesByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]model.Vehicle, error) {
	return s.repo.FindAllByOwnerID(ctx, ownerID)
}

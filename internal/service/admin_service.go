package service

import (
	"context"
	"errors"
	"sultra-otomotif-api/internal/model"
	"sultra-otomotif-api/internal/repository"

	"github.com/google/uuid"
)

type AdminService interface {
	GetVendors(ctx context.Context) ([]model.User, error)
	VerifyVendor(ctx context.Context, vendorID uuid.UUID) (model.User, error)
	GetAllUsers(ctx context.Context) ([]model.User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	GetAllVehicles(ctx context.Context) ([]model.Vehicle, error)
	DeleteVehicle(ctx context.Context, vehicleID uuid.UUID) error
}

type adminService struct {
	userRepo    repository.UserRepository
	vehicleRepo repository.VehicleRepository
}

func NewAdminService(userRepo repository.UserRepository, vehicleRepo repository.VehicleRepository) AdminService {
	return &adminService{userRepo: userRepo, vehicleRepo: vehicleRepo}
}

func (s *adminService) GetVendors(ctx context.Context) ([]model.User, error) {
	return s.userRepo.FindUsersByRole(ctx, "vendor")
}

func (s *adminService) VerifyVendor(ctx context.Context, vendorID uuid.UUID) (model.User, error) {
	vendor, err := s.userRepo.FindByID(ctx, vendorID)
	if err != nil {
		return model.User{}, errors.New("vendor not found")
	}

	if vendor.Role != "vendor" {
		return model.User{}, errors.New("this user is not a vendor")
	}

	if vendor.IsVerified {
		return vendor, nil
	}

	err = s.userRepo.UpdateVerificationStatus(ctx, vendorID, true)
	if err != nil {
		return model.User{}, err
	}

	updatedVendor, err := s.userRepo.FindByID(ctx, vendorID)
	if err != nil {
		return model.User{}, errors.New("failed to fetch updated vendor data")
	}

	return updatedVendor, nil
}

func (s *adminService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	return s.userRepo.FindAll(ctx)
}

func (s *adminService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	// Di sini bisa ditambahkan logika tambahan, misal logging siapa yang menghapus
	return s.userRepo.Delete(ctx, userID)
}

func (s *adminService) GetAllVehicles(ctx context.Context) ([]model.Vehicle, error) {
	return s.vehicleRepo.FindAllAdmin(ctx)
}

// FUNGSI BARU:
func (s *adminService) DeleteVehicle(ctx context.Context, vehicleID uuid.UUID) error {
	return s.vehicleRepo.Delete(ctx, vehicleID)
}

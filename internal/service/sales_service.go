package service

import (
	"context"
	"errors"
	"sultra-otomotif-api/internal/model"
	"sultra-otomotif-api/internal/repository"

	"github.com/google/uuid"
)

type SalesService interface {
	InitiatePurchase(ctx context.Context, vehicleID, buyerID uuid.UUID) (model.SalesTransaction, error)
	ConfirmSale(ctx context.Context, transactionID uuid.UUID) error
	GetPurchasesByBuyerID(ctx context.Context, buyerID uuid.UUID) ([]model.SalesTransaction, error)
	GetSalesBySellerID(ctx context.Context, sellerID uuid.UUID) ([]model.SalesTransaction, error)
}

type salesService struct {
	salesRepo   repository.SalesRepository
	vehicleRepo repository.VehicleRepository
}

func NewSalesService(salesRepo repository.SalesRepository, vehicleRepo repository.VehicleRepository) SalesService {
	return &salesService{salesRepo: salesRepo, vehicleRepo: vehicleRepo}
}

func (s *salesService) InitiatePurchase(ctx context.Context, vehicleID, buyerID uuid.UUID) (model.SalesTransaction, error) {
	vehicle, err := s.vehicleRepo.FindByID(ctx, vehicleID)
	if err != nil {
		return model.SalesTransaction{}, errors.New("vehicle not found")
	}

	if !vehicle.IsForSale {
		return model.SalesTransaction{}, errors.New("this vehicle is not for sale")
	}
	if vehicle.Status != "available" {
		return model.SalesTransaction{}, errors.New("this vehicle is no longer available")
	}
	if vehicle.OwnerID == buyerID {
		return model.SalesTransaction{}, errors.New("you cannot buy your own vehicle")
	}

	// PERBAIKAN DI SINI:
	// Cek apakah harga jual tidak NULL sebelum digunakan
	if vehicle.SalePrice == nil {
		return model.SalesTransaction{}, errors.New("sale price for this vehicle is not set")
	}
	// Ambil nilai dari pointer
	agreedPrice := *vehicle.SalePrice

	newTransaction := model.SalesTransaction{
		ID:          uuid.New(),
		VehicleID:   vehicleID,
		SellerID:    vehicle.OwnerID,
		BuyerID:     buyerID,
		AgreedPrice: agreedPrice, // Gunakan nilai yang sudah di-dereference
		Status:      "payment_pending",
	}

	return s.salesRepo.Create(ctx, newTransaction)
}

func (s *salesService) ConfirmSale(ctx context.Context, transactionID uuid.UUID) error {
	transaction, err := s.salesRepo.UpdateStatus(ctx, transactionID, "completed")
	if err != nil {
		return err
	}

	vehicle, err := s.vehicleRepo.FindByID(ctx, transaction.VehicleID)
	if err != nil {
		return err
	}

	vehicle.Status = "sold"
	vehicle.IsForSale = false
	vehicle.IsForRent = false

	_, err = s.vehicleRepo.Update(ctx, vehicle)
	return err
}

func (s *salesService) GetPurchasesByBuyerID(ctx context.Context, buyerID uuid.UUID) ([]model.SalesTransaction, error) {
	return s.salesRepo.FindByBuyerID(ctx, buyerID)
}

func (s *salesService) GetSalesBySellerID(ctx context.Context, sellerID uuid.UUID) ([]model.SalesTransaction, error) {
	return s.salesRepo.FindBySellerID(ctx, sellerID)
}

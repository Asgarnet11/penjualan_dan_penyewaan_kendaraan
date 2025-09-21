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

	newTransaction := model.SalesTransaction{
		ID:          uuid.New(),
		VehicleID:   vehicleID,
		SellerID:    vehicle.OwnerID,
		BuyerID:     buyerID,
		AgreedPrice: vehicle.SalePrice,
		Status:      "payment_pending",
	}

	return s.salesRepo.Create(ctx, newTransaction)
}

func (s *salesService) ConfirmSale(ctx context.Context, transactionID uuid.UUID) error {
	// 1. Update status transaksi penjualan menjadi 'completed'
	// Kita butuh vehicle_id yang dikembalikan dari repository
	transaction, err := s.salesRepo.UpdateStatus(ctx, transactionID, "completed")
	if err != nil {
		return err
	}

	// 2. KRUSIAL: Update status kendaraan menjadi 'sold'
	vehicle, err := s.vehicleRepo.FindByID(ctx, transaction.VehicleID)
	if err != nil {
		return err // Harusnya tidak terjadi jika data konsisten
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

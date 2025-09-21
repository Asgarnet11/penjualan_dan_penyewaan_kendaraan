package repository

import (
	"context"
	"sultra-otomotif-api/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SalesRepository interface {
	Create(ctx context.Context, transaction model.SalesTransaction) (model.SalesTransaction, error)
	UpdateStatus(ctx context.Context, transactionID uuid.UUID, status string) (model.SalesTransaction, error)
	FindByID(ctx context.Context, transactionID uuid.UUID) (model.SalesTransaction, error)
	FindByBuyerID(ctx context.Context, buyerID uuid.UUID) ([]model.SalesTransaction, error)
	FindBySellerID(ctx context.Context, sellerID uuid.UUID) ([]model.SalesTransaction, error)
}

type salesRepository struct {
	db *pgxpool.Pool
}

func NewSalesRepository(db *pgxpool.Pool) SalesRepository {
	return &salesRepository{db: db}
}

func (r *salesRepository) Create(ctx context.Context, t model.SalesTransaction) (model.SalesTransaction, error) {
	query := `INSERT INTO sales_transactions (id, vehicle_id, seller_id, buyer_id, agreed_price, status, payment_token, payment_url)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
              RETURNING created_at, updated_at`

	// Simulasi pembuatan token & url pembayaran
	t.PaymentToken = "SALE-TOKEN-" + t.ID.String()
	t.PaymentURL = "https://example.com/pay/sale/" + t.PaymentToken

	err := r.db.QueryRow(ctx, query, t.ID, t.VehicleID, t.SellerID, t.BuyerID, t.AgreedPrice, t.Status, t.PaymentToken, t.PaymentURL).Scan(&t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return model.SalesTransaction{}, err
	}

	return t, nil
}

func (r *salesRepository) UpdateStatus(ctx context.Context, transactionID uuid.UUID, status string) (model.SalesTransaction, error) {
	query := `UPDATE sales_transactions SET status = $1, updated_at = NOW() WHERE id = $2 RETURNING vehicle_id`

	var updatedTransaction model.SalesTransaction
	updatedTransaction.ID = transactionID
	updatedTransaction.Status = status

	err := r.db.QueryRow(ctx, query, status, transactionID).Scan(&updatedTransaction.VehicleID)
	if err != nil {
		return model.SalesTransaction{}, err
	}
	return updatedTransaction, nil
}

func (r *salesRepository) FindByID(ctx context.Context, transactionID uuid.UUID) (model.SalesTransaction, error) {
	var t model.SalesTransaction
	query := `SELECT id, vehicle_id, seller_id, buyer_id, agreed_price, status, created_at, updated_at FROM sales_transactions WHERE id = $1`
	err := r.db.QueryRow(ctx, query, transactionID).Scan(&t.ID, &t.VehicleID, &t.SellerID, &t.BuyerID, &t.AgreedPrice, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	return t, err
}

func (r *salesRepository) FindByBuyerID(ctx context.Context, buyerID uuid.UUID) ([]model.SalesTransaction, error) {
	var transactions []model.SalesTransaction
	query := `SELECT id, vehicle_id, seller_id, buyer_id, agreed_price, status, created_at, updated_at FROM sales_transactions WHERE buyer_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, buyerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t model.SalesTransaction
		if err := rows.Scan(&t.ID, &t.VehicleID, &t.SellerID, &t.BuyerID, &t.AgreedPrice, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (r *salesRepository) FindBySellerID(ctx context.Context, sellerID uuid.UUID) ([]model.SalesTransaction, error) {
	var transactions []model.SalesTransaction
	query := `SELECT id, vehicle_id, seller_id, buyer_id, agreed_price, status, created_at, updated_at FROM sales_transactions WHERE seller_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, sellerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t model.SalesTransaction
		if err := rows.Scan(&t.ID, &t.VehicleID, &t.SellerID, &t.BuyerID, &t.AgreedPrice, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

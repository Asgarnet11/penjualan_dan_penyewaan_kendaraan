package repository

import (
	"context"
	"sultra-otomotif-api/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository adalah interface yang akan digunakan oleh service
type UserRepository interface {
	Save(ctx context.Context, user model.User) (model.User, error)
	FindByEmail(ctx context.Context, email string) (model.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (model.User, error)
	FindUsersByRole(ctx context.Context, role string) ([]model.User, error)
	UpdateVerificationStatus(ctx context.Context, userID uuid.UUID, status bool) error
	FindAll(ctx context.Context) ([]model.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// userRepository adalah implementasi dari interface di atas
type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository adalah constructor untuk userRepository
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Save(ctx context.Context, user model.User) (model.User, error) {
	query := `INSERT INTO users (id, full_name, email, password_hash, phone_number, role)
              VALUES ($1, $2, $3, $4, $5, $6)
              RETURNING is_verified, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		user.ID,
		user.FullName,
		user.Email,
		user.PasswordHash,
		user.PhoneNumber,
		user.Role,
	).Scan(&user.IsVerified, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	query := `SELECT id, full_name, email, password_hash, phone_number, role, is_verified, verified_at, created_at, updated_at
              FROM users WHERE email = $1`

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.PhoneNumber,
		&user.Role,
		&user.IsVerified,
		&user.VerifiedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return user, err
}

// FindByID mencari satu pengguna berdasarkan ID-nya.
func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	var user model.User
	query := `SELECT id, full_name, email, password_hash, phone_number, role, is_verified, verified_at, created_at, updated_at
              FROM users WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.PhoneNumber,
		&user.Role,
		&user.IsVerified,
		&user.VerifiedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return user, err
}

// FindUsersByRole mencari semua pengguna dengan peran tertentu (misal: 'vendor').
func (r *userRepository) FindUsersByRole(ctx context.Context, role string) ([]model.User, error) {
	var users []model.User
	query := `SELECT id, full_name, email, phone_number, role, is_verified, verified_at, created_at
              FROM users WHERE role = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID,
			&user.FullName,
			&user.Email,
			&user.PhoneNumber,
			&user.Role,
			&user.IsVerified,
			&user.VerifiedAt,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// UpdateVerificationStatus mengubah status verifikasi seorang pengguna.
func (r *userRepository) UpdateVerificationStatus(ctx context.Context, userID uuid.UUID, status bool) error {
	var query string
	var err error

	if status {
		query = `UPDATE users SET is_verified = TRUE, verified_at = NOW() WHERE id = $1`
		_, err = r.db.Exec(ctx, query, userID)
	} else {
		query = `UPDATE users SET is_verified = FALSE, verified_at = NULL WHERE id = $1`
		_, err = r.db.Exec(ctx, query, userID)
	}

	return err
}

func (r *userRepository) FindAll(ctx context.Context) ([]model.User, error) {
	var users []model.User
	query := `SELECT id, full_name, email, phone_number, role, is_verified, verified_at, created_at FROM users ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.FullName, &user.Email, &user.PhoneNumber, &user.Role, &user.IsVerified, &user.VerifiedAt, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

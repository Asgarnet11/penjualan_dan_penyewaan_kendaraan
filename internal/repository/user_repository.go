package repository

import (
	"context"
	"sultra-otomotif-api/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Save(ctx context.Context, user model.User) (model.User, error)
	FindByEmail(ctx context.Context, email string) (model.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Save(ctx context.Context, user model.User) (model.User, error) {
	query := `INSERT INTO users (id, full_name, email, password_hash, phone_number, role)
              VALUES ($1, $2, $3, $4, $5, $6)
              RETURNING created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		user.ID,
		user.FullName,
		user.Email,
		user.PasswordHash,
		user.PhoneNumber,
		user.Role,
	).Scan(&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	query := `SELECT id, full_name, email, password_hash, phone_number, role, created_at, updated_at
              FROM users WHERE email = $1`

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.PhoneNumber,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return user, err
	}
	return user, nil
}

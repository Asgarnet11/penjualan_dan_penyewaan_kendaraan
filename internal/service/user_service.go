package service

import (
	"context"
	"errors"
	"sultra-otomotif-api/internal/auth"
	"sultra-otomotif-api/internal/helper"
	"sultra-otomotif-api/internal/model"
	"sultra-otomotif-api/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserService interface {
	RegisterUser(ctx context.Context, input model.RegisterUserInput) (model.User, error)
	LoginUser(ctx context.Context, input model.LoginUserInput, jwtSecret string) (string, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) RegisterUser(ctx context.Context, input model.RegisterUserInput) (model.User, error) {
	// Cek apakah email sudah ada
	_, err := s.repo.FindByEmail(ctx, input.Email)
	if err == nil { // Jika tidak ada error, berarti user ditemukan
		return model.User{}, errors.New("email already registered")
	}
	if !errors.Is(err, pgx.ErrNoRows) { // Jika errornya bukan karena tidak ada baris
		return model.User{}, err
	}

	passwordHash, err := helper.HashPassword(input.Password)
	if err != nil {
		return model.User{}, err
	}

	newUser := model.User{
		ID:           uuid.New(),
		FullName:     input.FullName,
		Email:        input.Email,
		PasswordHash: passwordHash,
		PhoneNumber:  input.PhoneNumber,
		Role:         input.Role,
	}

	createdUser, err := s.repo.Save(ctx, newUser)
	if err != nil {
		return model.User{}, err
	}

	return createdUser, nil
}

func (s *userService) LoginUser(ctx context.Context, input model.LoginUserInput, jwtSecret string) (string, error) {
	user, err := s.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("invalid email or password")
		}
		return "", err
	}

	isValidPassword := helper.CheckPasswordHash(input.Password, user.PasswordHash)
	if !isValidPassword {
		return "", errors.New("invalid email or password")
	}

	token, err := auth.GenerateToken(user.ID, user.Role, jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}

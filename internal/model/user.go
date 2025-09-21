package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `json:"id"`
	FullName     string     `json:"full_name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	PhoneNumber  string     `json:"phone_number"`
	Role         string     `json:"role"`
	IsVerified   bool       `json:"is_verified"`
	VerifiedAt   *time.Time `json:"verified_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type RegisterUserInput struct {
	FullName    string `json:"full_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Role        string `json:"role" binding:"required,oneof=customer vendor"`
}

type LoginUserInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

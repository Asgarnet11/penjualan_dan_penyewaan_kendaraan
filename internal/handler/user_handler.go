package handler

import (
	"net/http"
	"sultra-otomotif-api/internal/helper"
	"sultra-otomotif-api/internal/model"
	"sultra-otomotif-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService service.UserService
	jwtSecret   string
}

func NewUserHandler(userService service.UserService, jwtSecret string) *UserHandler {
	return &UserHandler{userService: userService, jwtSecret: jwtSecret}
}

func (h *UserHandler) Register(ctx *gin.Context) {
	var input model.RegisterUserInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(ctx, "Invalid input data", http.StatusBadRequest, err)
		return
	}

	user, err := h.userService.RegisterUser(ctx, input)
	if err != nil {
		helper.ErrorResponse(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	// Buat response tanpa password hash
	userResponse := model.User{
		ID:          user.ID,
		FullName:    user.FullName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	helper.APIResponse(ctx, "User registered successfully", http.StatusCreated, userResponse)
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var input model.LoginUserInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(ctx, "Invalid input data", http.StatusBadRequest, err)
		return
	}

	token, err := h.userService.LoginUser(ctx, input, h.jwtSecret)
	if err != nil {
		helper.ErrorResponse(ctx, err.Error(), http.StatusUnauthorized, err)
		return
	}

	tokenResponse := model.AuthResponse{Token: token}
	helper.APIResponse(ctx, "Login successful", http.StatusOK, tokenResponse)
}

func (h *UserHandler) GetMe(ctx *gin.Context) {
	currentUserID := ctx.MustGet("currentUserID").(uuid.UUID)

	user, err := h.userService.GetUserByID(ctx, currentUserID)
	if err != nil {
		helper.ErrorResponse(ctx, "Failed to get user profile", http.StatusNotFound, err)
		return
	}
	helper.APIResponse(ctx, "Successfully fetched user profile", http.StatusOK, user)
}

package handler

import (
	"net/http"
	"strings"
	"sultra-otomotif-api/internal/helper"
	"sultra-otomotif-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatHandler struct {
	chatService service.ChatService
}

func NewChatHandler(chatService service.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

func (h *ChatHandler) StartConversation(ctx *gin.Context) {
	vehicleID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helper.ErrorResponse(ctx, "Invalid vehicle ID", http.StatusBadRequest, err)
		return
	}

	customerID := ctx.MustGet("currentUserID").(uuid.UUID)

	conversation, err := h.chatService.StartConversation(ctx, customerID, vehicleID)
	if err != nil {
		helper.ErrorResponse(ctx, err.Error(), http.StatusConflict, err)
		return
	}

	helper.APIResponse(ctx, "Conversation started successfully", http.StatusCreated, conversation)
}

func (h *ChatHandler) ListConversations(ctx *gin.Context) {
	userID := ctx.MustGet("currentUserID").(uuid.UUID)
	conversations, err := h.chatService.GetConversationsForUser(ctx, userID)
	if err != nil {
		helper.ErrorResponse(ctx, "Failed to fetch conversations", http.StatusInternalServerError, err)
		return
	}
	helper.APIResponse(ctx, "Successfully fetched user conversations", http.StatusOK, conversations)
}

func (h *ChatHandler) GetMessages(ctx *gin.Context) {
	conversationID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helper.ErrorResponse(ctx, "Invalid conversation ID", http.StatusBadRequest, err)
		return
	}

	userID := ctx.MustGet("currentUserID").(uuid.UUID)
	messages, err := h.chatService.GetMessagesForConversation(ctx, conversationID, userID)
	if err != nil {
		if strings.HasPrefix(err.Error(), "forbidden") {
			helper.ErrorResponse(ctx, err.Error(), http.StatusForbidden, err)
		} else {
			helper.ErrorResponse(ctx, err.Error(), http.StatusNotFound, err)
		}
		return
	}
	helper.APIResponse(ctx, "Successfully fetched messages", http.StatusOK, messages)
}

package service

import (
	"context"
	"errors"
	"sultra-otomotif-api/internal/model"
	"sultra-otomotif-api/internal/repository"

	"github.com/google/uuid"
)

type ChatService interface {
	SaveMessage(ctx context.Context, msg model.Message) (model.Message, error)
	StartConversation(ctx context.Context, customerID, vehicleID uuid.UUID) (model.Conversation, error)
	GetConversationsForUser(ctx context.Context, userID uuid.UUID) ([]model.Conversation, error)
	GetMessagesForConversation(ctx context.Context, conversationID, userID uuid.UUID) ([]model.Message, error)
}

type chatService struct {
	chatRepo    repository.ChatRepository
	vehicleRepo repository.VehicleRepository
}

func NewChatService(chatRepo repository.ChatRepository, vehicleRepo repository.VehicleRepository) ChatService {
	return &chatService{chatRepo: chatRepo, vehicleRepo: vehicleRepo}
}

func (s *chatService) SaveMessage(ctx context.Context, msg model.Message) (model.Message, error) {
	return s.chatRepo.SaveMessage(ctx, msg)
}

func (s *chatService) StartConversation(ctx context.Context, customerID, vehicleID uuid.UUID) (model.Conversation, error) {
	vehicle, err := s.vehicleRepo.FindByID(ctx, vehicleID)
	if err != nil {
		return model.Conversation{}, errors.New("vehicle not found")
	}

	vendorID := vehicle.OwnerID
	if customerID == vendorID {
		return model.Conversation{}, errors.New("cannot start conversation with yourself")
	}

	return s.chatRepo.FindOrCreateConversation(ctx, customerID, vendorID, vehicleID)
}

func (s *chatService) GetConversationsForUser(ctx context.Context, userID uuid.UUID) ([]model.Conversation, error) {
	return s.chatRepo.FindConversationsByUserID(ctx, userID)
}

func (s *chatService) GetMessagesForConversation(ctx context.Context, conversationID, userID uuid.UUID) ([]model.Message, error) {
	// Validasi keamanan: pastikan user yang meminta adalah bagian dari percakapan
	convo, err := s.chatRepo.FindConversationByID(ctx, conversationID)
	if err != nil {
		return nil, errors.New("conversation not found")
	}

	if userID != convo.CustomerID && userID != convo.VendorID {
		return nil, errors.New("forbidden: you are not a participant in this conversation")
	}

	return s.chatRepo.FindMessagesByConversationID(ctx, conversationID)
}

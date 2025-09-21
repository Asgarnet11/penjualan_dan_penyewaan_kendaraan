package service

import (
	"context"
	"sultra-otomotif-api/internal/model"
	"sultra-otomotif-api/internal/repository"
)

type ChatService interface {
	SaveMessage(ctx context.Context, msg model.Message) (model.Message, error)
}

type chatService struct {
	chatRepo repository.ChatRepository
}

func NewChatService(chatRepo repository.ChatRepository) ChatService {
	return &chatService{chatRepo: chatRepo}
}

func (s *chatService) SaveMessage(ctx context.Context, msg model.Message) (model.Message, error) {
	// Di sini bisa ditambahkan validasi, misal cek apakah pengirim adalah bagian dari percakapan
	return s.chatRepo.SaveMessage(ctx, msg)
}

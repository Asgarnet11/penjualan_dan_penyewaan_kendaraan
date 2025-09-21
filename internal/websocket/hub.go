package websocket

import (
	"context"
	"encoding/json"
	"log"
	"sultra-otomotif-api/internal/model"
	"sultra-otomotif-api/internal/service"

	"github.com/google/uuid"
)

type Hub struct {
	Clients     map[uuid.UUID]*Client
	Broadcast   chan *model.Message
	Register    chan *Client
	Unregister  chan *Client
	chatService service.ChatService
}

func NewHub(chatService service.ChatService) *Hub {
	return &Hub{
		Clients:     make(map[uuid.UUID]*Client),
		Broadcast:   make(chan *model.Message),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		chatService: chatService,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.ID] = client // <-- Simpan client dengan key UserID
		case client := <-h.Unregister:
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID) // <-- Hapus client dengan key UserID
				close(client.Send)
			}
		case message := <-h.Broadcast:
			// 1. Simpan pesan ke database
			savedMessage, err := h.chatService.SaveMessage(context.Background(), *message)
			if err != nil {
				log.Printf("error saving message: %v", err)
				continue
			}

			// 2. Cek apakah penerima sedang online
			recipientClient, ok := h.Clients[savedMessage.RecipientID]
			if ok {
				// 3. Jika online, kirim pesan ke penerima
				jsonMessage, _ := json.Marshal(savedMessage)
				select {
				case recipientClient.Send <- jsonMessage:
				default:
					close(recipientClient.Send)
					delete(h.Clients, recipientClient.ID)
				}
			}
		}
	}
}

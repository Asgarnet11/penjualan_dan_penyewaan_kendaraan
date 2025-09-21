package websocket

import "github.com/google/uuid"

type Message struct {
	Content        string    `json:"content"`
	RecipientID    uuid.UUID `json:"recipient_id"`
	ConversationID uuid.UUID `json:"conversation_id"`
}

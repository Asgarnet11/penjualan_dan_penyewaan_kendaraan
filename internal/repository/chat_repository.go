package repository

import (
	"context"
	"sultra-otomotif-api/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatRepository interface {
	SaveMessage(ctx context.Context, msg model.Message) (model.Message, error)
	// Fungsi lain seperti GetConversations, GetMessagesByConversationID akan ditambahkan nanti
}

type chatRepository struct {
	db *pgxpool.Pool
}

func NewChatRepository(db *pgxpool.Pool) ChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) SaveMessage(ctx context.Context, msg model.Message) (model.Message, error) {
	query := `INSERT INTO messages (id, conversation_id, sender_id, recipient_id, content)
			  VALUES (uuid_generate_v4(), $1, $2, $3, $4)
			  RETURNING id, created_at`

	err := r.db.QueryRow(ctx, query, msg.ConversationID, msg.SenderID, msg.RecipientID, msg.Content).Scan(&msg.ID, &msg.CreatedAt)
	if err != nil {
		return model.Message{}, err
	}

	// Update timestamp 'updated_at' di tabel conversations
	_, err = r.db.Exec(ctx, `UPDATE conversations SET updated_at = NOW() WHERE id = $1`, msg.ConversationID)

	return msg, err
}

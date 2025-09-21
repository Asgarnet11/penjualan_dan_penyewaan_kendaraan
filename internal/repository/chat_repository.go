package repository

import (
	"context"
	"sultra-otomotif-api/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatRepository interface {
	SaveMessage(ctx context.Context, msg model.Message) (model.Message, error)
	FindOrCreateConversation(ctx context.Context, customerID, vendorID, vehicleID uuid.UUID) (model.Conversation, error)
	FindConversationByID(ctx context.Context, conversationID uuid.UUID) (model.Conversation, error)
	FindConversationsByUserID(ctx context.Context, userID uuid.UUID) ([]model.Conversation, error)
	FindMessagesByConversationID(ctx context.Context, conversationID uuid.UUID) ([]model.Message, error)
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

	_, err = r.db.Exec(ctx, `UPDATE conversations SET updated_at = NOW() WHERE id = $1`, msg.ConversationID)

	return msg, err
}

func (r *chatRepository) FindOrCreateConversation(ctx context.Context, customerID, vendorID, vehicleID uuid.UUID) (model.Conversation, error) {
	var convo model.Conversation
	querySelect := `SELECT id, customer_id, vendor_id, vehicle_id, created_at, updated_at FROM conversations WHERE customer_id = $1 AND vendor_id = $2 AND vehicle_id = $3`
	err := r.db.QueryRow(ctx, querySelect, customerID, vendorID, vehicleID).Scan(&convo.ID, &convo.CustomerID, &convo.VendorID, &convo.VehicleID, &convo.CreatedAt, &convo.UpdatedAt)

	if err == nil {
		return convo, nil
	}

	if err == pgx.ErrNoRows {
		queryInsert := `INSERT INTO conversations (id, customer_id, vendor_id, vehicle_id) VALUES (uuid_generate_v4(), $1, $2, $3) RETURNING id, customer_id, vendor_id, vehicle_id, created_at, updated_at`
		errInsert := r.db.QueryRow(ctx, queryInsert, customerID, vendorID, vehicleID).Scan(&convo.ID, &convo.CustomerID, &convo.VendorID, &convo.VehicleID, &convo.CreatedAt, &convo.UpdatedAt)
		return convo, errInsert
	}

	return model.Conversation{}, err
}

func (r *chatRepository) FindConversationByID(ctx context.Context, conversationID uuid.UUID) (model.Conversation, error) {
	var convo model.Conversation
	query := `SELECT id, customer_id, vendor_id, vehicle_id, created_at, updated_at FROM conversations WHERE id = $1`
	err := r.db.QueryRow(ctx, query, conversationID).Scan(&convo.ID, &convo.CustomerID, &convo.VendorID, &convo.VehicleID, &convo.CreatedAt, &convo.UpdatedAt)
	return convo, err
}

func (r *chatRepository) FindConversationsByUserID(ctx context.Context, userID uuid.UUID) ([]model.Conversation, error) {
	var conversations []model.Conversation
	query := `SELECT id, customer_id, vendor_id, vehicle_id, created_at, updated_at FROM conversations WHERE customer_id = $1 OR vendor_id = $1 ORDER BY updated_at DESC`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c model.Conversation
		if err := rows.Scan(&c.ID, &c.CustomerID, &c.VendorID, &c.VehicleID, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		conversations = append(conversations, c)
	}
	return conversations, nil
}

func (r *chatRepository) FindMessagesByConversationID(ctx context.Context, conversationID uuid.UUID) ([]model.Message, error) {
	var messages []model.Message
	query := `SELECT id, conversation_id, sender_id, recipient_id, content, is_read, created_at FROM messages WHERE conversation_id = $1 ORDER BY created_at ASC`
	rows, err := r.db.Query(ctx, query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m model.Message
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.RecipientID, &m.Content, &m.IsRead, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

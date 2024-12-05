package services

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID `json:"id"`
	ChannelID uuid.UUID `json:"channel_id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MessageService struct {
	db *sql.DB
}

func NewMessageService(db *sql.DB) *MessageService {
	return &MessageService{db: db}
}

func (s *MessageService) CreateMessage(channelID, userID uuid.UUID, content string) (*Message, error) {
	if content == "" {
		return nil, fmt.Errorf("message content cannot be empty")
	}

	var msg Message
	err := s.db.QueryRow(`
        INSERT INTO messages (channel_id, user_id, content)
        VALUES ($1, $2, $3)
        RETURNING id, channel_id, user_id, content, created_at, updated_at
    `, channelID, userID, content).Scan(
		&msg.ID,
		&msg.ChannelID,
		&msg.UserID,
		&msg.Content,
		&msg.CreatedAt,
		&msg.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (s *MessageService) GetChannelMessages(channelID uuid.UUID, limit int, before time.Time) ([]Message, error) {
	if limit <= 0 {
		limit = 50
	}

	log.Printf("Query params: channelID=%v, limit=%d, before=%v", channelID, limit, before)

	rows, err := s.db.Query(`
        SELECT id, channel_id, user_id, content, created_at, updated_at
        FROM messages
        WHERE channel_id = $1 AND created_at <= $2
        ORDER BY created_at DESC
        LIMIT $3
    `, channelID, before, limit)

	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(
			&msg.ID,
			&msg.ChannelID,
			&msg.UserID,
			&msg.Content,
			&msg.CreatedAt,
			&msg.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		messages = append(messages, msg)
	}

	log.Printf("Found %d messages", len(messages))
	return messages, rows.Err()
}

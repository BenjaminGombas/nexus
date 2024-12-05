package services

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type Channel struct {
	ID        uuid.UUID `json:"id"`
	ServerID  uuid.UUID `json:"server_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type ChannelService struct {
	db *sql.DB
}

func NewChannelService(db *sql.DB) *ChannelService {
	return &ChannelService{db: db}
}

func (s *ChannelService) CreateChannel(serverID uuid.UUID, name string, channelType string) (*Channel, error) {
	if name == "" {
		return nil, fmt.Errorf("channel name cannot be empty")
	}

	validTypes := map[string]bool{
		"text":  true,
		"voice": true,
	}
	if !validTypes[channelType] {
		return nil, fmt.Errorf("invalid channel type: must be 'text' or 'voice'")
	}

	var channel Channel
	err := s.db.QueryRow(`
        INSERT INTO channels (server_id, name, type)
        VALUES ($1, $2, $3)
        RETURNING id, server_id, name, type, created_at, updated_at
    `, serverID, name, channelType).Scan(
		&channel.ID,
		&channel.ServerID,
		&channel.Name,
		&channel.Type,
		&channel.CreatedAt,
		&channel.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &channel, nil
}

func (s *ChannelService) GetServerChannels(serverID uuid.UUID) ([]Channel, error) {
	rows, err := s.db.Query(`
		SELECT id, server_id, name, type, created_at, updated_at 
		FROM channels
		WHERE server_id = $1
		ORDER BY created_at ASC
	`, serverID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []Channel
	for rows.Next() {
		var channel Channel
		err := rows.Scan(
			&channel.ID,
			&channel.ServerID,
			&channel.Name,
			&channel.Type,
			&channel.CreatedAt,
			&channel.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		channels = append(channels, channel)
	}

	return channels, nil
}

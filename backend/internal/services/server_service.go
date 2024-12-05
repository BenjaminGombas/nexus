package services

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type Server struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	OwnerID   uuid.UUID `json:"owner_id"`
	IconURL   *string   `json:"icon_url,omitempty"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type ServerService struct {
	db *sql.DB
}

func NewServerService(db *sql.DB) *ServerService {
	return &ServerService{db: db}
}

func (s *ServerService) CreateServer(name string, ownerID uuid.UUID) (*Server, error) {
	if name == "" {
		return nil, fmt.Errorf("server name cannot be empty")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	var server Server
	err = tx.QueryRow(`
        INSERT INTO servers (name, owner_id)
        VALUES ($1, $2)
        RETURNING id, name, owner_id, icon_url, created_at, updated_at
    `, name, ownerID).Scan(
		&server.ID,
		&server.Name,
		&server.OwnerID,
		&server.IconURL,
		&server.CreatedAt,
		&server.UpdatedAt,
	)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Add owner to server_members
	_, err = tx.Exec(`
        INSERT INTO server_members (server_id, user_id)
        VALUES ($1, $2)
    `, server.ID, ownerID)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &server, tx.Commit()
}

func (s *ServerService) GetUserServers(userID uuid.UUID) ([]Server, error) {
	log.Printf("Getting servers for user: %v", userID)
	rows, err := s.db.Query(`
        SELECT s.id, s.name, s.owner_id, s.icon_url, s.created_at, s.updated_at
        FROM servers s
        JOIN server_members sm ON s.id = sm.server_id
        WHERE sm.user_id = $1
    `, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servers []Server
	for rows.Next() {
		var server Server
		err := rows.Scan(
			&server.ID,
			&server.Name,
			&server.OwnerID,
			&server.IconURL,
			&server.CreatedAt,
			&server.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		servers = append(servers, server)
	}

	log.Printf("Found %d servers", len(servers))
	return servers, rows.Err()
}

func (s *ServerService) JoinServer(serverID, userID uuid.UUID) error {
	_, err := s.db.Exec(`
		INSERT INTO server_members (server_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, serverID, userID)

	return err
}

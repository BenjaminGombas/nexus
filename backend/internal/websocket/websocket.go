package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"nexus/internal/services"

	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

type Conn interface {
	WriteJSON(v interface{}) error
	Close() error
}

type MessageServiceInterface interface {
	CreateMessage(channelID, userID uuid.UUID, content string) (*services.Message, error)
	GetChannelMessages(channelID uuid.UUID, limit int, before time.Time) ([]services.Message, error)
}

type Client struct {
	ID       string
	UserID   uuid.UUID
	Conn     Conn
	Hub      *Hub
	mu       sync.Mutex
	Channels map[uuid.UUID]bool
}

type Message struct {
	Type      string         `json:"type"`
	ChannelID uuid.UUID      `json:"channel_id,omitempty"`
	UserID    uuid.UUID      `json:"user_id,omitempty"`
	Content   string         `json:"content,omitempty"`
	Data      map[string]any `json:"data,omitempty"`
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	msgService MessageServiceInterface
	mu         sync.RWMutex
}

func NewHub(msgService MessageServiceInterface) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		msgService: msgService,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Conn.Close()
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			if message.Type == "message" {
				_, err := h.msgService.CreateMessage(message.ChannelID, message.UserID, message.Content)
				if err != nil {
					log.Printf("error persisting message: %v", err)
					continue
				}
			}

			h.mu.RLock()
			for client := range h.clients {
				if message.ChannelID != uuid.Nil {
					if !client.Channels[message.ChannelID] {
						continue
					}
				}

				client.mu.Lock()
				err := client.Conn.WriteJSON(message)
				client.mu.Unlock()

				if err != nil {
					log.Printf("error broadcasting message to client: %v", err)
					client.Conn.Close()
					h.unregister <- client
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (c *Client) SubscribeToChannel(channelID uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Channels[channelID] = true
}

func (c *Client) UnsubscribeFromChannel(channelID uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.Channels, channelID)
}

func HandleConnection(c *websocket.Conn, hub *Hub) {
	userID := c.Params("userID")
	if userID == "" {
		return
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return
	}

	client := &Client{
		ID:       uuid.New().String(),
		UserID:   uid,
		Conn:     c,
		Hub:      hub,
		Channels: make(map[uuid.UUID]bool),
	}

	hub.register <- client
	defer func() {
		hub.unregister <- client
	}()

	for {
		messageType, msg, err := c.ReadMessage()
		if err != nil {
			break
		}

		if messageType == websocket.TextMessage {
			var message Message
			if err := json.Unmarshal(msg, &message); err != nil {
				continue
			}

			message.UserID = client.UserID
			hub.broadcast <- message
		}
	}
}

package handlers

import (
	"time"

	"nexus/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MessageHandler struct {
	messageService *services.MessageService
}

func NewMessageHandler(messageService *services.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

// GetChannelMessages handles HTTP requests to get message history
func (h *MessageHandler) GetChannelMessages(c *fiber.Ctx) error {
	channelID, err := uuid.Parse(c.Params("channelId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid channel ID"})
	}

	limit := 50 // Default limit
	if c.Query("limit") != "" {
		// Parse limit from query if provided
	}

	messages, err := h.messageService.GetChannelMessages(channelID, limit, time.Now())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch messages"})
	}

	return c.JSON(messages)
}

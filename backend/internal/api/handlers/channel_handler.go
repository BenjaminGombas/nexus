package handlers

import (
	"nexus/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ChannelHandler struct {
	channelService *services.ChannelService
}

func NewChannelHandler(channelService *services.ChannelService) *ChannelHandler {
	return &ChannelHandler{
		channelService: channelService,
	}
}

func (h *ChannelHandler) CreateChannel(c *fiber.Ctx) error {
	var input struct {
		ServerID uuid.UUID `json:"server_id"`
		Name     string    `json:"name"`
		Type     string    `json:"type"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	channel, err := h.channelService.CreateChannel(input.ServerID, input.Name, input.Type)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create channel"})
	}

	return c.Status(201).JSON(channel)
}

func (h *ChannelHandler) GetServerChannels(c *fiber.Ctx) error {
	serverID, err := uuid.Parse(c.Params("serverId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid server ID"})
	}

	channels, err := h.channelService.GetServerChannels(serverID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not fetch channels"})
	}

	return c.JSON(channels)
}

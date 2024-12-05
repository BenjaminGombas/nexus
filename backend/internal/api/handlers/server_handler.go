package handlers

import (
	"nexus/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ServerHandler struct {
	serverService *services.ServerService
}

func NewServerHandler(serverService *services.ServerService) *ServerHandler {
	return &ServerHandler{
		serverService: serverService,
	}
}

func (h *ServerHandler) CreateServer(c *fiber.Ctx) error {
	var input struct {
		Name string `json:"name"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(uuid.UUID)

	server, err := h.serverService.CreateServer(input.Name, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create server"})
	}

	return c.Status(201).JSON(server)
}

func (h *ServerHandler) GetUserServers(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	servers, err := h.serverService.GetUserServers(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not fetch servers"})
	}

	return c.JSON(servers)
}

func (h *ServerHandler) JoinServer(c *fiber.Ctx) error {
	serverID, err := uuid.Parse(c.Params("serverId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid server ID"})
	}

	userID := c.Locals("userID").(uuid.UUID)

	if err := h.serverService.JoinServer(serverID, userID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not join server"})
	}

	return c.SendStatus(200)
}

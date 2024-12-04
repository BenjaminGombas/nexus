package handlers

import (
	"database/sql"
	"nexus/internal/auth"
	"nexus/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	db *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	input := new(models.CreateUserInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Check if email already exists
	var exists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", input.Email).Scan(&exists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	if exists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email already exists",
		})
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error processing password",
		})
	}

	// Create user
	userID := uuid.New()
	_, err = h.db.Exec(`
        INSERT INTO users (id, username, email, password_hash, created_at, updated_at)
        VALUES ($1, $2, $3, $4, NOW(), NOW())`,
		userID, input.Username, input.Email, hashedPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating user",
		})
	}

	// Generate JWT
	token, err := auth.GenerateJWT(userID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error generating token",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":       userID,
			"username": input.Username,
			"email":    input.Email,
		},
	})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	input := new(models.LoginInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Get user from database
	var user models.User
	var hashedPassword string
	err := h.db.QueryRow(`
        SELECT id, username, email, password_hash
        FROM users WHERE email = $1`,
		input.Email,
	).Scan(&user.ID, &user.Username, &user.Email, &hashedPassword)

	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Check password
	if !auth.CheckPasswordHash(input.Password, hashedPassword) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Generate JWT
	token, err := auth.GenerateJWT(user.ID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error generating token",
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"

	"nexus/internal/api/handlers"
	"nexus/internal/auth"
	"nexus/internal/database"
	"nexus/internal/services"
	wsHandler "nexus/internal/websocket"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize database connection
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Could not initialize database connection: %v", err)
	}
	defer db.Close()

	// Initialize services
	userService := services.NewUserService(db)
	messageService := services.NewMessageService(db)
	serverService := services.NewServerService(db)
	channelService := services.NewChannelService(db)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	messageHandler := handlers.NewMessageHandler(messageService)
	serverHandler := handlers.NewServerHandler(serverService)
	channelHandler := handlers.NewChannelHandler(channelService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Nexus API v1",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, HEAD, PUT, DELETE, PATCH",
	}))

	// Initialize WebSocket hub
	hub := wsHandler.NewHub(messageService)
	go hub.Run()

	// API routes
	api := app.Group("/api/v1")

	// Public routes
	authGroup := api.Group("/auth")
	authGroup.Post("/register", userHandler.Register)
	authGroup.Post("/login", userHandler.Login)

	// Protected routes
	api.Use(auth.Required())

	// Server routes
	servers := api.Group("/servers")
	servers.Post("/", serverHandler.CreateServer)
	servers.Get("/", serverHandler.GetUserServers)
	servers.Post("/:serverId/join", serverHandler.JoinServer)

	// Channel routes
	channels := api.Group("/channels")
	channels.Post("/", channelHandler.CreateChannel)
	channels.Get("/server/:serverId", channelHandler.GetServerChannels)

	// Message routes
	messages := api.Group("/messages")
	messages.Get("/channel/:channelId", messageHandler.GetChannelMessages)

	// WebSocket setup
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		if !auth.WSRequired()(c) {
			c.Close()
			return
		}
		wsHandler.HandleConnection(c, hub)
	}))

	// Get port from environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

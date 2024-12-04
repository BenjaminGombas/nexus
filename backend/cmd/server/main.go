package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"nexus/internal/api/handlers"
	"nexus/internal/database"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize database
	db, err := database.NewConnection()
	if err != nil {
		log.Fatal("Could not initialize database:", err)
	}
	defer db.Close()

	// Create new Fiber app
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

	// Initialize handlers
	userHandler := handlers.NewUserHandler(db)

	// Routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Auth routes
	auth := v1.Group("/auth")
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)

	// Get port from env or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

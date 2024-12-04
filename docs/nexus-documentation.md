# Nexus Project Documentation - Day 1

## Project Diary

### Initial Setup and Planning
Today I began the Nexus project, a Discord-inspired chat platform. The first major task was creating a comprehensive project plan, which included:
- Core features definition
- System architecture design
- Database schema planning
- Project structure layout
- Technical stack decisions
- Development phase planning

### Project Structure Creation
I created the basic project structure using a batch script. The structure follows a clean architecture pattern:
```
/nexus
├── /frontend           # React frontend (to be implemented)
├── /backend           # Go backend
│   ├── /cmd          # Main applications
│   ├── /internal     # Private packages
│   ├── /pkg          # Public packages
│   └── /config       # Configuration files
└── /docs             # Documentation
```

### Development Environment Setup
1. Initialized Go module:
```cmd
cd backend
go mod init nexus
```

2. Installed core dependencies:
```cmd
go get github.com/gofiber/fiber/v2
go get github.com/joho/godotenv
go get github.com/lib/pq
go get github.com/golang-jwt/jwt/v4
go get github.com/google/uuid
go get golang.org/x/crypto/bcrypt
```

3. Set up PostgreSQL using Docker:
```cmd
docker-compose up -d
```

## Code Documentation

### Main Server (cmd/server/main.go)
The entry point of our application.

**Key Components:**
1. Environment Loading
```go
if err := godotenv.Load(); err != nil {
    log.Printf("Warning: .env file not found")
}
```
Purpose: Loads environment variables from .env file

2. Fiber App Initialization
```go
app := fiber.New(fiber.Config{
    AppName: "Nexus Chat API v1",
})
```
Purpose: Creates the main web server instance

3. Middleware Setup
```go
app.Use(logger.New())
app.Use(cors.New(cors.Config{
    AllowOrigins: "*",
    AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    AllowMethods: "GET, POST, HEAD, PUT, DELETE, PATCH",
}))
```
Purpose: Configures logging and CORS policies

### Database Connection (internal/database/database.go)
Manages PostgreSQL database connections.

**Key Function:**
```go
func NewConnection() (*sql.DB, error)
```
Purpose: Creates and tests a new database connection
Returns: SQL database connection pool and error if any
Usage: Called during server initialization

### User Model (internal/models/user.go)
Defines the user data structure and related types.

**Structures:**
1. User
```go
type User struct {
    ID          uuid.UUID `json:"id"`
    Username    string    `json:"username"`
    Email       string    `json:"email"`
    Password    string    `json:"-"`
    AvatarURL   string    `json:"avatar_url,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```
Purpose: Represents a user in the system
Note: Password field is marked with `json:"-"` to prevent it from being included in JSON responses

2. CreateUserInput
```go
type CreateUserInput struct {
    Username string `json:"username" validate:"required,min=3,max=32"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}
```
Purpose: Validates user registration input

### Authentication (internal/auth/auth.go)
Handles JWT token generation and password hashing.

**Key Functions:**
1. GenerateJWT
```go
func GenerateJWT(userID string) (string, error)
```
Purpose: Creates a JWT token for authenticated users
Parameters: User ID
Returns: JWT token string and error if any

2. HashPassword
```go
func HashPassword(password string) (string, error)
```
Purpose: Securely hashes user passwords
Parameters: Plain text password
Returns: Hashed password and error if any

### User Handler (internal/api/handlers/user_handler.go)
Manages user-related HTTP endpoints.

**Key Methods:**
1. Register
```go
func (h *UserHandler) Register(c *fiber.Ctx) error
```
Purpose: Handles user registration
Endpoint: POST /api/v1/auth/register
Input: CreateUserInput JSON
Output: User data and JWT token

2. Login
```go
func (h *UserHandler) Login(c *fiber.Ctx) error
```
Purpose: Handles user authentication
Endpoint: POST /api/v1/auth/login
Input: LoginInput JSON
Output: User data and JWT token

## Database Schema
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(32) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    avatar_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

## API Endpoints Documentation

### Authentication Endpoints

#### Register User
- URL: `/api/v1/auth/register`
- Method: `POST`
- Body:
```json
{
    "username": "string",
    "email": "string",
    "password": "string"
}
```
- Success Response (201):
```json
{
    "token": "jwt_token",
    "user": {
        "id": "uuid",
        "username": "string",
        "email": "string"
    }
}
```

#### Login User
- URL: `/api/v1/auth/login`
- Method: `POST`
- Body:
```json
{
    "email": "string",
    "password": "string"
}
```
- Success Response (200):
```json
{
    "token": "jwt_token",
    "user": {
        "id": "uuid",
        "username": "string",
        "email": "string"
    }
}
```

## Progress Summary
- ✓ Completed project planning
- ✓ Set up basic project structure
- ✓ Implemented database connection
- ✓ Created user model and authentication
- ✓ Set up basic API endpoints
- ✓ Added environment configuration
- ✓ Implemented error handling

## Next Steps
Moving into Phase 2, I will:
1. Implement WebSocket server
2. Set up message handling
3. Create server and channel management
4. Implement real-time message broadcasting

## Notes and Observations
- Chose Fiber over Echo for better performance and simpler API
- Using UUID for IDs to ensure uniqueness across distributed systems
- Implemented CORS middleware for future frontend integration
- Added basic logging for debugging purposes
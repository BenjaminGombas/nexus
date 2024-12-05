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

# Nexus Project Documentation - Day 2

### WebSocket Implementation and Authentication
Today I focused on implementing real-time communication features and securing both HTTP and WebSocket endpoints:
- Implemented WebSocket server with Hub pattern
- Added message persistence
- Created server and channel management
- Added authentication for both HTTP and WebSocket endpoints
- Added comprehensive tests for WebSocket functionality

### Major Components Implemented

#### WebSocket Infrastructure
1. Implemented Hub for managing connections:
```go
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan Message
    register   chan *Client
    unregister chan *Client
    msgService MessageServiceInterface
    mu         sync.RWMutex
}
```

2. Client management for individual connections:
```go
type Client struct {
    ID       string
    UserID   uuid.UUID
    Conn     Conn
    Hub      *Hub
    mu       sync.Mutex
    Channels map[uuid.UUID]bool
}
```

#### Authentication System
1. JWT-based authentication for HTTP endpoints:
```go
func Required() fiber.Handler {
    return func(c *fiber.Ctx) error {
        auth := c.Get("Authorization")
        // ... JWT validation logic
    }
}
```

2. WebSocket-specific authentication:
```go
func WSRequired() func(*websocket.Conn) bool {
    return func(conn *websocket.Conn) bool {
        token := conn.Query("token")
        // ... WebSocket auth logic
    }
}
```

#### Service Layer Implementation
Created services for core functionality:
- UserService: User management
- ServerService: Server creation and management
- ChannelService: Channel operations
- MessageService: Message persistence and retrieval

## API Documentation

### Server Management Endpoints

#### Create Server
- URL: `/api/v1/servers`
- Method: `POST`
- Auth Required: Yes
- Body:
```json
{
    "name": "string"
}
```
- Success Response (201):
```json
{
    "id": "uuid",
    "name": "string",
    "owner_id": "uuid",
    "created_at": "timestamp"
}
```

#### Get User's Servers
- URL: `/api/v1/servers`
- Method: `GET`
- Auth Required: Yes
- Success Response (200):
```json
[
    {
        "id": "uuid",
        "name": "string",
        "owner_id": "uuid",
        "created_at": "timestamp"
    }
]
```

### Channel Management Endpoints

#### Create Channel
- URL: `/api/v1/channels`
- Method: `POST`
- Auth Required: Yes
- Body:
```json
{
    "server_id": "uuid",
    "name": "string",
    "type": "text|voice"
}
```
- Success Response (201):
```json
{
    "id": "uuid",
    "server_id": "uuid",
    "name": "string",
    "type": "string",
    "created_at": "timestamp"
}
```

#### Get Server Channels
- URL: `/api/v1/channels/server/:serverId`
- Method: `GET`
- Auth Required: Yes
- Success Response (200):
```json
[
    {
        "id": "uuid",
        "server_id": "uuid",
        "name": "string",
        "type": "string",
        "created_at": "timestamp"
    }
]
```

### WebSocket Endpoints

#### Connect to WebSocket
- URL: `/ws`
- Query Parameters:
  - `token`: JWT token for authentication
- Protocol: `WebSocket`

#### Message Types
1. Chat Message:
```json
{
    "type": "message",
    "channel_id": "uuid",
    "content": "string"
}
```

2. Channel Subscription:
```json
{
    "type": "subscribe",
    "channel_id": "uuid"
}
```

## Testing

Added comprehensive test suite for WebSocket functionality:
- Client registration/unregistration
- Message broadcasting
- Channel subscription
- Message persistence

Example test:
```go
func TestHub(t *testing.T) {
    mockMsgService := &MockMessageService{}
    hub := NewHub(mockMsgService)
    // ... test implementation
}
```

## Progress Summary
- ✓ Implemented WebSocket server
- ✓ Added authentication system
- ✓ Created server/channel management
- ✓ Added message persistence
- ✓ Implemented comprehensive testing
- ✓ Added API documentation

## Next Steps
1. Frontend Foundation:
   - Set up React project with Vite
   - Implement authentication UI
   - Create basic layouts and components
   - Set up Redux store
2. Additional Backend Features:
   - Add user presence system
   - Implement message history pagination
   - Add server member management
   - Add rate limiting

## Notes and Observations
- Used interface for MessageService to improve testability
- Implemented concurrent-safe operations with mutexes
- Added comprehensive error handling for WebSocket operations
- Used JWT for both HTTP and WebSocket authentication
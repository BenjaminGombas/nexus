# Nexus

A modern real-time chat platform inspired by Discord, focusing on seamless communication and community building.

## Features

### Current Features (v0.1)
- User authentication (registration/login)
- JWT-based authorization
- RESTful API endpoints
- Database integration with PostgreSQL

### Upcoming Features
- Real-time messaging using WebSocket
- Server (workspace) management
- Channel system
- User presence
- Direct messaging
- User profiles

## Tech Stack

### Backend
- Go
- Fiber (web framework)
- PostgreSQL (database)
- JWT (authentication)
- WebSocket (coming soon)

### Frontend (Coming Soon)
- React
- Redux Toolkit
- TailwindCSS
- shadcn/ui

## Getting Started

### Prerequisites
- Go 1.21 or later
- Docker and Docker Compose
- PostgreSQL (via Docker)

## API Endpoints

### Authentication
- **Register**: `POST /api/v1/auth/register`
- **Login**: `POST /api/v1/auth/login`

## Project Structure
```
/nexus
├── /frontend           # React frontend (coming soon)
├── /backend
│   ├── /cmd           # Main applications
│   ├── /internal      # Private application packages
│   ├── /pkg           # Public packages
│   └── /config        # Configuration files
└── /docs              # Documentation
```

## Development

### Running Tests
```cmd
cd backend
go test ./...
```

### Database Migrations
Currently handled through raw SQL in `internal/database/migrations/`.

## License
This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments
- Inspired by Discord's functionality and user experience
- Built using modern Go practices and patterns
- Utilizing high-performance web frameworks and libraries

## Contact
Benjamin Gombas - benhgombas@gmail.com

Project Link: [repository-url]
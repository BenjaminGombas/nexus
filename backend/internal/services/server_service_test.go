package services

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupServerTest(t *testing.T) (*ServerService, *UserService) {
	_, err := testDB.Exec(`
		TRUNCATE servers CASCADE;
		TRUNCATE users CASCADE;
	`)
	if err != nil {
		t.Fatalf("Could not clean test database: %v", err)
	}

	return NewServerService(testDB), NewUserService(testDB)
}

func TestCreateServer(t *testing.T) {
	serverService, userService := setupServerTest(t)

	user, err := userService.CreateUser("testuser", "test@example.com", "password123")
	assert.NoError(t, err)

	tests := []struct {
		name          string
		serverName    string
		ownerID       uuid.UUID
		expectedError bool
	}{
		{
			name:          "Valid server creation",
			serverName:    "Test Server",
			ownerID:       user.ID,
			expectedError: false,
		},
		{
			name:          "Empty server name",
			serverName:    "",
			ownerID:       user.ID,
			expectedError: true,
		},
		{
			name:          "Invalid owner ID",
			serverName:    "Test Server",
			ownerID:       uuid.New(),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, err := serverService.CreateServer(tt.serverName, tt.ownerID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, server)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, server)
				assert.Equal(t, tt.serverName, server.Name)
				assert.Equal(t, tt.ownerID, server.OwnerID)
			}
		})
	}
}

func TestGetUserServers(t *testing.T) {
	serverService, userService := setupServerTest(t)

	user, err := userService.CreateUser("testuser", "test@example.com", "password123")
	assert.NoError(t, err)

	otherUser, err := userService.CreateUser("other", "other@example.com", "password123")
	assert.NoError(t, err)

	for i := 0; i < 3; i++ {
		_, err := serverService.CreateServer("Test Server", user.ID)
		assert.NoError(t, err)
	}

	tests := []struct {
		name          string
		userID        uuid.UUID
		expectedCount int
		expectedError bool
	}{
		{
			name:          "User with servers",
			userID:        user.ID,
			expectedCount: 3,
			expectedError: false,
		},
		{
			name:          "User without servers",
			userID:        otherUser.ID,
			expectedCount: 0,
			expectedError: false,
		},
		{
			name:          "Non-existent user",
			userID:        uuid.New(),
			expectedCount: 0,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			servers, err := serverService.GetUserServers(tt.userID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, servers, tt.expectedCount)
			}
		})
	}
}

func TestJoinServer(t *testing.T) {
	serverService, userService := setupServerTest(t)

	owner, err := userService.CreateUser("owner", "owner@example.com", "password123")
	assert.NoError(t, err)

	user, err := userService.CreateUser("user", "user@example.com", "password123")
	assert.NoError(t, err)

	server, err := serverService.CreateServer("Test Server", owner.ID)
	assert.NoError(t, err)

	tests := []struct {
		name          string
		serverID      uuid.UUID
		userID        uuid.UUID
		expectedError bool
	}{
		{
			name:          "Valid server join",
			serverID:      server.ID,
			userID:        user.ID,
			expectedError: false,
		},
		{
			name:          "Already joined server",
			serverID:      server.ID,
			userID:        owner.ID,
			expectedError: false,
		},
		{
			name:          "Non-existent server",
			serverID:      uuid.New(),
			userID:        user.ID,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := serverService.JoinServer(tt.serverID, tt.userID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

package services

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupChannelTest(t *testing.T) (*ChannelService, *ServerService, *UserService) {
	_, err := testDB.Exec(`
		TRUNCATE channels CASCADE;
		TRUNCATE servers CASCADE;
		TRUNCATE users CASCADE;
	`)
	if err != nil {
		t.Fatalf("Could not clean test database: %v", err)
	}

	return NewChannelService(testDB), NewServerService(testDB), NewUserService(testDB)
}

func TestCreateChannel(t *testing.T) {
	channelService, serverService, userService := setupChannelTest(t)

	user, err := userService.CreateUser("testuser", "test@example.com", "password123")
	assert.NoError(t, err)

	server, err := serverService.CreateServer("Test Server", user.ID)
	assert.NoError(t, err)

	tests := []struct {
		name          string
		serverID      uuid.UUID
		channelName   string
		channelType   string
		expectedError bool
	}{
		{
			name:          "Valid text channel",
			serverID:      server.ID,
			channelName:   "general",
			channelType:   "text",
			expectedError: false,
		},
		{
			name:          "Valid voice channel",
			serverID:      server.ID,
			channelName:   "voice-chat",
			channelType:   "voice",
			expectedError: false,
		},
		{
			name:          "Empty channel name",
			serverID:      server.ID,
			channelName:   "",
			channelType:   "text",
			expectedError: true,
		},
		{
			name:          "Invalid channel type",
			serverID:      server.ID,
			channelName:   "test",
			channelType:   "invalid",
			expectedError: true,
		},
		{
			name:          "Non-existent server",
			serverID:      uuid.New(),
			channelName:   "test",
			channelType:   "text",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			channel, err := channelService.CreateChannel(tt.serverID, tt.channelName, tt.channelType)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, channel)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, channel)
				assert.Equal(t, tt.channelName, channel.Name)
				assert.Equal(t, tt.channelType, channel.Type)
				assert.Equal(t, tt.serverID, channel.ServerID)
			}
		})
	}
}

func TestGetServerChannels(t *testing.T) {
	channelService, serverService, userService := setupChannelTest(t)

	user, err := userService.CreateUser("testuser", "test@example.com", "password123")
	assert.NoError(t, err)

	server, err := serverService.CreateServer("Test Server", user.ID)
	assert.NoError(t, err)

	otherServer, err := serverService.CreateServer("Other Server", user.ID)
	assert.NoError(t, err)

	// Create channels for test server
	for i := 0; i < 3; i++ {
		_, err := channelService.CreateChannel(server.ID, "test-channel", "text")
		assert.NoError(t, err)
	}

	tests := []struct {
		name          string
		serverID      uuid.UUID
		expectedCount int
		expectedError bool
	}{
		{
			name:          "Server with channels",
			serverID:      server.ID,
			expectedCount: 3,
			expectedError: false,
		},
		{
			name:          "Server without channels",
			serverID:      otherServer.ID,
			expectedCount: 0,
			expectedError: false,
		},
		{
			name:          "Non-existent server",
			serverID:      uuid.New(),
			expectedCount: 0,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			channels, err := channelService.GetServerChannels(tt.serverID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, channels, tt.expectedCount)
			}
		})
	}
}

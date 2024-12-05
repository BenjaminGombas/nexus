package services

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupMessageTest(t *testing.T) (*MessageService, *UserService, *ChannelService, *ServerService) {
	// Clear test data
	_, err := testDB.Exec(`
		TRUNCATE messages CASCADE;
		TRUNCATE channels CASCADE;
		TRUNCATE servers CASCADE;
		TRUNCATE users CASCADE;
	`)
	if err != nil {
		t.Fatalf("Could not clean test database: %v", err)
	}

	return NewMessageService(testDB),
		NewUserService(testDB),
		NewChannelService(testDB),
		NewServerService(testDB)
}

func TestCreateMessage(t *testing.T) {
	messageService, userService, channelService, serverService := setupMessageTest(t)

	user, err := userService.CreateUser("testuser", "test@example.com", "password123")
	assert.NoError(t, err)

	server, err := serverService.CreateServer("Test Server", user.ID)
	assert.NoError(t, err)

	channel, err := channelService.CreateChannel(server.ID, "general", "text")
	assert.NoError(t, err)

	tests := []struct {
		name          string
		channelID     uuid.UUID
		userID        uuid.UUID
		content       string
		expectedError bool
	}{
		{
			name:          "Valid message creation",
			channelID:     channel.ID,
			userID:        user.ID,
			content:       "Hello, world!",
			expectedError: false,
		},
		{
			name:          "Invalid channel ID",
			channelID:     uuid.New(),
			userID:        user.ID,
			content:       "Hello, world!",
			expectedError: true,
		},
		{
			name:          "Empty content",
			channelID:     channel.ID,
			userID:        user.ID,
			content:       "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message, err := messageService.CreateMessage(tt.channelID, tt.userID, tt.content)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, message)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, message)
				assert.Equal(t, tt.channelID, message.ChannelID)
				assert.Equal(t, tt.userID, message.UserID)
				assert.Equal(t, tt.content, message.Content)
			}
		})
	}
}

func TestGetChannelMessages(t *testing.T) {
	messageService, userService, channelService, serverService := setupMessageTest(t)

	user, err := userService.CreateUser("testuser", "test@example.com", "password123")
	assert.NoError(t, err)

	server, err := serverService.CreateServer("Test Server", user.ID)
	assert.NoError(t, err)

	channel, err := channelService.CreateChannel(server.ID, "general", "text")
	assert.NoError(t, err)

	createdMessages := make([]Message, 0)
	for i := 0; i < 5; i++ {
		msg, err := messageService.CreateMessage(channel.ID, user.ID, "Test message")
		assert.NoError(t, err)
		createdMessages = append(createdMessages, *msg)
		t.Logf("Created message with time: %v", msg.CreatedAt)
	}

	queryTime := time.Now().UTC().Add(time.Second) // Add buffer for any time differences

	var latestTime time.Time
	err = testDB.QueryRow("SELECT created_at FROM messages WHERE channel_id = $1 ORDER BY created_at DESC LIMIT 1", channel.ID).Scan(&latestTime)
	assert.NoError(t, err)
	t.Logf("Latest message time in DB: %v", latestTime)
	t.Logf("Query time used: %v", queryTime)

	tests := []struct {
		name          string
		channelID     uuid.UUID
		limit         int
		expectedCount int
		expectedError bool
	}{
		{
			name:          "Get all messages",
			channelID:     channel.ID,
			limit:         10,
			expectedCount: 5,
			expectedError: false,
		},
		{
			name:          "Get limited messages",
			channelID:     channel.ID,
			limit:         3,
			expectedCount: 3,
			expectedError: false,
		},
		{
			name:          "Invalid channel ID",
			channelID:     uuid.New(),
			limit:         10,
			expectedCount: 0,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messages, err := messageService.GetChannelMessages(tt.channelID, tt.limit, queryTime)
			if err != nil {
				t.Logf("Error getting messages: %v", err)
			}
			if len(messages) > 0 {
				t.Logf("First message time: %v", messages[0].CreatedAt)
			}

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, messages, tt.expectedCount)
			}
		})
	}
}

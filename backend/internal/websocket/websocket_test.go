package websocket

import (
	"testing"
	"time"

	"nexus/internal/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockConn struct {
	mock.Mock
	closed bool
}

func (m *MockConn) WriteJSON(v interface{}) error {
	args := m.Called(v)
	return args.Error(0)
}

func (m *MockConn) Close() error {
	m.closed = true
	return nil
}

type MockMessageService struct {
	mock.Mock
	MessageServiceInterface
}

func (m *MockMessageService) CreateMessage(channelID, userID uuid.UUID, content string) (*services.Message, error) {
	args := m.Called(channelID, userID, content)
	if msg, ok := args.Get(0).(*services.Message); ok {
		return msg, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMessageService) GetChannelMessages(channelID uuid.UUID, limit int, before time.Time) ([]services.Message, error) {
	args := m.Called(channelID, limit, before)
	if msgs, ok := args.Get(0).([]services.Message); ok {
		return msgs, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestHub(t *testing.T) {
	mockMsgService := &MockMessageService{}
	mockMsgService.On("CreateMessage", mock.Anything, mock.Anything, mock.Anything).Return(
		&services.Message{
			ID:      uuid.New(),
			Content: "test message",
		}, nil)

	hub := NewHub(mockMsgService)
	go hub.Run()

	t.Run("Client Registration", func(t *testing.T) {
		conn := &MockConn{}
		client := &Client{
			ID:       uuid.New().String(),
			UserID:   uuid.New(),
			Conn:     conn,
			Hub:      hub,
			Channels: make(map[uuid.UUID]bool),
		}

		hub.register <- client
		time.Sleep(100 * time.Millisecond)
		assert.Contains(t, hub.clients, client)
	})

	t.Run("Client Unregistration", func(t *testing.T) {
		conn := &MockConn{}
		client := &Client{
			ID:       uuid.New().String(),
			UserID:   uuid.New(),
			Conn:     conn,
			Hub:      hub,
			Channels: make(map[uuid.UUID]bool),
		}

		hub.register <- client
		time.Sleep(100 * time.Millisecond)
		hub.unregister <- client
		time.Sleep(100 * time.Millisecond)
		assert.NotContains(t, hub.clients, client)
		assert.True(t, conn.closed)
	})

	t.Run("Message Broadcast", func(t *testing.T) {
		channelID := uuid.New()
		conn1, conn2 := &MockConn{}, &MockConn{}

		client1 := &Client{
			ID:       uuid.New().String(),
			UserID:   uuid.New(),
			Conn:     conn1,
			Hub:      hub,
			Channels: map[uuid.UUID]bool{channelID: true},
		}

		client2 := &Client{
			ID:       uuid.New().String(),
			UserID:   uuid.New(),
			Conn:     conn2,
			Hub:      hub,
			Channels: map[uuid.UUID]bool{channelID: true},
		}

		hub.register <- client1
		hub.register <- client2
		time.Sleep(100 * time.Millisecond)

		msg := Message{
			Type:      "message",
			ChannelID: channelID,
			Content:   "test message",
		}

		conn1.On("WriteJSON", mock.Anything).Return(nil)
		conn2.On("WriteJSON", mock.Anything).Return(nil)

		hub.broadcast <- msg
		time.Sleep(100 * time.Millisecond)

		conn1.AssertCalled(t, "WriteJSON", msg)
		conn2.AssertCalled(t, "WriteJSON", msg)
	})

	t.Run("Channel Subscription", func(t *testing.T) {
		channelID := uuid.New()
		conn := &MockConn{}
		client := &Client{
			ID:       uuid.New().String(),
			UserID:   uuid.New(),
			Conn:     conn,
			Hub:      hub,
			Channels: make(map[uuid.UUID]bool),
		}

		hub.register <- client
		time.Sleep(100 * time.Millisecond)

		client.SubscribeToChannel(channelID)
		assert.True(t, client.Channels[channelID])

		client.UnsubscribeFromChannel(channelID)
		assert.False(t, client.Channels[channelID])
	})
}

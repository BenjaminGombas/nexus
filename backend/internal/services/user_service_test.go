package services

import (
	"database/sql"
	"log"
	"nexus/internal/database"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Set test database name
	os.Setenv("DB_NAME", "nexus_test")

	var err error
	testDB, err = database.NewConnection()
	if err != nil {
		log.Fatalf("Could not connect to test database: %v", err)
	}

	code := m.Run()

	testDB.Close()
	os.Exit(code)
}

func setupTest(t *testing.T) *UserService {
	// Clear test data
	_, err := testDB.Exec(`
		TRUNCATE users CASCADE;
	`)
	if err != nil {
		t.Fatalf("Could not clean test database: %v", err)
	}

	return NewUserService(testDB)
}

func TestCreateUser(t *testing.T) {
	userService := setupTest(t)

	tests := []struct {
		name          string
		username      string
		email         string
		password      string
		expectedError bool
	}{
		{
			name:          "Valid user creation",
			username:      "testuser",
			email:         "test@example.com",
			password:      "password123",
			expectedError: false,
		},
		{
			name:          "Duplicate email",
			username:      "testuser2",
			email:         "test@example.com",
			password:      "password123",
			expectedError: true,
		},
		{
			name:          "Empty username",
			username:      "",
			email:         "test2@example.com",
			password:      "password123",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := userService.CreateUser(tt.username, tt.email, tt.password)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.username, user.Username)
				assert.Equal(t, tt.email, user.Email)
			}
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	userService := setupTest(t)

	// Create test user
	testUser, err := userService.CreateUser("testuser", "test@example.com", "password123")
	assert.NoError(t, err)

	tests := []struct {
		name          string
		email         string
		expectedError bool
	}{
		{
			name:          "Existing user",
			email:         "test@example.com",
			expectedError: false,
		},
		{
			name:          "Non-existent user",
			email:         "nonexistent@example.com",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := userService.GetUserByEmail(tt.email)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, testUser.ID, user.ID)
				assert.Equal(t, testUser.Email, user.Email)
			}
		})
	}
}

package services

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"educ-retro/internal/models"
	"educ-retro/internal/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// MockUserRepository é um mock simples do UserRepository
type MockUserRepository struct {
	users  map[uuid.UUID]*models.User
	emails map[string]*models.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:  make(map[uuid.UUID]*models.User),
		emails: make(map[string]*models.User),
	}
}

func (m *MockUserRepository) Create(user *models.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	if _, exists := m.emails[user.Email]; exists {
		return errors.New("duplicate key value violates unique constraint")
	}

	m.users[user.ID] = user
	m.emails[user.Email] = user
	return nil
}

func (m *MockUserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, sql.ErrNoRows
	}
	userCopy := *user
	return &userCopy, nil
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	user, exists := m.emails[email]
	if !exists {
		return nil, sql.ErrNoRows
	}
	userCopy := *user
	return &userCopy, nil
}

func (m *MockUserRepository) Update(user *models.User) error {
	if _, exists := m.users[user.ID]; !exists {
		return sql.ErrNoRows
	}
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	m.emails[user.Email] = user
	return nil
}

func (m *MockUserRepository) Delete(id uuid.UUID) error {
	user, exists := m.users[id]
	if !exists {
		return sql.ErrNoRows
	}
	delete(m.users, id)
	delete(m.emails, user.Email)
	return nil
}

func (m *MockUserRepository) GetUsersByIDs(ids []uuid.UUID) ([]models.User, error) {
	var users []models.User
	for _, id := range ids {
		if user, exists := m.users[id]; exists {
			users = append(users, *user)
		}
	}
	return users, nil
}

func TestNewUserService(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewUserService(mockRepo)

	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.userRepo)
}

func TestUserService_Register(t *testing.T) {
	tests := []struct {
		name          string
		request       *models.UserCreateRequest
		setupMock     func(*MockUserRepository)
		expectedError string
	}{
		{
			name: "successful registration",
			request: &models.UserCreateRequest{
				Email:    "test@example.com",
				Name:     "Test User",
				Password: "password123",
			},
			setupMock: func(m *MockUserRepository) {
				// Mock vazio - sem setup necessário
			},
			expectedError: "",
		},
		{
			name: "user already exists",
			request: &models.UserCreateRequest{
				Email:    "existing@example.com",
				Name:     "Existing User",
				Password: "password123",
			},
			setupMock: func(m *MockUserRepository) {
				// Simular usuário existente
				existingUser := &models.User{
					ID:    uuid.New(),
					Email: "existing@example.com",
					Name:  "Existing User",
				}
				m.emails["existing@example.com"] = existingUser
			},
			expectedError: "user with this email already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockUserRepository()
			tt.setupMock(mockRepo)

			service := NewUserService(mockRepo)
			userResponse, token, err := service.Register(tt.request)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, userResponse)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, userResponse)
				assert.NotEmpty(t, token)
				assert.Equal(t, tt.request.Email, userResponse.Email)
				assert.Equal(t, tt.request.Name, userResponse.Name)
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
	tests := []struct {
		name          string
		request       *models.UserLoginRequest
		setupMock     func(*MockUserRepository)
		expectedError string
	}{
		{
			name: "successful login",
			request: &models.UserLoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(m *MockUserRepository) {
				hashedPassword, _ := utils.HashPassword("password123")
				user := &models.User{
					ID:        uuid.New(),
					Email:     "test@example.com",
					Name:      "Test User",
					Password:  hashedPassword,
					CreatedAt: time.Now(),
				}
				m.emails["test@example.com"] = user
			},
			expectedError: "",
		},
		{
			name: "user not found",
			request: &models.UserLoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			setupMock: func(m *MockUserRepository) {
				// Mock vazio - usuário não existe
			},
			expectedError: "invalid credentials",
		},
		{
			name: "wrong password",
			request: &models.UserLoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			setupMock: func(m *MockUserRepository) {
				hashedPassword, _ := utils.HashPassword("password123")
				user := &models.User{
					ID:        uuid.New(),
					Email:     "test@example.com",
					Name:      "Test User",
					Password:  hashedPassword,
					CreatedAt: time.Now(),
				}
				m.emails["test@example.com"] = user
			},
			expectedError: "invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockUserRepository()
			tt.setupMock(mockRepo)

			service := NewUserService(mockRepo)
			userResponse, token, err := service.Login(tt.request)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, userResponse)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, userResponse)
				assert.NotEmpty(t, token)
				assert.Equal(t, tt.request.Email, userResponse.Email)
			}
		})
	}
}

func TestUserService_GetProfile(t *testing.T) {
	tests := []struct {
		name          string
		userID        uuid.UUID
		setupMock     func(*MockUserRepository)
		expectedError string
	}{
		{
			name:   "successful get profile",
			userID: uuid.New(),
			setupMock: func(m *MockUserRepository) {
				// O userID será definido no teste
			},
			expectedError: "",
		},
		{
			name:   "user not found",
			userID: uuid.New(),
			setupMock: func(m *MockUserRepository) {
				// Mock vazio - usuário não existe
			},
			expectedError: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockUserRepository()
			tt.setupMock(mockRepo)

			// Para o teste de sucesso, criar o usuário com o userID correto
			if tt.name == "successful get profile" {
				user := &models.User{
					ID:        tt.userID,
					Email:     "test@example.com",
					Name:      "Test User",
					CreatedAt: time.Now(),
				}
				mockRepo.users[tt.userID] = user
			}

			service := NewUserService(mockRepo)
			userResponse, err := service.GetProfile(tt.userID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, userResponse)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, userResponse)
				assert.NotEqual(t, uuid.Nil, userResponse.ID)
			}
		})
	}
}

func TestUserService_UpdateProfile(t *testing.T) {
	tests := []struct {
		name          string
		userID        uuid.UUID
		userName      string
		avatar        string
		setupMock     func(*MockUserRepository)
		expectedError string
	}{
		{
			name:     "successful update",
			userID:   uuid.New(),
			userName: "Updated Name",
			avatar:   "https://example.com/avatar.jpg",
			setupMock: func(m *MockUserRepository) {
				// O userID será definido no teste
			},
			expectedError: "",
		},
		{
			name:     "user not found",
			userID:   uuid.New(),
			userName: "Updated Name",
			avatar:   "",
			setupMock: func(m *MockUserRepository) {
				// Mock vazio - usuário não existe
			},
			expectedError: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockUserRepository()
			tt.setupMock(mockRepo)

			// Para o teste de sucesso, criar o usuário com o userID correto
			if tt.name == "successful update" {
				user := &models.User{
					ID:        tt.userID,
					Email:     "test@example.com",
					Name:      "Original Name",
					CreatedAt: time.Now(),
				}
				mockRepo.users[tt.userID] = user
			}

			service := NewUserService(mockRepo)
			userResponse, err := service.UpdateProfile(tt.userID, tt.userName, tt.avatar)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, userResponse)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, userResponse)
				assert.Equal(t, tt.userName, userResponse.Name)
			}
		})
	}
}

// Testes de edge cases simplificados
func TestUserService_EdgeCases(t *testing.T) {
	t.Run("register with special characters", func(t *testing.T) {
		mockRepo := NewMockUserRepository()
		service := NewUserService(mockRepo)

		request := &models.UserCreateRequest{
			Email:    "test+special@example.com",
			Name:     "José da Silva & Cia.",
			Password: "password123",
		}

		userResponse, token, err := service.Register(request)

		assert.NoError(t, err)
		assert.NotNil(t, userResponse)
		assert.NotEmpty(t, token)
		assert.Equal(t, request.Email, userResponse.Email)
		assert.Equal(t, request.Name, userResponse.Name)
	})

	t.Run("login case sensitivity", func(t *testing.T) {
		mockRepo := NewMockUserRepository()
		service := NewUserService(mockRepo)

		// Criar usuário com email lowercase
		hashedPassword, _ := utils.HashPassword("password123")
		user := &models.User{
			ID:        uuid.New(),
			Email:     "test@example.com",
			Name:      "Test User",
			Password:  hashedPassword,
			CreatedAt: time.Now(),
		}
		mockRepo.emails["test@example.com"] = user

		// Tentar login com email uppercase
		request := &models.UserLoginRequest{
			Email:    "TEST@EXAMPLE.COM",
			Password: "password123",
		}

		userResponse, token, err := service.Login(request)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid credentials")
		assert.Nil(t, userResponse)
		assert.Empty(t, token)
	})
}

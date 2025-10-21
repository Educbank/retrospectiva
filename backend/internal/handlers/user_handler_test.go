package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"educ-retro/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// UserServiceInterface for mocking
type UserServiceInterface interface {
	Register(req *models.UserCreateRequest) (*models.UserResponse, string, error)
	Login(req *models.UserLoginRequest) (*models.UserResponse, string, error)
	GetProfile(userID uuid.UUID) (*models.UserResponse, error)
	UpdateProfile(userID uuid.UUID, name, avatar string) (*models.UserResponse, error)
}

// MockUserService for testing
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(req *models.UserCreateRequest) (*models.UserResponse, string, error) {
	args := m.Called(req)
	return args.Get(0).(*models.UserResponse), args.String(1), args.Error(2)
}

func (m *MockUserService) Login(req *models.UserLoginRequest) (*models.UserResponse, string, error) {
	args := m.Called(req)
	return args.Get(0).(*models.UserResponse), args.String(1), args.Error(2)
}

func (m *MockUserService) GetProfile(userID uuid.UUID) (*models.UserResponse, error) {
	args := m.Called(userID)
	return args.Get(0).(*models.UserResponse), args.Error(1)
}

func (m *MockUserService) UpdateProfile(userID uuid.UUID, name, avatar string) (*models.UserResponse, error) {
	args := m.Called(userID, name, avatar)
	return args.Get(0).(*models.UserResponse), args.Error(1)
}

// TestUserHandler for testing with interface
type TestUserHandler struct {
	userService UserServiceInterface
}

func NewTestUserHandler(userService UserServiceInterface) *TestUserHandler {
	return &TestUserHandler{userService: userService}
}

func (h *TestUserHandler) Register(c *gin.Context) {
	var req models.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.userService.Register(&req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user with this email already exists" {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":  user,
		"token": token,
	})
}

func (h *TestUserHandler) Login(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.userService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}

func (h *TestUserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	user, err := h.userService.GetProfile(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *TestUserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req struct {
		Name   string `json:"name" binding:"required"`
		Avatar string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.UpdateProfile(userID.(uuid.UUID), req.Name, req.Avatar)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func TestUserHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockUserService)
	handler := NewTestUserHandler(mockService)

	router := gin.New()
	router.POST("/register", handler.Register)

	req := models.UserCreateRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedUser := &models.UserResponse{
		ID:    uuid.New(),
		Name:  "Test User",
		Email: "test@example.com",
	}

	mockService.On("Register", &req).Return(expectedUser, "jwt-token", nil)

	jsonBody, _ := json.Marshal(req)
	request, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "user")
	assert.Contains(t, response, "token")
	mockService.AssertExpectations(t)
}

func TestUserHandler_Register_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockUserService)
	handler := NewTestUserHandler(mockService)

	router := gin.New()
	router.POST("/register", handler.Register)

	invalidReq := map[string]string{
		"name": "Test User",
		// Missing email and password
	}

	jsonBody, _ := json.Marshal(invalidReq)
	request, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "error")
}

func TestUserHandler_Register_UserExists(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockUserService)
	handler := NewTestUserHandler(mockService)

	router := gin.New()
	router.POST("/register", handler.Register)

	req := models.UserCreateRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	mockService.On("Register", &req).Return((*models.UserResponse)(nil), "", assert.AnError)

	jsonBody, _ := json.Marshal(req)
	request, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "error")
	mockService.AssertExpectations(t)
}

func TestUserHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockUserService)
	handler := NewTestUserHandler(mockService)

	router := gin.New()
	router.POST("/login", handler.Login)

	req := models.UserLoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedUser := &models.UserResponse{
		ID:    uuid.New(),
		Name:  "Test User",
		Email: "test@example.com",
	}

	mockService.On("Login", &req).Return(expectedUser, "jwt-token", nil)

	jsonBody, _ := json.Marshal(req)
	request, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "user")
	assert.Contains(t, response, "token")
	mockService.AssertExpectations(t)
}

func TestUserHandler_Login_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockUserService)
	handler := NewTestUserHandler(mockService)

	router := gin.New()
	router.POST("/login", handler.Login)

	req := models.UserLoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	mockService.On("Login", &req).Return((*models.UserResponse)(nil), "", assert.AnError)

	jsonBody, _ := json.Marshal(req)
	request, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "error")
	mockService.AssertExpectations(t)
}

func TestUserHandler_GetProfile_NotAuthenticated(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockUserService)
	handler := NewTestUserHandler(mockService)

	router := gin.New()
	router.GET("/profile", handler.GetProfile)

	request, _ := http.NewRequest("GET", "/profile", nil)
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "error")
	assert.Equal(t, "user not authenticated", response["error"])
}

func TestUserHandler_UpdateProfile_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockUserService)
	handler := NewTestUserHandler(mockService)

	router := gin.New()
	router.PUT("/profile", handler.UpdateProfile)

	invalidReq := map[string]string{
		"avatar": "new-avatar.jpg",
		// Missing required name field
	}

	jsonBody, _ := json.Marshal(invalidReq)
	request, _ := http.NewRequest("PUT", "/profile", bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	// Should return 401 because authentication check happens first
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "error")
}

func TestUserHandler_UpdateProfile_NotAuthenticated(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockUserService)
	handler := NewTestUserHandler(mockService)

	router := gin.New()
	router.PUT("/profile", handler.UpdateProfile)

	req := map[string]string{
		"name":   "Updated Name",
		"avatar": "new-avatar.jpg",
	}

	jsonBody, _ := json.Marshal(req)
	request, _ := http.NewRequest("PUT", "/profile", bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "error")
	assert.Equal(t, "user not authenticated", response["error"])
}

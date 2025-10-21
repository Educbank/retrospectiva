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

// TeamServiceInterface for mocking
type TeamServiceInterface interface {
	CreateTeam(userID uuid.UUID, req *models.TeamCreateRequest) (*models.Team, error)
	GetUserTeams(userID uuid.UUID) ([]models.TeamWithCounts, error)
	GetTeam(teamID, userID uuid.UUID) (*models.Team, error)
	UpdateTeam(teamID, userID uuid.UUID, req *models.TeamCreateRequest) (*models.Team, error)
	DeleteTeam(teamID, userID uuid.UUID) error
	AddMember(teamID, userID uuid.UUID, req *models.TeamInviteRequest) error
	RemoveMember(teamID, userID, memberID uuid.UUID) error
	UpdateMemberRole(teamID, userID, memberID uuid.UUID, role string) error
}

// MockTeamService for testing
type MockTeamService struct {
	mock.Mock
}

func (m *MockTeamService) CreateTeam(userID uuid.UUID, req *models.TeamCreateRequest) (*models.Team, error) {
	args := m.Called(userID, req)
	return args.Get(0).(*models.Team), args.Error(1)
}

func (m *MockTeamService) GetUserTeams(userID uuid.UUID) ([]models.TeamWithCounts, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.TeamWithCounts), args.Error(1)
}

func (m *MockTeamService) GetTeam(teamID, userID uuid.UUID) (*models.Team, error) {
	args := m.Called(teamID, userID)
	return args.Get(0).(*models.Team), args.Error(1)
}

func (m *MockTeamService) UpdateTeam(teamID, userID uuid.UUID, req *models.TeamCreateRequest) (*models.Team, error) {
	args := m.Called(teamID, userID, req)
	return args.Get(0).(*models.Team), args.Error(1)
}

func (m *MockTeamService) DeleteTeam(teamID, userID uuid.UUID) error {
	args := m.Called(teamID, userID)
	return args.Error(0)
}

func (m *MockTeamService) AddMember(teamID, userID uuid.UUID, req *models.TeamInviteRequest) error {
	args := m.Called(teamID, userID, req)
	return args.Error(0)
}

func (m *MockTeamService) RemoveMember(teamID, userID, memberID uuid.UUID) error {
	args := m.Called(teamID, userID, memberID)
	return args.Error(0)
}

func (m *MockTeamService) UpdateMemberRole(teamID, userID, memberID uuid.UUID, role string) error {
	args := m.Called(teamID, userID, memberID, role)
	return args.Error(0)
}

// TestTeamHandler for testing with interface
type TestTeamHandler struct {
	teamService TeamServiceInterface
}

func NewTestTeamHandler(teamService TeamServiceInterface) *TestTeamHandler {
	return &TestTeamHandler{teamService: teamService}
}

func (h *TestTeamHandler) CreateTeam(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req models.TeamCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team, err := h.teamService.CreateTeam(userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, team)
}

func (h *TestTeamHandler) GetUserTeams(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	teams, err := h.teamService.GetUserTeams(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teams)
}

func (h *TestTeamHandler) GetTeam(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	teamIDStr := c.Param("id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	team, err := h.teamService.GetTeam(teamID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, team)
}

func (h *TestTeamHandler) UpdateTeam(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	teamIDStr := c.Param("id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	var req models.TeamCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team, err := h.teamService.UpdateTeam(teamID, userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, team)
}

func (h *TestTeamHandler) DeleteTeam(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	teamIDStr := c.Param("id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	err = h.teamService.DeleteTeam(teamID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team deleted successfully"})
}

func TestTeamHandler_CreateTeam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockTeamService)
	handler := NewTestTeamHandler(mockService)

	router := gin.New()
	router.POST("/teams", handler.CreateTeam)

	req := models.TeamCreateRequest{
		Name:        "Test Team",
		Description: "Test Description",
	}

	expectedTeam := &models.Team{
		ID:          uuid.New(),
		Name:        "Test Team",
		Description: &req.Description,
		OwnerID:     uuid.New(),
	}

	mockService.On("CreateTeam", mock.AnythingOfType("uuid.UUID"), &req).Return(expectedTeam, nil)

	jsonBody, _ := json.Marshal(req)
	request, _ := http.NewRequest("POST", "/teams", bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	// Should return 401 because no authentication
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "error")
	assert.Equal(t, "user not authenticated", response["error"])
}

func TestTeamHandler_CreateTeam_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockTeamService)
	handler := NewTestTeamHandler(mockService)

	router := gin.New()
	router.POST("/teams", handler.CreateTeam)

	invalidReq := map[string]string{
		"description": "Test Description",
		// Missing required name field
	}

	jsonBody, _ := json.Marshal(invalidReq)
	request, _ := http.NewRequest("POST", "/teams", bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	// Should return 401 because authentication check happens first
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "error")
}

func TestTeamHandler_GetUserTeams_NotAuthenticated(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockTeamService)
	handler := NewTestTeamHandler(mockService)

	router := gin.New()
	router.GET("/teams", handler.GetUserTeams)

	request, _ := http.NewRequest("GET", "/teams", nil)
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "error")
	assert.Equal(t, "user not authenticated", response["error"])
}

func TestTeamHandler_GetTeam_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockTeamService)
	handler := NewTestTeamHandler(mockService)

	router := gin.New()
	router.GET("/teams/:id", handler.GetTeam)

	request, _ := http.NewRequest("GET", "/teams/invalid-id", nil)
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	// Should return 401 because authentication check happens first
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "error")
}

func TestTeamHandler_UpdateTeam_NotAuthenticated(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockTeamService)
	handler := NewTestTeamHandler(mockService)

	router := gin.New()
	router.PUT("/teams/:id", handler.UpdateTeam)

	req := models.TeamCreateRequest{
		Name:        "Updated Team",
		Description: "Updated Description",
	}

	jsonBody, _ := json.Marshal(req)
	request, _ := http.NewRequest("PUT", "/teams/123e4567-e89b-12d3-a456-426614174000", bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "error")
	assert.Equal(t, "user not authenticated", response["error"])
}

func TestTeamHandler_DeleteTeam_NotAuthenticated(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockTeamService)
	handler := NewTestTeamHandler(mockService)

	router := gin.New()
	router.DELETE("/teams/:id", handler.DeleteTeam)

	request, _ := http.NewRequest("DELETE", "/teams/123e4567-e89b-12d3-a456-426614174000", nil)
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "error")
	assert.Equal(t, "user not authenticated", response["error"])
}

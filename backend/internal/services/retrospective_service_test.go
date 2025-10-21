package services

import (
	"database/sql"
	"testing"
	"time"

	"educ-retro/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// MockRetrospectiveRepository é um mock simples do RetrospectiveRepository
type MockRetrospectiveRepository struct {
	retrospectives map[uuid.UUID]*models.Retrospective
	details        map[uuid.UUID]*models.RetrospectiveWithDetails
}

func NewMockRetrospectiveRepository() *MockRetrospectiveRepository {
	return &MockRetrospectiveRepository{
		retrospectives: make(map[uuid.UUID]*models.Retrospective),
		details:        make(map[uuid.UUID]*models.RetrospectiveWithDetails),
	}
}

func (m *MockRetrospectiveRepository) Create(retrospective *models.Retrospective) error {
	if retrospective.ID == uuid.Nil {
		retrospective.ID = uuid.New()
	}
	now := time.Now()
	retrospective.CreatedAt = now
	retrospective.UpdatedAt = now
	m.retrospectives[retrospective.ID] = retrospective
	return nil
}

func (m *MockRetrospectiveRepository) GetByID(id uuid.UUID) (*models.Retrospective, error) {
	retrospective, exists := m.retrospectives[id]
	if !exists {
		return nil, sql.ErrNoRows
	}
	retrospectiveCopy := *retrospective
	return &retrospectiveCopy, nil
}

func (m *MockRetrospectiveRepository) GetAllRetrospectives() ([]models.Retrospective, error) {
	var retrospectives []models.Retrospective
	for _, retro := range m.retrospectives {
		retrospectives = append(retrospectives, *retro)
	}
	return retrospectives, nil
}

func (m *MockRetrospectiveRepository) GetRetrospectiveWithDetails(id uuid.UUID) (*models.RetrospectiveWithDetails, error) {
	details, exists := m.details[id]
	if !exists {
		return nil, sql.ErrNoRows
	}
	detailsCopy := *details
	return &detailsCopy, nil
}

func (m *MockRetrospectiveRepository) Update(retrospective *models.Retrospective) error {
	if _, exists := m.retrospectives[retrospective.ID]; !exists {
		return sql.ErrNoRows
	}
	retrospective.UpdatedAt = time.Now()
	m.retrospectives[retrospective.ID] = retrospective
	return nil
}

func (m *MockRetrospectiveRepository) Delete(id uuid.UUID) error {
	if _, exists := m.retrospectives[id]; !exists {
		return sql.ErrNoRows
	}
	delete(m.retrospectives, id)
	return nil
}

// Implementações vazias para os outros métodos
func (m *MockRetrospectiveRepository) UpdateStatus(id uuid.UUID, status models.RetrospectiveStatus) error {
	return nil
}
func (m *MockRetrospectiveRepository) GetRetrospectiveCount(id uuid.UUID) (int, error)   { return 1, nil }
func (m *MockRetrospectiveRepository) GetActionItemCount(id uuid.UUID) (int, error)      { return 0, nil }
func (m *MockRetrospectiveRepository) AddItem(item *models.RetrospectiveItem) error      { return nil }
func (m *MockRetrospectiveRepository) VoteItem(itemID, userID uuid.UUID) error           { return nil }
func (m *MockRetrospectiveRepository) AddActionItem(actionItem *models.ActionItem) error { return nil }
func (m *MockRetrospectiveRepository) RegisterParticipant(retrospectiveID, userID uuid.UUID) error {
	return nil
}
func (m *MockRetrospectiveRepository) GetParticipants(retrospectiveID uuid.UUID) ([]models.RetrospectiveParticipant, error) {
	return []models.RetrospectiveParticipant{}, nil
}
func (m *MockRetrospectiveRepository) GetItemByID(id uuid.UUID) (*models.RetrospectiveItem, error) {
	return nil, sql.ErrNoRows
}
func (m *MockRetrospectiveRepository) DeleteItem(id uuid.UUID) error { return nil }
func (m *MockRetrospectiveRepository) ReopenRetrospective(id uuid.UUID) error { return nil }
func (m *MockRetrospectiveRepository) CreateGroup(group *models.RetrospectiveGroup, itemIDs []uuid.UUID) error {
	return nil
}
func (m *MockRetrospectiveRepository) VoteGroup(groupID, userID uuid.UUID) error { return nil }
func (m *MockRetrospectiveRepository) GetGroupByID(id uuid.UUID) (*models.RetrospectiveGroup, error) {
	return nil, sql.ErrNoRows
}
func (m *MockRetrospectiveRepository) DeleteGroup(id uuid.UUID) error { return nil }
func (m *MockRetrospectiveRepository) MergeItems(sourceItemID, targetItemID uuid.UUID) (*models.RetrospectiveItem, error) {
	return nil, nil
}
func (m *MockRetrospectiveRepository) GetActionItemByID(id uuid.UUID) (*models.ActionItem, error) {
	return nil, sql.ErrNoRows
}
func (m *MockRetrospectiveRepository) UpdateActionItem(actionItemID uuid.UUID, req *models.ActionItemUpdateRequest) (*models.ActionItem, error) {
	return nil, nil
}
func (m *MockRetrospectiveRepository) DeleteActionItem(id uuid.UUID) error { return nil }
func (m *MockRetrospectiveRepository) UpdateTimer(id uuid.UUID, req models.TimerUpdateRequest) error {
	return nil
}

// MockTeamRepository para testes
type MockTeamRepositoryForRetro struct {
	teams map[uuid.UUID]*models.Team
}

func NewMockTeamRepositoryForRetro() *MockTeamRepositoryForRetro {
	return &MockTeamRepositoryForRetro{
		teams: make(map[uuid.UUID]*models.Team),
	}
}

func (m *MockTeamRepositoryForRetro) Create(team *models.Team) error { return nil }
func (m *MockTeamRepositoryForRetro) GetByID(id uuid.UUID) (*models.Team, error) {
	return nil, sql.ErrNoRows
}
func (m *MockTeamRepositoryForRetro) GetByUserID(userID uuid.UUID) ([]models.TeamWithCounts, error) {
	return []models.TeamWithCounts{}, nil
}
func (m *MockTeamRepositoryForRetro) Update(team *models.Team) error { return nil }
func (m *MockTeamRepositoryForRetro) Delete(id uuid.UUID) error      { return nil }
func (m *MockTeamRepositoryForRetro) IsMember(teamID, userID uuid.UUID) (bool, error) {
	return true, nil
}
func (m *MockTeamRepositoryForRetro) GetMemberRole(teamID, userID uuid.UUID) (string, error) {
	return "owner", nil
}
func (m *MockTeamRepositoryForRetro) GetMembers(teamID uuid.UUID) ([]models.TeamMemberWithUser, error) {
	return []models.TeamMemberWithUser{}, nil
}
func (m *MockTeamRepositoryForRetro) AddMember(teamID, userID uuid.UUID, role string) error {
	return nil
}
func (m *MockTeamRepositoryForRetro) RemoveMember(teamID, userID uuid.UUID) error { return nil }
func (m *MockTeamRepositoryForRetro) UpdateMemberRole(teamID, userID uuid.UUID, newRole string) error {
	return nil
}

func TestNewRetrospectiveService(t *testing.T) {
	mockRetroRepo := NewMockRetrospectiveRepository()
	mockTeamRepo := NewMockTeamRepositoryForRetro()
	service := NewRetrospectiveService(mockRetroRepo, mockTeamRepo)

	assert.NotNil(t, service)
	assert.Equal(t, mockRetroRepo, service.retroRepo)
	assert.Equal(t, mockTeamRepo, service.teamRepo)
}

func TestRetrospectiveService_CreateRetrospective(t *testing.T) {
	mockRetroRepo := NewMockRetrospectiveRepository()
	mockTeamRepo := NewMockTeamRepositoryForRetro()
	service := NewRetrospectiveService(mockRetroRepo, mockTeamRepo)

	userID := uuid.New()
	request := &models.RetrospectiveCreateRequest{
		Title:       "Test Retrospective",
		Description: "Test Description",
		Template:    "start_stop_continue",
	}

	retrospective, err := service.CreateRetrospective(userID, request)

	assert.NoError(t, err)
	assert.NotNil(t, retrospective)
	assert.Equal(t, request.Title, retrospective.Title)
	assert.Equal(t, userID, retrospective.CreatedBy)
	assert.Equal(t, models.RetroStatusPlanned, retrospective.Status)
}

func TestRetrospectiveService_GetUserRetrospectives(t *testing.T) {
	mockRetroRepo := NewMockRetrospectiveRepository()
	mockTeamRepo := NewMockTeamRepositoryForRetro()
	service := NewRetrospectiveService(mockRetroRepo, mockTeamRepo)

	userID := uuid.New()

	// Criar retrospectivas
	retro1 := &models.Retrospective{
		ID:        uuid.New(),
		Title:     "Retro 1",
		Status:    models.RetroStatusPlanned,
		CreatedBy: userID,
	}
	retro2 := &models.Retrospective{
		ID:        uuid.New(),
		Title:     "Retro 2",
		Status:    models.RetroStatusActive,
		CreatedBy: uuid.New(), // Outro usuário
	}

	mockRetroRepo.retrospectives[retro1.ID] = retro1
	mockRetroRepo.retrospectives[retro2.ID] = retro2

	// Setup details
	mockRetroRepo.details[retro1.ID] = &models.RetrospectiveWithDetails{
		Retrospective: *retro1,
	}
	mockRetroRepo.details[retro2.ID] = &models.RetrospectiveWithDetails{
		Retrospective: *retro2,
	}

	retrospectives, err := service.GetUserRetrospectives(userID)

	assert.NoError(t, err)
	assert.Len(t, retrospectives, 2) // Ambas devem aparecer
}

func TestRetrospectiveService_GetRetrospective(t *testing.T) {
	mockRetroRepo := NewMockRetrospectiveRepository()
	mockTeamRepo := NewMockTeamRepositoryForRetro()
	service := NewRetrospectiveService(mockRetroRepo, mockTeamRepo)

	retrospectiveID := uuid.New()
	userID := uuid.New()

	// Setup retrospective
	retrospective := &models.Retrospective{
		ID:        retrospectiveID,
		Title:     "Test Retrospective",
		CreatedBy: userID,
	}
	mockRetroRepo.retrospectives[retrospectiveID] = retrospective

	retrievedRetrospective, err := service.GetRetrospective(retrospectiveID, userID)

	assert.NoError(t, err)
	assert.NotNil(t, retrievedRetrospective)
	assert.Equal(t, retrospective.Title, retrievedRetrospective.Title)
}

func TestRetrospectiveService_GetRetrospective_AccessDenied(t *testing.T) {
	mockRetroRepo := NewMockRetrospectiveRepository()
	mockTeamRepo := NewMockTeamRepositoryForRetro()
	service := NewRetrospectiveService(mockRetroRepo, mockTeamRepo)

	retrospectiveID := uuid.New()
	userID := uuid.New()
	otherUserID := uuid.New()

	// Setup retrospective created by other user
	retrospective := &models.Retrospective{
		ID:        retrospectiveID,
		Title:     "Test Retrospective",
		CreatedBy: otherUserID,
	}
	mockRetroRepo.retrospectives[retrospectiveID] = retrospective

	retrievedRetrospective, err := service.GetRetrospective(retrospectiveID, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access denied")
	assert.Nil(t, retrievedRetrospective)
}

func TestRetrospectiveService_UpdateRetrospective(t *testing.T) {
	mockRetroRepo := NewMockRetrospectiveRepository()
	mockTeamRepo := NewMockTeamRepositoryForRetro()
	service := NewRetrospectiveService(mockRetroRepo, mockTeamRepo)

	retrospectiveID := uuid.New()
	userID := uuid.New()

	// Setup retrospective
	retrospective := &models.Retrospective{
		ID:        retrospectiveID,
		Title:     "Original Title",
		CreatedBy: userID,
	}
	mockRetroRepo.retrospectives[retrospectiveID] = retrospective

	request := &models.RetrospectiveCreateRequest{
		Title:       "Updated Title",
		Description: "Updated Description",
		Template:    "start_stop_continue",
	}

	updatedRetrospective, err := service.UpdateRetrospective(retrospectiveID, userID, request)

	assert.NoError(t, err)
	assert.NotNil(t, updatedRetrospective)
	assert.Equal(t, request.Title, updatedRetrospective.Title)
}

func TestRetrospectiveService_UpdateRetrospective_AccessDenied(t *testing.T) {
	mockRetroRepo := NewMockRetrospectiveRepository()
	mockTeamRepo := NewMockTeamRepositoryForRetro()
	service := NewRetrospectiveService(mockRetroRepo, mockTeamRepo)

	retrospectiveID := uuid.New()
	userID := uuid.New()
	otherUserID := uuid.New()

	// Setup retrospective created by other user
	retrospective := &models.Retrospective{
		ID:        retrospectiveID,
		Title:     "Original Title",
		CreatedBy: otherUserID,
	}
	mockRetroRepo.retrospectives[retrospectiveID] = retrospective

	request := &models.RetrospectiveCreateRequest{
		Title:       "Updated Title",
		Description: "Updated Description",
		Template:    "start_stop_continue",
	}

	updatedRetrospective, err := service.UpdateRetrospective(retrospectiveID, userID, request)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access denied")
	assert.Nil(t, updatedRetrospective)
}

func TestRetrospectiveService_DeleteRetrospective(t *testing.T) {
	mockRetroRepo := NewMockRetrospectiveRepository()
	mockTeamRepo := NewMockTeamRepositoryForRetro()
	service := NewRetrospectiveService(mockRetroRepo, mockTeamRepo)

	retrospectiveID := uuid.New()
	userID := uuid.New()

	// Setup retrospective
	retrospective := &models.Retrospective{
		ID:        retrospectiveID,
		Title:     "Test Retrospective",
		CreatedBy: userID,
	}
	mockRetroRepo.retrospectives[retrospectiveID] = retrospective

	err := service.DeleteRetrospective(retrospectiveID, userID)

	assert.NoError(t, err)
}

func TestRetrospectiveService_DeleteRetrospective_AccessDenied(t *testing.T) {
	mockRetroRepo := NewMockRetrospectiveRepository()
	mockTeamRepo := NewMockTeamRepositoryForRetro()
	service := NewRetrospectiveService(mockRetroRepo, mockTeamRepo)
	
	retrospectiveID := uuid.New()
	userID := uuid.New()
	otherUserID := uuid.New()
	
	// Setup retrospective created by other user
	retrospective := &models.Retrospective{
		ID:        retrospectiveID,
		Title:     "Test Retrospective",
		CreatedBy: otherUserID,
	}
	mockRetroRepo.retrospectives[retrospectiveID] = retrospective
	
	err := service.DeleteRetrospective(retrospectiveID, userID)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access denied")
}

func TestRetrospectiveService_ReopenRetrospective(t *testing.T) {
	mockRetroRepo := NewMockRetrospectiveRepository()
	mockTeamRepo := NewMockTeamRepositoryForRetro()
	service := NewRetrospectiveService(mockRetroRepo, mockTeamRepo)
	
	retrospectiveID := uuid.New()
	userID := uuid.New()
	
	// Setup closed retrospective
	retrospective := &models.Retrospective{
		ID:        retrospectiveID,
		Title:     "Test Retrospective",
		Status:    models.RetroStatusClosed,
		CreatedBy: userID,
	}
	mockRetroRepo.retrospectives[retrospectiveID] = retrospective
	
	err := service.ReopenRetrospective(retrospectiveID, userID)
	
	assert.NoError(t, err)
}

func TestRetrospectiveService_ReopenRetrospective_AccessDenied(t *testing.T) {
	mockRetroRepo := NewMockRetrospectiveRepository()
	mockTeamRepo := NewMockTeamRepositoryForRetro()
	service := NewRetrospectiveService(mockRetroRepo, mockTeamRepo)
	
	retrospectiveID := uuid.New()
	userID := uuid.New()
	otherUserID := uuid.New()
	
	// Setup retrospective created by other user
	retrospective := &models.Retrospective{
		ID:        retrospectiveID,
		Title:     "Test Retrospective",
		Status:    models.RetroStatusClosed,
		CreatedBy: otherUserID,
	}
	mockRetroRepo.retrospectives[retrospectiveID] = retrospective
	
	err := service.ReopenRetrospective(retrospectiveID, userID)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access denied")
}

func TestRetrospectiveService_ReopenRetrospective_NotClosed(t *testing.T) {
	mockRetroRepo := NewMockRetrospectiveRepository()
	mockTeamRepo := NewMockTeamRepositoryForRetro()
	service := NewRetrospectiveService(mockRetroRepo, mockTeamRepo)
	
	retrospectiveID := uuid.New()
	userID := uuid.New()
	
	// Setup active retrospective (not closed)
	retrospective := &models.Retrospective{
		ID:        retrospectiveID,
		Title:     "Test Retrospective",
		Status:    models.RetroStatusActive,
		CreatedBy: userID,
	}
	mockRetroRepo.retrospectives[retrospectiveID] = retrospective
	
	err := service.ReopenRetrospective(retrospectiveID, userID)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "retrospective is not closed")
}

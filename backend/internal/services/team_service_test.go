package services

import (
	"database/sql"
	"testing"
	"time"

	"educ-retro/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// MockTeamRepository é um mock simples do TeamRepository
type MockTeamRepository struct {
	teams       map[uuid.UUID]*models.Team
	memberRoles map[string]string // "teamID:userID" -> role
}

func NewMockTeamRepository() *MockTeamRepository {
	return &MockTeamRepository{
		teams:       make(map[uuid.UUID]*models.Team),
		memberRoles: make(map[string]string),
	}
}

func (m *MockTeamRepository) Create(team *models.Team) error {
	if team.ID == uuid.Nil {
		team.ID = uuid.New()
	}
	now := time.Now()
	team.CreatedAt = now
	team.UpdatedAt = now
	m.teams[team.ID] = team
	return nil
}

func (m *MockTeamRepository) GetByID(id uuid.UUID) (*models.Team, error) {
	team, exists := m.teams[id]
	if !exists {
		return nil, sql.ErrNoRows
	}
	teamCopy := *team
	return &teamCopy, nil
}

func (m *MockTeamRepository) GetByUserID(userID uuid.UUID) ([]models.TeamWithCounts, error) {
	var teams []models.TeamWithCounts
	for _, team := range m.teams {
		if team.OwnerID == userID {
			teams = append(teams, models.TeamWithCounts{
				Team:               *team,
				MemberCount:        1,
				RetrospectiveCount: 0,
			})
		}
	}
	return teams, nil
}

func (m *MockTeamRepository) Update(team *models.Team) error {
	if _, exists := m.teams[team.ID]; !exists {
		return sql.ErrNoRows
	}
	team.UpdatedAt = time.Now()
	m.teams[team.ID] = team
	return nil
}

func (m *MockTeamRepository) Delete(id uuid.UUID) error {
	if _, exists := m.teams[id]; !exists {
		return sql.ErrNoRows
	}
	delete(m.teams, id)
	return nil
}

func (m *MockTeamRepository) IsMember(teamID, userID uuid.UUID) (bool, error) {
	key := teamID.String() + ":" + userID.String()
	_, exists := m.memberRoles[key]
	return exists, nil
}

func (m *MockTeamRepository) GetMemberRole(teamID, userID uuid.UUID) (string, error) {
	key := teamID.String() + ":" + userID.String()
	role, exists := m.memberRoles[key]
	if !exists {
		return "", sql.ErrNoRows
	}
	return role, nil
}

func (m *MockTeamRepository) GetMembers(teamID uuid.UUID) ([]models.TeamMemberWithUser, error) {
	return []models.TeamMemberWithUser{}, nil
}

func (m *MockTeamRepository) AddMember(teamID, userID uuid.UUID, role string) error {
	key := teamID.String() + ":" + userID.String()
	m.memberRoles[key] = role
	return nil
}

func (m *MockTeamRepository) RemoveMember(teamID, userID uuid.UUID) error {
	key := teamID.String() + ":" + userID.String()
	delete(m.memberRoles, key)
	return nil
}

func (m *MockTeamRepository) UpdateMemberRole(teamID, userID uuid.UUID, newRole string) error {
	key := teamID.String() + ":" + userID.String()
	if _, exists := m.memberRoles[key]; !exists {
		return sql.ErrNoRows
	}
	m.memberRoles[key] = newRole
	return nil
}

// MockUserRepository para testes
type MockUserRepositoryForTeam struct {
	users  map[uuid.UUID]*models.User
	emails map[string]*models.User
}

func NewMockUserRepositoryForTeam() *MockUserRepositoryForTeam {
	return &MockUserRepositoryForTeam{
		users:  make(map[uuid.UUID]*models.User),
		emails: make(map[string]*models.User),
	}
}

func (m *MockUserRepositoryForTeam) GetByID(id uuid.UUID) (*models.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, sql.ErrNoRows
	}
	userCopy := *user
	return &userCopy, nil
}

func (m *MockUserRepositoryForTeam) GetByEmail(email string) (*models.User, error) {
	user, exists := m.emails[email]
	if !exists {
		return nil, sql.ErrNoRows
	}
	userCopy := *user
	return &userCopy, nil
}

func (m *MockUserRepositoryForTeam) Create(user *models.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	m.users[user.ID] = user
	m.emails[user.Email] = user
	return nil
}

func (m *MockUserRepositoryForTeam) Update(user *models.User) error {
	if _, exists := m.users[user.ID]; !exists {
		return sql.ErrNoRows
	}
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	m.emails[user.Email] = user
	return nil
}

func (m *MockUserRepositoryForTeam) Delete(id uuid.UUID) error {
	user, exists := m.users[id]
	if !exists {
		return sql.ErrNoRows
	}
	delete(m.users, id)
	delete(m.emails, user.Email)
	return nil
}

func (m *MockUserRepositoryForTeam) GetUsersByIDs(ids []uuid.UUID) ([]models.User, error) {
	var users []models.User
	for _, id := range ids {
		if user, exists := m.users[id]; exists {
			users = append(users, *user)
		}
	}
	return users, nil
}

func TestNewTeamService(t *testing.T) {
	mockTeamRepo := NewMockTeamRepository()
	mockUserRepo := NewMockUserRepositoryForTeam()
	service := NewTeamService(mockTeamRepo, mockUserRepo)

	assert.NotNil(t, service)
	assert.Equal(t, mockTeamRepo, service.teamRepo)
	assert.Equal(t, mockUserRepo, service.userRepo)
}

func TestTeamService_CreateTeam(t *testing.T) {
	mockTeamRepo := NewMockTeamRepository()
	mockUserRepo := NewMockUserRepositoryForTeam()
	service := NewTeamService(mockTeamRepo, mockUserRepo)

	userID := uuid.New()
	request := &models.TeamCreateRequest{
		Name:        "Test Team",
		Description: "Test Description",
	}

	team, err := service.CreateTeam(userID, request)

	assert.NoError(t, err)
	assert.NotNil(t, team)
	assert.Equal(t, request.Name, team.Name)
	assert.Equal(t, userID, team.OwnerID)
}

func TestTeamService_GetUserTeams(t *testing.T) {
	mockTeamRepo := NewMockTeamRepository()
	mockUserRepo := NewMockUserRepositoryForTeam()
	service := NewTeamService(mockTeamRepo, mockUserRepo)

	userID := uuid.New()

	// Criar um time para o usuário
	team := &models.Team{
		ID:      uuid.New(),
		Name:    "Test Team",
		OwnerID: userID,
	}
	mockTeamRepo.teams[team.ID] = team

	teams, err := service.GetUserTeams(userID)

	assert.NoError(t, err)
	assert.Len(t, teams, 1)
	assert.Equal(t, team.Name, teams[0].Name)
}

func TestTeamService_GetTeam(t *testing.T) {
	mockTeamRepo := NewMockTeamRepository()
	mockUserRepo := NewMockUserRepositoryForTeam()
	service := NewTeamService(mockTeamRepo, mockUserRepo)

	teamID := uuid.New()
	userID := uuid.New()
	ownerID := uuid.New()

	// Setup team
	team := &models.Team{
		ID:      teamID,
		Name:    "Test Team",
		OwnerID: ownerID,
	}
	mockTeamRepo.teams[teamID] = team

	// Setup member role
	key := teamID.String() + ":" + userID.String()
	mockTeamRepo.memberRoles[key] = "member"

	// Setup owner
	owner := &models.User{
		ID:    ownerID,
		Email: "owner@example.com",
		Name:  "Owner",
	}
	mockUserRepo.users[ownerID] = owner

	teamWithMembers, err := service.GetTeam(teamID, userID)

	assert.NoError(t, err)
	assert.NotNil(t, teamWithMembers)
	assert.Equal(t, team.Name, teamWithMembers.Team.Name)
}

func TestTeamService_GetTeam_AccessDenied(t *testing.T) {
	mockTeamRepo := NewMockTeamRepository()
	mockUserRepo := NewMockUserRepositoryForTeam()
	service := NewTeamService(mockTeamRepo, mockUserRepo)

	teamID := uuid.New()
	userID := uuid.New()

	// Setup team
	team := &models.Team{
		ID:      teamID,
		Name:    "Test Team",
		OwnerID: uuid.New(),
	}
	mockTeamRepo.teams[teamID] = team

	// No member role setup - user is not a member

	teamWithMembers, err := service.GetTeam(teamID, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access denied")
	assert.Nil(t, teamWithMembers)
}

func TestTeamService_UpdateTeam(t *testing.T) {
	mockTeamRepo := NewMockTeamRepository()
	mockUserRepo := NewMockUserRepositoryForTeam()
	service := NewTeamService(mockTeamRepo, mockUserRepo)

	teamID := uuid.New()
	userID := uuid.New()

	// Setup team
	team := &models.Team{
		ID:      teamID,
		Name:    "Original Team",
		OwnerID: userID,
	}
	mockTeamRepo.teams[teamID] = team

	// Setup owner role
	key := teamID.String() + ":" + userID.String()
	mockTeamRepo.memberRoles[key] = "owner"

	request := &models.TeamCreateRequest{
		Name:        "Updated Team",
		Description: "Updated Description",
	}

	updatedTeam, err := service.UpdateTeam(teamID, userID, request)

	assert.NoError(t, err)
	assert.NotNil(t, updatedTeam)
	assert.Equal(t, request.Name, updatedTeam.Name)
}

func TestTeamService_UpdateTeam_AccessDenied(t *testing.T) {
	mockTeamRepo := NewMockTeamRepository()
	mockUserRepo := NewMockUserRepositoryForTeam()
	service := NewTeamService(mockTeamRepo, mockUserRepo)

	teamID := uuid.New()
	userID := uuid.New()

	// Setup team
	team := &models.Team{
		ID:      teamID,
		Name:    "Original Team",
		OwnerID: userID,
	}
	mockTeamRepo.teams[teamID] = team

	// Setup member role (not owner)
	key := teamID.String() + ":" + userID.String()
	mockTeamRepo.memberRoles[key] = "member"

	request := &models.TeamCreateRequest{
		Name:        "Updated Team",
		Description: "Updated Description",
	}

	updatedTeam, err := service.UpdateTeam(teamID, userID, request)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "only team owner can update team")
	assert.Nil(t, updatedTeam)
}

func TestTeamService_DeleteTeam(t *testing.T) {
	mockTeamRepo := NewMockTeamRepository()
	mockUserRepo := NewMockUserRepositoryForTeam()
	service := NewTeamService(mockTeamRepo, mockUserRepo)

	teamID := uuid.New()
	userID := uuid.New()

	// Setup team
	team := &models.Team{
		ID:      teamID,
		Name:    "Test Team",
		OwnerID: userID,
	}
	mockTeamRepo.teams[teamID] = team

	// Setup owner role
	key := teamID.String() + ":" + userID.String()
	mockTeamRepo.memberRoles[key] = "owner"

	err := service.DeleteTeam(teamID, userID)

	assert.NoError(t, err)
}

func TestTeamService_DeleteTeam_AccessDenied(t *testing.T) {
	mockTeamRepo := NewMockTeamRepository()
	mockUserRepo := NewMockUserRepositoryForTeam()
	service := NewTeamService(mockTeamRepo, mockUserRepo)

	teamID := uuid.New()
	userID := uuid.New()

	// Setup team
	team := &models.Team{
		ID:      teamID,
		Name:    "Test Team",
		OwnerID: userID,
	}
	mockTeamRepo.teams[teamID] = team

	// Setup member role (not owner)
	key := teamID.String() + ":" + userID.String()
	mockTeamRepo.memberRoles[key] = "member"

	err := service.DeleteTeam(teamID, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "only team owner can delete team")
}

func TestTeamService_AddMember(t *testing.T) {
	mockTeamRepo := NewMockTeamRepository()
	mockUserRepo := NewMockUserRepositoryForTeam()
	service := NewTeamService(mockTeamRepo, mockUserRepo)

	teamID := uuid.New()
	userID := uuid.New()
	newUserID := uuid.New()

	// Setup current user as member
	key := teamID.String() + ":" + userID.String()
	mockTeamRepo.memberRoles[key] = "member"

	// Setup new user
	newUser := &models.User{
		ID:    newUserID,
		Email: "newmember@example.com",
		Name:  "New Member",
	}
	mockUserRepo.emails["newmember@example.com"] = newUser

	request := &models.TeamInviteRequest{
		Email: "newmember@example.com",
		Role:  "member",
	}

	err := service.AddMember(teamID, userID, request)

	assert.NoError(t, err)
}

func TestTeamService_AddMember_UserNotFound(t *testing.T) {
	mockTeamRepo := NewMockTeamRepository()
	mockUserRepo := NewMockUserRepositoryForTeam()
	service := NewTeamService(mockTeamRepo, mockUserRepo)

	teamID := uuid.New()
	userID := uuid.New()

	// Setup current user as member
	key := teamID.String() + ":" + userID.String()
	mockTeamRepo.memberRoles[key] = "member"

	request := &models.TeamInviteRequest{
		Email: "nonexistent@example.com",
		Role:  "member",
	}

	err := service.AddMember(teamID, userID, request)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestTeamService_AddMember_InsufficientPermissions(t *testing.T) {
	mockTeamRepo := NewMockTeamRepository()
	mockUserRepo := NewMockUserRepositoryForTeam()
	service := NewTeamService(mockTeamRepo, mockUserRepo)

	teamID := uuid.New()
	userID := uuid.New()

	// Setup user with no permissions (viewer role)
	key := teamID.String() + ":" + userID.String()
	mockTeamRepo.memberRoles[key] = "viewer"

	request := &models.TeamInviteRequest{
		Email: "newmember@example.com",
		Role:  "member",
	}

	err := service.AddMember(teamID, userID, request)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient permissions")
}

func TestTeamService_RemoveMember(t *testing.T) {
	mockTeamRepo := NewMockTeamRepository()
	mockUserRepo := NewMockUserRepositoryForTeam()
	service := NewTeamService(mockTeamRepo, mockUserRepo)

	teamID := uuid.New()
	userID := uuid.New()
	targetUserID := uuid.New()

	// Setup owner role
	key := teamID.String() + ":" + userID.String()
	mockTeamRepo.memberRoles[key] = "owner"

	// Setup target member
	targetKey := teamID.String() + ":" + targetUserID.String()
	mockTeamRepo.memberRoles[targetKey] = "member"

	err := service.RemoveMember(teamID, userID, targetUserID)

	assert.NoError(t, err)
}

func TestTeamService_RemoveMember_InsufficientPermissions(t *testing.T) {
	mockTeamRepo := NewMockTeamRepository()
	mockUserRepo := NewMockUserRepositoryForTeam()
	service := NewTeamService(mockTeamRepo, mockUserRepo)

	teamID := uuid.New()
	userID := uuid.New()
	targetUserID := uuid.New()

	// Setup member role
	key := teamID.String() + ":" + userID.String()
	mockTeamRepo.memberRoles[key] = "member"

	err := service.RemoveMember(teamID, userID, targetUserID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient permissions")
}

func TestTeamService_UpdateMemberRole(t *testing.T) {
	mockTeamRepo := NewMockTeamRepository()
	mockUserRepo := NewMockUserRepositoryForTeam()
	service := NewTeamService(mockTeamRepo, mockUserRepo)

	teamID := uuid.New()
	userID := uuid.New()
	targetUserID := uuid.New()

	// Setup owner role
	key := teamID.String() + ":" + userID.String()
	mockTeamRepo.memberRoles[key] = "owner"

	// Setup target member
	targetKey := teamID.String() + ":" + targetUserID.String()
	mockTeamRepo.memberRoles[targetKey] = "member"

	err := service.UpdateMemberRole(teamID, userID, targetUserID, "viewer")

	assert.NoError(t, err)
}

func TestTeamService_UpdateMemberRole_AccessDenied(t *testing.T) {
	mockTeamRepo := NewMockTeamRepository()
	mockUserRepo := NewMockUserRepositoryForTeam()
	service := NewTeamService(mockTeamRepo, mockUserRepo)

	teamID := uuid.New()
	userID := uuid.New()
	targetUserID := uuid.New()

	// Setup member role (not owner)
	key := teamID.String() + ":" + userID.String()
	mockTeamRepo.memberRoles[key] = "member"

	err := service.UpdateMemberRole(teamID, userID, targetUserID, "viewer")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "only team owner can update member roles")
}

func TestTeamService_UpdateMemberRole_CannotChangeToOwner(t *testing.T) {
	mockTeamRepo := NewMockTeamRepository()
	mockUserRepo := NewMockUserRepositoryForTeam()
	service := NewTeamService(mockTeamRepo, mockUserRepo)

	teamID := uuid.New()
	userID := uuid.New()
	targetUserID := uuid.New()

	// Setup owner role
	key := teamID.String() + ":" + userID.String()
	mockTeamRepo.memberRoles[key] = "owner"

	err := service.UpdateMemberRole(teamID, userID, targetUserID, "owner")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot change role to owner")
}

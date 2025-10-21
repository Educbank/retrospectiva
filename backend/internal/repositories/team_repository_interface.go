package repositories

import (
	"educ-retro/internal/models"
	"github.com/google/uuid"
)

// TeamRepositoryInterface define a interface para o TeamRepository
type TeamRepositoryInterface interface {
	Create(team *models.Team) error
	GetByID(id uuid.UUID) (*models.Team, error)
	GetByUserID(userID uuid.UUID) ([]models.TeamWithCounts, error)
	Update(team *models.Team) error
	Delete(id uuid.UUID) error
	IsMember(teamID, userID uuid.UUID) (bool, error)
	GetMemberRole(teamID, userID uuid.UUID) (string, error)
	GetMembers(teamID uuid.UUID) ([]models.TeamMemberWithUser, error)
	AddMember(teamID, userID uuid.UUID, role string) error
	RemoveMember(teamID, userID uuid.UUID) error
	UpdateMemberRole(teamID, userID uuid.UUID, newRole string) error
}

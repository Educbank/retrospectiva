package repositories

import (
	"educ-retro/internal/models"

	"github.com/google/uuid"
)

// RetrospectiveRepositoryInterface define a interface para o RetrospectiveRepository
type RetrospectiveRepositoryInterface interface {
	Create(retrospective *models.Retrospective) error
	GetByID(id uuid.UUID) (*models.Retrospective, error)
	GetAllRetrospectives() ([]models.Retrospective, error)
	GetRetrospectiveWithDetails(id uuid.UUID) (*models.RetrospectiveWithDetails, error)
	Update(retrospective *models.Retrospective) error
	Delete(id uuid.UUID) error
	UpdateStatus(id uuid.UUID, status models.RetrospectiveStatus) error
	GetRetrospectiveCount(id uuid.UUID) (int, error)
	GetActionItemCount(id uuid.UUID) (int, error)
	AddItem(item *models.RetrospectiveItem) error
	VoteItem(itemID, userID uuid.UUID) error
	AddActionItem(actionItem *models.ActionItem) error
	RegisterParticipant(retrospectiveID, userID uuid.UUID) error
	GetParticipants(retrospectiveID uuid.UUID) ([]models.RetrospectiveParticipant, error)
	GetItemByID(id uuid.UUID) (*models.RetrospectiveItem, error)
	DeleteItem(id uuid.UUID) error
	ReopenRetrospective(id uuid.UUID) error
	CreateGroup(group *models.RetrospectiveGroup, itemIDs []uuid.UUID) error
	VoteGroup(groupID, userID uuid.UUID) error
	GetGroupByID(id uuid.UUID) (*models.RetrospectiveGroup, error)
	DeleteGroup(id uuid.UUID) error
	MergeItems(sourceItemID, targetItemID uuid.UUID) (*models.RetrospectiveItem, error)
	GetActionItemByID(id uuid.UUID) (*models.ActionItem, error)
	UpdateActionItem(actionItemID uuid.UUID, req *models.ActionItemUpdateRequest) (*models.ActionItem, error)
	DeleteActionItem(id uuid.UUID) error
	UpdateTimer(id uuid.UUID, req models.TimerUpdateRequest) error
}

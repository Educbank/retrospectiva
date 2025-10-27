package services

import (
	"errors"
	"time"

	"educ-retro/internal/models"
	"educ-retro/internal/repositories"

	"github.com/google/uuid"
)

type RetrospectiveService struct {
	retroRepo repositories.RetrospectiveRepositoryInterface
}

func NewRetrospectiveService(retroRepo repositories.RetrospectiveRepositoryInterface) *RetrospectiveService {
	return &RetrospectiveService{
		retroRepo: retroRepo,
	}
}

func (s *RetrospectiveService) CreateRetrospective(userID uuid.UUID, req *models.RetrospectiveCreateRequest) (*models.Retrospective, error) {
	retrospective := &models.Retrospective{
		Title:       req.Title,
		Description: &req.Description,
		Template:    req.Template,
		Status:      models.RetroStatusPlanned,
		CreatedBy:   userID,
	}

	err := s.retroRepo.Create(retrospective)
	if err != nil {
		return nil, err
	}

	return retrospective, nil
}

func (s *RetrospectiveService) GetUserRetrospectives(userID uuid.UUID) ([]models.RetrospectiveWithDetails, error) {
	// Get all retrospectives with full details including action items
	retrospectives, err := s.retroRepo.GetAllRetrospectives()
	if err != nil {
		return nil, err
	}

	// Filter retrospectives: show "planned" only to creator, others to everyone
	var filteredRetrospectives []models.Retrospective
	for _, retro := range retrospectives {
		// Show retrospectives that are not "planned" OR are "planned" and created by the current user
		if retro.Status != models.RetroStatusPlanned || retro.CreatedBy == userID {
			filteredRetrospectives = append(filteredRetrospectives, retro)
		}
	}

	// Convert to RetrospectiveWithDetails
	var retrospectivesWithDetails []models.RetrospectiveWithDetails
	for _, retro := range filteredRetrospectives {
		details, err := s.retroRepo.GetRetrospectiveWithDetails(retro.ID)
		if err != nil {
			// If we can't get details for one, continue with others
			continue
		}
		retrospectivesWithDetails = append(retrospectivesWithDetails, *details)
	}

	return retrospectivesWithDetails, nil
}

func (s *RetrospectiveService) GetRetrospective(retrospectiveID, userID uuid.UUID) (*models.Retrospective, error) {
	retrospective, err := s.retroRepo.GetByID(retrospectiveID)
	if err != nil {
		return nil, err
	}

	// Check if user is the creator
	if retrospective.CreatedBy != userID {
		return nil, errors.New("access denied")
	}

	return retrospective, nil
}

func (s *RetrospectiveService) UpdateRetrospective(retrospectiveID, userID uuid.UUID, req *models.RetrospectiveCreateRequest) (*models.Retrospective, error) {
	retrospective, err := s.retroRepo.GetByID(retrospectiveID)
	if err != nil {
		return nil, err
	}

	// Check if user is the creator
	if retrospective.CreatedBy != userID {
		return nil, errors.New("access denied")
	}

	// Update fields
	retrospective.Title = req.Title
	retrospective.Description = &req.Description
	retrospective.Template = req.Template

	err = s.retroRepo.Update(retrospective)
	if err != nil {
		return nil, err
	}

	return retrospective, nil
}

func (s *RetrospectiveService) DeleteRetrospective(retrospectiveID, userID uuid.UUID) error {
	retrospective, err := s.retroRepo.GetByID(retrospectiveID)
	if err != nil {
		return err
	}

	// Check if user is the creator
	if retrospective.CreatedBy != userID {
		return errors.New("access denied")
	}

	// Check if retrospective is closed
	if retrospective.Status == models.RetroStatusClosed {
		return errors.New("retrospective is closed")
	}

	return s.retroRepo.Delete(retrospectiveID)
}

func (s *RetrospectiveService) StartRetrospective(retrospectiveID, userID uuid.UUID) error {
	retrospective, err := s.retroRepo.GetByID(retrospectiveID)
	if err != nil {
		return err
	}

	// Check if user is the creator
	if retrospective.CreatedBy != userID {
		return errors.New("access denied")
	}

	return s.retroRepo.UpdateStatus(retrospectiveID, models.RetroStatusActive)
}

func (s *RetrospectiveService) EndRetrospective(retrospectiveID, userID uuid.UUID) error {
	retrospective, err := s.retroRepo.GetByID(retrospectiveID)
	if err != nil {
		return err
	}

	// Check if user is the creator
	if retrospective.CreatedBy != userID {
		return errors.New("access denied")
	}

	return s.retroRepo.UpdateStatus(retrospectiveID, models.RetroStatusClosed)
}

func (s *RetrospectiveService) GetRetrospectiveStats(userID uuid.UUID) (map[string]int, error) {
	retroCount, err := s.retroRepo.GetRetrospectiveCount(userID)
	if err != nil {
		return nil, err
	}

	actionCount, err := s.retroRepo.GetActionItemCount(userID)
	if err != nil {
		return nil, err
	}

	return map[string]int{
		"retrospectives": retroCount,
		"action_items":   actionCount,
	}, nil
}

func (s *RetrospectiveService) AddItem(retrospectiveID, userID uuid.UUID, req *models.RetrospectiveItemCreateRequest) (*models.RetrospectiveItem, error) {
	// Allow any authenticated user to add items
	item := &models.RetrospectiveItem{
		ID:              uuid.New(),
		RetrospectiveID: retrospectiveID,
		Category:        req.Category,
		Content:         req.Content,
		AuthorID:        &userID,
		IsAnonymous:     req.IsAnonymous,
		Votes:           0,
	}

	if req.IsAnonymous {
		item.AuthorID = nil
	}

	err := s.retroRepo.AddItem(item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *RetrospectiveService) VoteItem(itemID, userID uuid.UUID) error {
	return s.retroRepo.VoteItem(itemID, userID)
}

func (s *RetrospectiveService) AddActionItem(retrospectiveID, userID uuid.UUID, req *models.ActionItemCreateRequest) (*models.ActionItem, error) {
	// Allow any authenticated user to add action items
	actionItem := &models.ActionItem{
		ID:              uuid.New(),
		RetrospectiveID: retrospectiveID,
		Title:           req.Title,
		Description:     &req.Description,
		Status:          "todo",
		CreatedBy:       userID,
	}

	if req.ItemID != nil {
		itemID, err := uuid.Parse(*req.ItemID)
		if err != nil {
			return nil, errors.New("invalid item_id")
		}
		actionItem.ItemID = &itemID
	}

	if req.AssignedTo != nil {
		assignedToID, err := uuid.Parse(*req.AssignedTo)
		if err != nil {
			return nil, errors.New("invalid assigned_to")
		}
		actionItem.AssignedTo = &assignedToID
	}

	if req.DueDate != nil && *req.DueDate != "" {
		dueDate, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			return nil, errors.New("invalid due_date format")
		}
		actionItem.DueDate = &dueDate
	}

	err := s.retroRepo.AddActionItem(actionItem)
	if err != nil {
		return nil, err
	}

	return actionItem, nil
}

func (s *RetrospectiveService) GetRetrospectiveWithDetails(retrospectiveID, userID uuid.UUID) (*models.RetrospectiveWithDetails, error) {
	// Allow all users to view retrospectives - no access control needed
	return s.retroRepo.GetRetrospectiveWithDetails(retrospectiveID)
}

func (s *RetrospectiveService) RegisterParticipant(retrospectiveID, userID uuid.UUID) error {
	// First, register the participant
	err := s.retroRepo.RegisterParticipant(retrospectiveID, userID)
	if err != nil {
		return err
	}

	// Check if retrospective should be started automatically
	retrospective, err := s.retroRepo.GetByID(retrospectiveID)
	if err != nil {
		return err
	}

	// Only auto-start if retrospective is in "planned" status
	if retrospective.Status == models.RetroStatusPlanned {
		// Get current participant count
		participants, err := s.retroRepo.GetParticipants(retrospectiveID)
		if err != nil {
			return err
		}

		// If there's at least 2 participants, start the retrospective automatically
		if len(participants) > 1 {
			err = s.retroRepo.UpdateStatus(retrospectiveID, models.RetroStatusActive)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *RetrospectiveService) GetParticipants(retrospectiveID uuid.UUID) ([]models.RetrospectiveParticipant, error) {
	return s.retroRepo.GetParticipants(retrospectiveID)
}

func (s *RetrospectiveService) GetItemByID(itemID uuid.UUID) (*models.RetrospectiveItem, error) {
	return s.retroRepo.GetItemByID(itemID)
}

func (s *RetrospectiveService) DeleteItem(itemID uuid.UUID) error {
	return s.retroRepo.DeleteItem(itemID)
}

func (s *RetrospectiveService) ReopenRetrospective(retrospectiveID, userID uuid.UUID) error {
	retrospective, err := s.retroRepo.GetByID(retrospectiveID)
	if err != nil {
		return err
	}

	// Check if user is the creator
	if retrospective.CreatedBy != userID {
		return errors.New("access denied")
	}

	// Only allow reopening if retrospective is closed
	if retrospective.Status != models.RetroStatusClosed {
		return errors.New("retrospective is not closed")
	}

	return s.retroRepo.ReopenRetrospective(retrospectiveID)
}

// Group methods
func (s *RetrospectiveService) CreateGroup(retrospectiveID, userID uuid.UUID, req *models.GroupCreateRequest) (*models.RetrospectiveGroup, error) {
	// Verify retrospective exists and user has access
	retrospective, err := s.retroRepo.GetByID(retrospectiveID)
	if err != nil {
		return nil, err
	}

	// Only allow grouping if retrospective is active
	if retrospective.Status != models.RetroStatusActive {
		return nil, errors.New("can only create groups for active retrospectives")
	}

	// Parse item IDs
	var itemIDs []uuid.UUID
	for _, itemIDStr := range req.ItemIDs {
		itemID, err := uuid.Parse(itemIDStr)
		if err != nil {
			return nil, errors.New("invalid item ID: " + itemIDStr)
		}
		itemIDs = append(itemIDs, itemID)
	}

	// Verify all items exist and belong to the retrospective
	for _, itemID := range itemIDs {
		item, err := s.retroRepo.GetItemByID(itemID)
		if err != nil {
			return nil, errors.New("item not found: " + itemID.String())
		}
		if item.RetrospectiveID != retrospectiveID {
			return nil, errors.New("item does not belong to this retrospective")
		}
	}

	group := &models.RetrospectiveGroup{
		ID:              uuid.New(),
		RetrospectiveID: retrospectiveID,
		Name:            req.Name,
		Description:     req.Description,
		Votes:           0,
		CreatedBy:       userID,
	}

	err = s.retroRepo.CreateGroup(group, itemIDs)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (s *RetrospectiveService) VoteGroup(groupID, userID uuid.UUID) error {
	return s.retroRepo.VoteGroup(groupID, userID)
}

func (s *RetrospectiveService) GetGroupByID(groupID uuid.UUID) (*models.RetrospectiveGroup, error) {
	return s.retroRepo.GetGroupByID(groupID)
}

func (s *RetrospectiveService) DeleteGroup(groupID, userID uuid.UUID) error {
	// First get the group to check permissions - we need to find which retrospective it belongs to
	// Since we don't have a direct GetGroupByID method, we'll add one to the repository
	group, err := s.retroRepo.GetGroupByID(groupID)
	if err != nil {
		return err
	}

	// Check if user is the creator
	if group.CreatedBy != userID {
		return errors.New("access denied")
	}

	return s.retroRepo.DeleteGroup(groupID)
}

func (s *RetrospectiveService) MergeItems(sourceItemID, targetItemID, userID uuid.UUID) (*models.RetrospectiveItem, error) {
	// Get source item to verify permissions and retrospective
	sourceItem, err := s.retroRepo.GetItemByID(sourceItemID)
	if err != nil {
		return nil, err
	}

	// Get target item to verify permissions and retrospective
	targetItem, err := s.retroRepo.GetItemByID(targetItemID)
	if err != nil {
		return nil, err
	}

	// Verify both items belong to the same retrospective
	if sourceItem.RetrospectiveID != targetItem.RetrospectiveID {
		return nil, errors.New("items must belong to the same retrospective")
	}

	// Verify retrospective is active (can only merge items in active retrospectives)
	retrospective, err := s.retroRepo.GetByID(sourceItem.RetrospectiveID)
	if err != nil {
		return nil, err
	}

	if retrospective.Status != models.RetroStatusActive {
		return nil, errors.New("can only merge items in active retrospectives")
	}

	// Merge items
	return s.retroRepo.MergeItems(sourceItemID, targetItemID)
}

// Action Item methods
func (s *RetrospectiveService) GetActionItemByID(actionItemID uuid.UUID) (*models.ActionItem, error) {
	return s.retroRepo.GetActionItemByID(actionItemID)
}

func (s *RetrospectiveService) UpdateActionItem(actionItemID, userID uuid.UUID, req *models.ActionItemUpdateRequest) (*models.ActionItem, error) {
	// Get the action item to check permissions
	actionItem, err := s.retroRepo.GetActionItemByID(actionItemID)
	if err != nil {
		return nil, err
	}

	// Check if user is the creator or has access to the retrospective
	retrospective, err := s.retroRepo.GetByID(actionItem.RetrospectiveID)
	if err != nil {
		return nil, err
	}

	// Allow creator of action item or creator of retrospective to update
	if actionItem.CreatedBy != userID && retrospective.CreatedBy != userID {
		return nil, errors.New("access denied")
	}

	// Validate status if provided
	if req.Status != nil {
		validStatuses := []string{"todo", "in_progress", "done"}
		isValid := false
		for _, status := range validStatuses {
			if *req.Status == status {
				isValid = true
				break
			}
		}
		if !isValid {
			return nil, errors.New("invalid status. Must be one of: todo, in_progress, done")
		}

		// If status is being changed to "done" and completed_at is not provided, set it to now
		if *req.Status == "done" && req.CompletedAt == nil {
			now := time.Now().Format("2006-01-02T15:04:05Z")
			req.CompletedAt = &now
		}

		// If status is being changed from "done" to something else, clear completed_at
		if *req.Status != "done" && actionItem.Status == "done" {
			empty := ""
			req.CompletedAt = &empty
		}
	}

	return s.retroRepo.UpdateActionItem(actionItemID, req)
}

func (s *RetrospectiveService) DeleteActionItem(actionItemID, userID uuid.UUID) error {
	// Get the action item to check permissions
	actionItem, err := s.retroRepo.GetActionItemByID(actionItemID)
	if err != nil {
		return err
	}

	// Check if user is the creator or has access to the retrospective
	retrospective, err := s.retroRepo.GetByID(actionItem.RetrospectiveID)
	if err != nil {
		return err
	}

	// Allow creator of action item or creator of retrospective to delete
	if actionItem.CreatedBy != userID && retrospective.CreatedBy != userID {
		return errors.New("access denied")
	}

	return s.retroRepo.DeleteActionItem(actionItemID)
}

package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"educ-retro/internal/models"
	"educ-retro/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
)

type RetrospectiveHandler struct {
	retrospectiveService *services.RetrospectiveService
	realtimeService      *services.RealtimeService
}

func NewRetrospectiveHandler(retrospectiveService *services.RetrospectiveService, realtimeService *services.RealtimeService) *RetrospectiveHandler {
	return &RetrospectiveHandler{
		retrospectiveService: retrospectiveService,
		realtimeService:      realtimeService,
	}
}

// CreateRetrospective godoc
// @Summary Create a new retrospective
// @Description Create a new retrospective for a team
// @Tags Retrospectives
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param retrospective body models.RetrospectiveCreateRequest true "Retrospective creation data"
// @Success 201 {object} models.Retrospective "Retrospective created successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied"
// @Router /retrospectives [post]
func (h *RetrospectiveHandler) CreateRetrospective(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req models.RetrospectiveCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	retrospective, err := h.retrospectiveService.CreateRetrospective(userID.(uuid.UUID), &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "access denied" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, retrospective)
}

// GetUserRetrospectives godoc
// @Summary Get user's retrospectives
// @Description Get all retrospectives where the current user is a team member
// @Tags Retrospectives
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Retrospective "User's retrospectives"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /retrospectives [get]
func (h *RetrospectiveHandler) GetUserRetrospectives(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	retrospectives, err := h.retrospectiveService.GetUserRetrospectives(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, retrospectives)
}

// GetRetrospective godoc
// @Summary Get retrospective details
// @Description Get detailed information about a specific retrospective
// @Tags Retrospectives
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Retrospective ID"
// @Success 200 {object} models.Retrospective "Retrospective details"
// @Failure 400 {object} map[string]string "Invalid retrospective ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied"
// @Failure 404 {object} map[string]string "Retrospective not found"
// @Router /retrospectives/{id} [get]
func (h *RetrospectiveHandler) GetRetrospective(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	retrospectiveIDStr := c.Param("id")
	retrospectiveID, err := uuid.Parse(retrospectiveIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid retrospective ID"})
		return
	}

	retrospective, err := h.retrospectiveService.GetRetrospectiveWithDetails(retrospectiveID, userID.(uuid.UUID))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "access denied" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, retrospective)
}

// UpdateRetrospective godoc
// @Summary Update retrospective
// @Description Update retrospective information
// @Tags Retrospectives
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Retrospective ID"
// @Param retrospective body models.RetrospectiveCreateRequest true "Retrospective update data"
// @Success 200 {object} models.Retrospective "Updated retrospective"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied"
// @Router /retrospectives/{id} [put]
func (h *RetrospectiveHandler) UpdateRetrospective(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	retrospectiveIDStr := c.Param("id")
	retrospectiveID, err := uuid.Parse(retrospectiveIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid retrospective ID"})
		return
	}

	var req models.RetrospectiveCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	retrospective, err := h.retrospectiveService.UpdateRetrospective(retrospectiveID, userID.(uuid.UUID), &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "access denied" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, retrospective)
}

// DeleteRetrospective godoc
// @Summary Delete retrospective
// @Description Delete a retrospective
// @Tags Retrospectives
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Retrospective ID"
// @Success 204 "Retrospective deleted successfully"
// @Failure 400 {object} map[string]string "Invalid retrospective ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied"
// @Router /retrospectives/{id} [delete]
func (h *RetrospectiveHandler) DeleteRetrospective(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	retrospectiveIDStr := c.Param("id")
	retrospectiveID, err := uuid.Parse(retrospectiveIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid retrospective ID"})
		return
	}

	err = h.retrospectiveService.DeleteRetrospective(retrospectiveID, userID.(uuid.UUID))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "access denied" {
			status = http.StatusForbidden
		} else if err.Error() == "retrospective is closed" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// StartRetrospective godoc
// @Summary Start retrospective
// @Description Start a retrospective session
// @Tags Retrospectives
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Retrospective ID"
// @Success 200 {object} map[string]string "Retrospective started successfully"
// @Failure 400 {object} map[string]string "Invalid retrospective ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied"
// @Router /retrospectives/{id}/start [post]
func (h *RetrospectiveHandler) StartRetrospective(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	retrospectiveIDStr := c.Param("id")
	retrospectiveID, err := uuid.Parse(retrospectiveIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid retrospective ID"})
		return
	}

	err = h.retrospectiveService.StartRetrospective(retrospectiveID, userID.(uuid.UUID))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "access denied" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Retrospective started successfully"})
}

// EndRetrospective godoc
// @Summary End retrospective
// @Description End a retrospective session
// @Tags Retrospectives
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Retrospective ID"
// @Success 200 {object} map[string]string "Retrospective ended successfully"
// @Failure 400 {object} map[string]string "Invalid retrospective ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied"
// @Router /retrospectives/{id}/end [post]
func (h *RetrospectiveHandler) EndRetrospective(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	retrospectiveIDStr := c.Param("id")
	retrospectiveID, err := uuid.Parse(retrospectiveIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid retrospective ID"})
		return
	}

	err = h.retrospectiveService.EndRetrospective(retrospectiveID, userID.(uuid.UUID))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "access denied" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Retrospective ended successfully"})
}

func (h *RetrospectiveHandler) ReopenRetrospective(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	retrospectiveIDStr := c.Param("id")
	retrospectiveID, err := uuid.Parse(retrospectiveIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid retrospective ID"})
		return
	}

	err = h.retrospectiveService.ReopenRetrospective(retrospectiveID, userID.(uuid.UUID))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "access denied" {
			status = http.StatusForbidden
		} else if err.Error() == "retrospective is not closed" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Retrospective reopened successfully"})
}

func (h *RetrospectiveHandler) AddItem(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	retrospectiveIDStr := c.Param("id")
	retrospectiveID, err := uuid.Parse(retrospectiveIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid retrospective ID"})
		return
	}

	var req models.RetrospectiveItemCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.retrospectiveService.AddItem(retrospectiveID, userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send real-time update via SSE
	if h.realtimeService != nil {
		h.realtimeService.BroadcastToRetrospective(retrospectiveID, "item_added", map[string]interface{}{
			"item": item,
		})
	}

	c.JSON(http.StatusCreated, item)
}

func (h *RetrospectiveHandler) VoteItem(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	itemIDStr := c.Param("itemId")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	// Get the item to find the retrospective ID
	item, err := h.retrospectiveService.GetItemByID(itemID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	// Get the retrospective to check its status
	retrospective, err := h.retrospectiveService.GetRetrospectiveWithDetails(item.RetrospectiveID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Retrospective not found"})
		return
	}

	// Check if retrospective is closed
	if retrospective.Status == "closed" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot vote on items in a closed retrospective"})
		return
	}

	err = h.retrospectiveService.VoteItem(itemID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send real-time update via SSE
	if h.realtimeService != nil {
		// Get the updated item to send complete data
		item, err := h.retrospectiveService.GetItemByID(itemID)
		if err == nil {
			h.realtimeService.BroadcastToRetrospective(item.RetrospectiveID, "item_voted", map[string]interface{}{
				"item": item,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote recorded successfully"})
}

func (h *RetrospectiveHandler) DeleteItem(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	itemIDStr := c.Param("itemId")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	// Get item to check retrospective and permissions
	item, err := h.retrospectiveService.GetItemByID(itemID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	// Check if user can delete item (must be the author or retrospective creator)
	retrospective, err := h.retrospectiveService.GetRetrospective(item.RetrospectiveID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Retrospective not found"})
		return
	}

	canDelete := false
	if retrospective.CreatedBy == userID.(uuid.UUID) {
		canDelete = true // Creator can delete any item
	} else if item.AuthorID != nil && *item.AuthorID == userID.(uuid.UUID) {
		canDelete = true // Author can delete their own item
	}

	if !canDelete {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own items or be the retrospective creator"})
		return
	}

	err = h.retrospectiveService.DeleteItem(itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send real-time update via SSE
	if h.realtimeService != nil {
		h.realtimeService.BroadcastToRetrospective(item.RetrospectiveID, "item_deleted", map[string]interface{}{
			"item_id": itemID,
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}

func (h *RetrospectiveHandler) AddActionItem(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	retrospectiveIDStr := c.Param("id")
	retrospectiveID, err := uuid.Parse(retrospectiveIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid retrospective ID"})
		return
	}

	var req models.ActionItemCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actionItem, err := h.retrospectiveService.AddActionItem(retrospectiveID, userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send real-time update via SSE
	if h.realtimeService != nil {
		h.realtimeService.BroadcastToRetrospective(retrospectiveID, "action_item_added", map[string]interface{}{
			"action_item": actionItem,
		})
	}

	c.JSON(http.StatusCreated, actionItem)
}

func (h *RetrospectiveHandler) UpdateActionItem(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	actionItemIDStr := c.Param("actionItemId")
	actionItemID, err := uuid.Parse(actionItemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action item ID"})
		return
	}

	var req models.ActionItemUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actionItem, err := h.retrospectiveService.UpdateActionItem(actionItemID, userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send real-time update via SSE
	if h.realtimeService != nil {
		h.realtimeService.BroadcastToRetrospective(actionItem.RetrospectiveID, "action_item_updated", map[string]interface{}{
			"action_item": actionItem,
		})
	}

	c.JSON(http.StatusOK, actionItem)
}

func (h *RetrospectiveHandler) DeleteActionItem(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	actionItemIDStr := c.Param("actionItemId")
	actionItemID, err := uuid.Parse(actionItemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action item ID"})
		return
	}

	// Get action item before deleting to get retrospective ID for SSE
	actionItem, err := h.retrospectiveService.GetActionItemByID(actionItemID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Action item not found"})
		return
	}

	err = h.retrospectiveService.DeleteActionItem(actionItemID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send real-time update via SSE
	if h.realtimeService != nil {
		h.realtimeService.BroadcastToRetrospective(actionItem.RetrospectiveID, "action_item_deleted", map[string]interface{}{
			"action_item_id": actionItemID,
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Action item deleted successfully"})
}

func (h *RetrospectiveHandler) JoinRetrospective(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	retrospectiveIDStr := c.Param("id")
	retrospectiveID, err := uuid.Parse(retrospectiveIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid retrospective ID"})
		return
	}

	// Register user access to the retrospective
	err = h.retrospectiveService.RegisterParticipant(retrospectiveID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Joined retrospective successfully"})
}

func (h *RetrospectiveHandler) GetParticipants(c *gin.Context) {
	retrospectiveIDStr := c.Param("id")
	retrospectiveID, err := uuid.Parse(retrospectiveIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid retrospective ID"})
		return
	}

	participants, err := h.retrospectiveService.GetParticipants(retrospectiveID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"participants": participants})
}

func (h *RetrospectiveHandler) CreateGroup(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	retrospectiveIDStr := c.Param("id")
	retrospectiveID, err := uuid.Parse(retrospectiveIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid retrospective ID"})
		return
	}

	var req models.GroupCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group, err := h.retrospectiveService.CreateGroup(retrospectiveID, userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send real-time update via SSE
	if h.realtimeService != nil {
		h.realtimeService.BroadcastToRetrospective(retrospectiveID, "group_created", map[string]interface{}{
			"group": group,
		})
	}

	c.JSON(http.StatusCreated, group)
}

func (h *RetrospectiveHandler) VoteGroup(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	groupIDStr := c.Param("groupId")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// Get the group to find the retrospective ID
	group, err := h.retrospectiveService.GetGroupByID(groupID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	// Get the retrospective to check its status
	retrospective, err := h.retrospectiveService.GetRetrospectiveWithDetails(group.RetrospectiveID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Retrospective not found"})
		return
	}

	// Check if retrospective is closed
	if retrospective.Status == "closed" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot vote on groups in a closed retrospective"})
		return
	}

	err = h.retrospectiveService.VoteGroup(groupID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the group to find the retrospective ID for broadcasting
	group, err = h.retrospectiveService.GetGroupByID(groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send real-time update via SSE
	if h.realtimeService != nil {
		h.realtimeService.BroadcastToRetrospective(group.RetrospectiveID, "group_voted", map[string]interface{}{
			"group_id": groupID,
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote recorded successfully"})
}

func (h *RetrospectiveHandler) DeleteGroup(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	groupIDStr := c.Param("groupId")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// Get the group first to find the retrospective ID for broadcasting
	group, err := h.retrospectiveService.GetGroupByID(groupID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	err = h.retrospectiveService.DeleteGroup(groupID, userID.(uuid.UUID))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "access denied" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// Send real-time update via SSE
	if h.realtimeService != nil {
		h.realtimeService.BroadcastToRetrospective(group.RetrospectiveID, "group_deleted", map[string]interface{}{
			"group_id": groupID,
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group deleted successfully"})
}

func (h *RetrospectiveHandler) MergeItems(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	retrospectiveIDStr := c.Param("id")
	retrospectiveID, err := uuid.Parse(retrospectiveIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid retrospective ID"})
		return
	}

	var req models.MergeItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sourceItemID, err := uuid.Parse(req.SourceItemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid source item ID"})
		return
	}

	targetItemID, err := uuid.Parse(req.TargetItemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target item ID"})
		return
	}

	mergedItem, err := h.retrospectiveService.MergeItems(sourceItemID, targetItemID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send real-time update via SSE
	if h.realtimeService != nil {
		h.realtimeService.BroadcastToRetrospective(retrospectiveID, "items_merged", map[string]interface{}{
			"merged_item":    mergedItem,
			"source_item_id": sourceItemID,
			"target_item_id": targetItemID,
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Items merged successfully", "merged_item": mergedItem})
}

func (h *RetrospectiveHandler) SetupRoutes(r *gin.RouterGroup) {
	retrospectives := r.Group("/retrospectives")
	retrospectives.Use(authMiddleware)
	{
		retrospectives.POST("", h.CreateRetrospective)
		retrospectives.GET("", h.GetUserRetrospectives)
		// Action Items routes (must be before /:id routes to avoid conflicts)
		retrospectives.PUT("/action-items/:actionItemId", h.UpdateActionItem)
		retrospectives.DELETE("/action-items/:actionItemId", h.DeleteActionItem)
		retrospectives.DELETE("/items/:itemId", h.DeleteItem)
		retrospectives.POST("/items/:itemId/vote", h.VoteItem)
		retrospectives.POST("/groups/:groupId/vote", h.VoteGroup)
		retrospectives.DELETE("/groups/:groupId", h.DeleteGroup)
		// Retrospective-specific routes
		retrospectives.GET("/:id", h.GetRetrospective)
		retrospectives.PUT("/:id", h.UpdateRetrospective)
		retrospectives.DELETE("/:id", h.DeleteRetrospective)
		retrospectives.POST("/:id/start", h.StartRetrospective)
		retrospectives.POST("/:id/end", h.EndRetrospective)
		retrospectives.POST("/:id/reopen", h.ReopenRetrospective)
		retrospectives.POST("/:id/items", h.AddItem)
		retrospectives.POST("/:id/action-items", h.AddActionItem)
		retrospectives.POST("/:id/join", h.JoinRetrospective)
		retrospectives.GET("/:id/participants", h.GetParticipants)
		retrospectives.POST("/:id/groups", h.CreateGroup)
		retrospectives.POST("/:id/merge-items", h.MergeItems)
		retrospectives.PUT("/:id/blur", h.ToggleBlur)
		retrospectives.GET("/:id/export", h.ExportRetrospective)
	}
}

// ToggleBlur handles blur toggle for a retrospective
func (h *RetrospectiveHandler) ToggleBlur(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	retrospectiveIDStr := c.Param("id")
	retrospectiveID, err := uuid.Parse(retrospectiveIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid retrospective ID"})
		return
	}

	var req struct {
		Blurred bool `json:"blurred"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify user is the owner of the retrospective
	retrospective, err := h.retrospectiveService.GetRetrospectiveWithDetails(retrospectiveID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "retrospective not found"})
		return
	}

	if retrospective.CreatedBy != userID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only retrospective owner can toggle blur"})
		return
	}

	// Store blur state in realtime service
	h.realtimeService.SetBlurState(retrospectiveID, req.Blurred)

	// Broadcast blur state to all participants
	h.realtimeService.BroadcastToRetrospective(retrospectiveID, "blur_toggled", map[string]interface{}{
		"blurred": req.Blurred,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Blur state updated successfully", "blurred": req.Blurred})
}

// ExportRetrospective godoc
// @Summary Export retrospective to PDF
// @Description Export a retrospective to PDF format (only accessible by the creator)
// @Tags Retrospectives
// @Accept json
// @Produce application/pdf
// @Security BearerAuth
// @Param id path string true "Retrospective ID"
// @Success 200 {file} binary "PDF file"
// @Failure 400 {object} map[string]string "Invalid retrospective ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied - only creator can export"
// @Failure 404 {object} map[string]string "Retrospective not found"
// @Router /retrospectives/{id}/export [get]
func (h *RetrospectiveHandler) ExportRetrospective(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	retrospectiveIDStr := c.Param("id")
	retrospectiveID, err := uuid.Parse(retrospectiveIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid retrospective ID"})
		return
	}

	// Get retrospective with full details
	retrospective, err := h.retrospectiveService.GetRetrospectiveWithDetails(retrospectiveID, userID.(uuid.UUID))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "access denied" {
			status = http.StatusForbidden
		} else if err.Error() == "retrospective not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// Check if user is the creator
	if retrospective.CreatedBy != userID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only the retrospective creator can export"})
		return
	}

	// Generate PDF content
	pdfContent, err := h.generateRetrospectivePDF(retrospective)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate PDF"})
		return
	}

	// Set headers for PDF download
	filename := fmt.Sprintf("retrospective_%s_%s.pdf", retrospective.Title, time.Now().Format("2006-01-02"))
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Length", fmt.Sprintf("%d", len(pdfContent)))

	c.Data(http.StatusOK, "application/pdf", pdfContent)
}

// generateRetrospectivePDF creates a PDF representation of the retrospective
func (h *RetrospectiveHandler) generateRetrospectivePDF(retrospective *models.RetrospectiveWithDetails) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Title
	pdf.Cell(0, 10, fmt.Sprintf("RETROSPECTIVA: %s", retrospective.Title))
	pdf.Ln(15)

	// Basic info
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Status: %s", retrospective.Status))
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("Template: %s", retrospective.Template))
	pdf.Ln(8)

	if retrospective.EndedAt != nil {
		pdf.Cell(0, 8, fmt.Sprintf("Data de Termino: %s", retrospective.EndedAt.Format("02/01/2006 15:04")))
		pdf.Ln(8)
	}

	pdf.Ln(5)

	// Groups
	if len(retrospective.Groups) > 0 {
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(0, 10, "GRUPOS")
		pdf.Ln(10)

		pdf.SetFont("Arial", "", 10)
		for _, group := range retrospective.Groups {
			pdf.Cell(0, 6, fmt.Sprintf("- %s (%d votos)", group.Name, group.Votes))
			pdf.Ln(6)
		}
		pdf.Ln(5)
	}

	// Items
	if len(retrospective.Items) > 0 {
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(0, 10, "ITENS")
		pdf.Ln(10)

		pdf.SetFont("Arial", "", 10)
		for _, item := range retrospective.Items {
			// Format like in the attachment: - [category] content (votes votos)
			content := fmt.Sprintf("- [%s] %s (%d votos)", item.Category, item.Content, item.Votes)
			pdf.Cell(0, 6, content)
			pdf.Ln(6)
		}
		pdf.Ln(5)
	}

	// Action Items
	if len(retrospective.ActionItems) > 0 {
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(0, 10, "ITENS DE ACAO")
		pdf.Ln(10)

		pdf.SetFont("Arial", "", 10)
		for _, actionItem := range retrospective.ActionItems {
			status := "Pendente"
			if actionItem.Status == "done" {
				status = "Concluido"
			} else if actionItem.Status == "in_progress" {
				status = "Em Progresso"
			}

			// Format like in the attachment: - Title - Status
			pdf.Cell(0, 6, fmt.Sprintf("- %s - %s", actionItem.Title, status))
			pdf.Ln(6)
		}
		pdf.Ln(5)
	}

	// Footer
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(0, 6, fmt.Sprintf("Relatorio gerado em: %s", time.Now().Format("02/01/2006 15:04")))

	// Generate PDF bytes
	var buf strings.Builder
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return []byte(buf.String()), nil
}

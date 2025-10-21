package handlers

import (
	"net/http"

	"educ-retro/internal/models"
	"educ-retro/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TeamHandler struct {
	teamService *services.TeamService
}

func NewTeamHandler(teamService *services.TeamService) *TeamHandler {
	return &TeamHandler{teamService: teamService}
}

// CreateTeam godoc
// @Summary Create a new team
// @Description Create a new team with the current user as owner
// @Tags Teams
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param team body models.TeamCreateRequest true "Team creation data"
// @Success 201 {object} models.Team "Team created successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /teams [post]
func (h *TeamHandler) CreateTeam(c *gin.Context) {
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

// GetUserTeams godoc
// @Summary Get user's teams
// @Description Get all teams where the current user is a member
// @Tags Teams
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Team "User's teams"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /teams [get]
func (h *TeamHandler) GetUserTeams(c *gin.Context) {
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

// GetTeam godoc
// @Summary Get team details
// @Description Get detailed information about a specific team including members
// @Tags Teams
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Team ID"
// @Success 200 {object} models.TeamWithMembers "Team details"
// @Failure 400 {object} map[string]string "Invalid team ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied"
// @Failure 404 {object} map[string]string "Team not found"
// @Router /teams/{id} [get]
func (h *TeamHandler) GetTeam(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	teamIDStr := c.Param("id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team ID"})
		return
	}

	team, err := h.teamService.GetTeam(teamID, userID.(uuid.UUID))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "access denied" {
			status = http.StatusForbidden
		} else if err.Error() == "team not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, team)
}

// UpdateTeam godoc
// @Summary Update team
// @Description Update team information (only team owner can update)
// @Tags Teams
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Team ID"
// @Param team body models.TeamCreateRequest true "Team update data"
// @Success 200 {object} models.Team "Updated team"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied"
// @Router /teams/{id} [put]
func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	teamIDStr := c.Param("id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team ID"})
		return
	}

	var req models.TeamCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team, err := h.teamService.UpdateTeam(teamID, userID.(uuid.UUID), &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "only team owner can update team" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, team)
}

// DeleteTeam godoc
// @Summary Delete team
// @Description Delete a team (only team owner can delete)
// @Tags Teams
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Team ID"
// @Success 204 "Team deleted successfully"
// @Failure 400 {object} map[string]string "Invalid team ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied"
// @Router /teams/{id} [delete]
func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	teamIDStr := c.Param("id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team ID"})
		return
	}

	err = h.teamService.DeleteTeam(teamID, userID.(uuid.UUID))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "only team owner can delete team" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// AddMember godoc
// @Summary Add team member
// @Description Add a new member to the team
// @Tags Teams
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Team ID"
// @Param invite body models.TeamInviteRequest true "Member invitation data"
// @Success 200 {object} map[string]string "Member added successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied"
// @Router /teams/{id}/members [post]
func (h *TeamHandler) AddMember(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	teamIDStr := c.Param("id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team ID"})
		return
	}

	var req models.TeamInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.teamService.AddMember(teamID, userID.(uuid.UUID), &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "insufficient permissions" || err.Error() == "user not found" || err.Error() == "user is already a member of this team" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member added successfully"})
}

// RemoveMember godoc
// @Summary Remove team member
// @Description Remove a member from the team
// @Tags Teams
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Team ID"
// @Param userId path string true "User ID"
// @Success 200 {object} map[string]string "Member removed successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied"
// @Router /teams/{id}/members/{userId} [delete]
func (h *TeamHandler) RemoveMember(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	teamIDStr := c.Param("id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team ID"})
		return
	}

	targetUserIDStr := c.Param("userId")
	targetUserID, err := uuid.Parse(targetUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	err = h.teamService.RemoveMember(teamID, userID.(uuid.UUID), targetUserID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "insufficient permissions" || err.Error() == "cannot remove the only owner of the team" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed successfully"})
}

func (h *TeamHandler) SetupRoutes(r *gin.RouterGroup) {
	teams := r.Group("/teams")
	teams.Use(authMiddleware)
	{
		teams.POST("", h.CreateTeam)
		teams.GET("", h.GetUserTeams)
		teams.GET("/:id", h.GetTeam)
		teams.PUT("/:id", h.UpdateTeam)
		teams.DELETE("/:id", h.DeleteTeam)
		teams.POST("/:id/members", h.AddMember)
		teams.DELETE("/:id/members/:userId", h.RemoveMember)
	}
}

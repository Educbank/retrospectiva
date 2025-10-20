package services

import (
	"errors"

	"educ-retro/internal/models"
	"educ-retro/internal/repositories"

	"github.com/google/uuid"
)

type TeamService struct {
	teamRepo *repositories.TeamRepository
	userRepo *repositories.UserRepository
}

func NewTeamService(teamRepo *repositories.TeamRepository, userRepo *repositories.UserRepository) *TeamService {
	return &TeamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

func (s *TeamService) CreateTeam(userID uuid.UUID, req *models.TeamCreateRequest) (*models.Team, error) {
	team := &models.Team{
		Name:        req.Name,
		Description: &req.Description,
		OwnerID:     userID,
	}

	err := s.teamRepo.Create(team)
	if err != nil {
		return nil, err
	}

	return team, nil
}

func (s *TeamService) GetUserTeams(userID uuid.UUID) ([]models.TeamWithCounts, error) {
	teams, err := s.teamRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	return teams, nil
}

func (s *TeamService) GetTeam(teamID uuid.UUID, userID uuid.UUID) (*models.TeamWithMembers, error) {
	// Check if user is member of the team
	isMember, err := s.teamRepo.IsMember(teamID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("access denied")
	}

	// Get team
	team, err := s.teamRepo.GetByID(teamID)
	if err != nil {
		return nil, err
	}

	// Get members
	members, err := s.teamRepo.GetMembers(teamID)
	if err != nil {
		return nil, err
	}

	// Get owner
	owner, err := s.userRepo.GetByID(team.OwnerID)
	if err != nil {
		return nil, err
	}

	ownerResponse := models.UserResponse{
		ID:        owner.ID,
		Email:     owner.Email,
		Name:      owner.Name,
		Avatar:    owner.Avatar,
		CreatedAt: owner.CreatedAt,
	}

	return &models.TeamWithMembers{
		Team:    *team,
		Members: members,
		Owner:   ownerResponse,
	}, nil
}

func (s *TeamService) UpdateTeam(teamID, userID uuid.UUID, req *models.TeamCreateRequest) (*models.Team, error) {
	// Check if user is owner
	role, err := s.teamRepo.GetMemberRole(teamID, userID)
	if err != nil {
		return nil, err
	}
	if role != "owner" {
		return nil, errors.New("only team owner can update team")
	}

	// Get team
	team, err := s.teamRepo.GetByID(teamID)
	if err != nil {
		return nil, err
	}

	// Update team
	team.Name = req.Name
	team.Description = &req.Description

	err = s.teamRepo.Update(team)
	if err != nil {
		return nil, err
	}

	return team, nil
}

func (s *TeamService) DeleteTeam(teamID, userID uuid.UUID) error {
	// Check if user is owner
	role, err := s.teamRepo.GetMemberRole(teamID, userID)
	if err != nil {
		return err
	}
	if role != "owner" {
		return errors.New("only team owner can delete team")
	}

	return s.teamRepo.Delete(teamID)
}

func (s *TeamService) AddMember(teamID, userID uuid.UUID, req *models.TeamInviteRequest) error {
	// Check if current user is owner or member
	role, err := s.teamRepo.GetMemberRole(teamID, userID)
	if err != nil {
		return err
	}
	if role != "owner" && role != "member" {
		return errors.New("insufficient permissions")
	}

	// Check if user exists
	invitedUser, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return errors.New("user not found")
	}

	// Check if user is already a member
	isMember, err := s.teamRepo.IsMember(teamID, invitedUser.ID)
	if err != nil {
		return err
	}
	if isMember {
		return errors.New("user is already a member of this team")
	}

	return s.teamRepo.AddMember(teamID, invitedUser.ID, req.Role)
}

func (s *TeamService) RemoveMember(teamID, userID, targetUserID uuid.UUID) error {
	// Check if current user is owner or removing themselves
	role, err := s.teamRepo.GetMemberRole(teamID, userID)
	if err != nil {
		return err
	}
	if role != "owner" && userID != targetUserID {
		return errors.New("insufficient permissions")
	}

	// Don't allow owner to remove themselves if they're the only owner
	if role == "owner" && userID == targetUserID {
		members, err := s.teamRepo.GetMembers(teamID)
		if err != nil {
			return err
		}

		ownerCount := 0
		for _, member := range members {
			if member.Role == "owner" {
				ownerCount++
			}
		}

		if ownerCount == 1 {
			return errors.New("cannot remove the only owner of the team")
		}
	}

	return s.teamRepo.RemoveMember(teamID, targetUserID)
}

func (s *TeamService) UpdateMemberRole(teamID, userID, targetUserID uuid.UUID, newRole string) error {
	// Check if current user is owner
	role, err := s.teamRepo.GetMemberRole(teamID, userID)
	if err != nil {
		return err
	}
	if role != "owner" {
		return errors.New("only team owner can update member roles")
	}

	// Don't allow changing owner role
	if newRole == "owner" {
		return errors.New("cannot change role to owner")
	}

	return s.teamRepo.UpdateMemberRole(teamID, targetUserID, newRole)
}

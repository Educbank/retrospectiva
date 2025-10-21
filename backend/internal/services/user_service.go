package services

import (
	"errors"

	"educ-retro/internal/auth"
	"educ-retro/internal/models"
	"educ-retro/internal/repositories"
	"educ-retro/internal/utils"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo repositories.UserRepositoryInterface
}

func NewUserService(userRepo repositories.UserRepositoryInterface) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) Register(req *models.UserCreateRequest) (*models.UserResponse, string, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, "", errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, "", err
	}

	// Create user
	user := &models.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: hashedPassword,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, "", err
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Email, user.Name)
	if err != nil {
		return nil, "", err
	}

	// Return user response (without password)
	userResponse := &models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
	}

	return userResponse, token, nil
}

func (s *UserService) Login(req *models.UserLoginRequest) (*models.UserResponse, string, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Email, user.Name)
	if err != nil {
		return nil, "", err
	}

	// Return user response (without password)
	userResponse := &models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
	}

	return userResponse, token, nil
}

func (s *UserService) GetProfile(userID uuid.UUID) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	userResponse := &models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
	}

	return userResponse, nil
}

func (s *UserService) UpdateProfile(userID uuid.UUID, name, avatar string) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.Name = name
	if avatar != "" {
		user.Avatar = &avatar
	}

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	userResponse := &models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
	}

	return userResponse, nil
}

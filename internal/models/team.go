package models

import (
	"time"

	"github.com/google/uuid"
)

type Team struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	OwnerID     uuid.UUID `json:"owner_id" db:"owner_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type TeamMember struct {
	ID       uuid.UUID `json:"id" db:"id"`
	TeamID   uuid.UUID `json:"team_id" db:"team_id"`
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	Role     string    `json:"role" db:"role"` // owner, member, viewer
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
}

type TeamMemberWithUser struct {
	ID        uuid.UUID `json:"id" db:"id"`
	TeamID    uuid.UUID `json:"team_id" db:"team_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	UserName  string    `json:"user_name" db:"user_name"`
	UserEmail string    `json:"user_email" db:"user_email"`
	Role      string    `json:"role" db:"role"` // owner, member, viewer
	JoinedAt  time.Time `json:"joined_at" db:"joined_at"`
}

type TeamWithMembers struct {
	Team    Team                 `json:"team"`
	Members []TeamMemberWithUser `json:"members"`
	Owner   UserResponse         `json:"owner"`
}

type TeamWithCounts struct {
	Team
	MemberCount        int `json:"member_count" db:"member_count"`
	RetrospectiveCount int `json:"retrospective_count" db:"retrospective_count"`
}

type TeamCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type TeamInviteRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required,oneof=member viewer"`
}

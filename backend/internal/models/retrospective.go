package models

import (
	"time"

	"github.com/google/uuid"
)

type RetrospectiveStatus string

const (
	RetroStatusPlanned    RetrospectiveStatus = "planned"
	RetroStatusActive     RetrospectiveStatus = "active"
	RetroStatusCollecting RetrospectiveStatus = "collecting"
	RetroStatusVoting     RetrospectiveStatus = "voting"
	RetroStatusDiscussing RetrospectiveStatus = "discussing"
	RetroStatusClosed     RetrospectiveStatus = "closed"
)

type RetrospectiveTemplate string

const (
	TemplateStartStopContinue RetrospectiveTemplate = "start_stop_continue"
	Template4Ls               RetrospectiveTemplate = "4ls" // Liked, Learned, Lacked, Longed for
	TemplateMadSadGlad        RetrospectiveTemplate = "mad_sad_glad"
	TemplateSailboat          RetrospectiveTemplate = "sailboat"
	TemplateWentWellToImprove RetrospectiveTemplate = "went_well_to_improve"
)

type Retrospective struct {
	ID          uuid.UUID             `json:"id" db:"id"`
	TeamID      uuid.UUID             `json:"team_id,omitempty" db:"team_id" swaggerignore:"true"`
	Title       string                `json:"title" db:"title"`
	Description *string               `json:"description" db:"description"`
	Template    RetrospectiveTemplate `json:"template" db:"template"`
	Status      RetrospectiveStatus   `json:"status" db:"status"`
	ScheduledAt *time.Time            `json:"scheduled_at" db:"scheduled_at"`
	StartedAt   *time.Time            `json:"started_at" db:"started_at"`
	EndedAt     *time.Time            `json:"ended_at" db:"ended_at"`
	CreatedBy   uuid.UUID             `json:"created_by" db:"created_by"`
	CreatedAt   time.Time             `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at" db:"updated_at"`
}

type RetrospectiveItem struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	RetrospectiveID uuid.UUID  `json:"retrospective_id" db:"retrospective_id"`
	Category        string     `json:"category" db:"category"` // depends on template
	Content         string     `json:"content" db:"content"`
	AuthorID        *uuid.UUID `json:"author_id" db:"author_id"` // null if anonymous
	IsAnonymous     bool       `json:"is_anonymous" db:"is_anonymous"`
	Votes           int        `json:"votes" db:"votes"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

type RetrospectiveVote struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ItemID    uuid.UUID `json:"item_id" db:"item_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ActionItem struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	RetrospectiveID uuid.UUID  `json:"retrospective_id" db:"retrospective_id"`
	ItemID          *uuid.UUID `json:"item_id" db:"item_id"` // can be null for general action items
	Title           string     `json:"title" db:"title"`
	Description     *string    `json:"description" db:"description"`
	AssignedTo      *uuid.UUID `json:"assigned_to" db:"assigned_to"`
	Status          string     `json:"status" db:"status"` // todo, in_progress, done
	DueDate         *time.Time `json:"due_date" db:"due_date"`
	CompletedAt     *time.Time `json:"completed_at" db:"completed_at"`
	CreatedBy       uuid.UUID  `json:"created_by" db:"created_by"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

type RetrospectiveCreateRequest struct {
	Title       string                `json:"title" binding:"required"`
	Description string                `json:"description"`
	Template    RetrospectiveTemplate `json:"template" binding:"required"`
}

type RetrospectiveItemCreateRequest struct {
	Category    string `json:"category" binding:"required"`
	Content     string `json:"content" binding:"required"`
	IsAnonymous bool   `json:"is_anonymous"`
}

type ActionItemCreateRequest struct {
	ItemID      *string `json:"item_id"`
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	AssignedTo  *string `json:"assigned_to"`
	DueDate     *string `json:"due_date"`
}

type ActionItemUpdateRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
	AssignedTo  *string `json:"assigned_to"`
	DueDate     *string `json:"due_date"`
	CompletedAt *string `json:"completed_at"`
}

type RetrospectiveParticipant struct {
	ID              uuid.UUID `json:"id" db:"id"`
	RetrospectiveID uuid.UUID `json:"retrospective_id" db:"retrospective_id"`
	UserID          uuid.UUID `json:"user_id" db:"user_id"`
	JoinedAt        time.Time `json:"joined_at" db:"joined_at"`
	LastSeen        time.Time `json:"last_seen" db:"last_seen"`
}

type RetrospectiveGroup struct {
	ID              uuid.UUID `json:"id" db:"id"`
	RetrospectiveID uuid.UUID `json:"retrospective_id" db:"retrospective_id"`
	Name            string    `json:"name" db:"name"`
	Description     *string   `json:"description" db:"description"`
	Votes           int       `json:"votes" db:"votes"`
	CreatedBy       uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type RetrospectiveGroupItem struct {
	ID      uuid.UUID `json:"id" db:"id"`
	GroupID uuid.UUID `json:"group_id" db:"group_id"`
	ItemID  uuid.UUID `json:"item_id" db:"item_id"`
}

type GroupCreateRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description *string  `json:"description"`
	ItemIDs     []string `json:"item_ids"`
}

type GroupVoteRequest struct {
	GroupID string `json:"group_id" binding:"required"`
}

type MergeItemsRequest struct {
	SourceItemID string `json:"source_item_id" binding:"required"`
	TargetItemID string `json:"target_item_id" binding:"required"`
}

type RetrospectiveWithDetails struct {
	Retrospective
	Items        []RetrospectiveItem        `json:"items"`
	ActionItems  []ActionItem               `json:"action_items"`
	Participants []RetrospectiveParticipant `json:"participants"`
	Groups       []RetrospectiveGroup       `json:"groups"`
}

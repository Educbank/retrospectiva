package schema

import (
	"encoding/json"
	"testing"
	"time"

	"educ-retro/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUserSchema validates the User model structure
func TestUserSchema(t *testing.T) {
	tests := []struct {
		name     string
		user     models.User
		validate func(t *testing.T, original, unmarshaled models.User)
	}{
		{
			name: "Complete User with all fields",
			user: models.User{
				ID:        uuid.New(),
				Email:     "test@example.com",
				Name:      "Test User",
				Password:  "hashedpassword",
				Avatar:    stringPtr("avatar.jpg"),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			validate: func(t *testing.T, original, unmarshaled models.User) {
				assert.Equal(t, original.ID, unmarshaled.ID)
				assert.Equal(t, original.Email, unmarshaled.Email)
				assert.Equal(t, original.Name, unmarshaled.Name)
				assert.Equal(t, original.Avatar, unmarshaled.Avatar)
				assert.Equal(t, original.CreatedAt.Unix(), unmarshaled.CreatedAt.Unix())
				assert.Equal(t, original.UpdatedAt.Unix(), unmarshaled.UpdatedAt.Unix())
				// Password should be hidden in JSON
				assert.Empty(t, unmarshaled.Password)
			},
		},
		{
			name: "User with nil avatar",
			user: models.User{
				ID:        uuid.New(),
				Email:     "test2@example.com",
				Name:      "Test User 2",
				Password:  "hashedpassword2",
				Avatar:    nil,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			validate: func(t *testing.T, original, unmarshaled models.User) {
				assert.Equal(t, original.ID, unmarshaled.ID)
				assert.Equal(t, original.Email, unmarshaled.Email)
				assert.Equal(t, original.Name, unmarshaled.Name)
				assert.Nil(t, unmarshaled.Avatar)
				assert.Empty(t, unmarshaled.Password)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.user)
			require.NoError(t, err)

			// Test JSON unmarshaling
			var unmarshaledUser models.User
			err = json.Unmarshal(jsonData, &unmarshaledUser)
			require.NoError(t, err)

			// Validate fields
			tt.validate(t, tt.user, unmarshaledUser)
		})
	}
}

// TestUserCreateRequestSchema validates the UserCreateRequest structure
func TestUserCreateRequestSchema(t *testing.T) {
	tests := []struct {
		name     string
		req      models.UserCreateRequest
		validate func(t *testing.T, original, unmarshaled models.UserCreateRequest)
	}{
		{
			name: "Complete UserCreateRequest",
			req: models.UserCreateRequest{
				Email:    "test@example.com",
				Name:     "Test User",
				Password: "password123",
			},
			validate: func(t *testing.T, original, unmarshaled models.UserCreateRequest) {
				assert.Equal(t, original.Email, unmarshaled.Email)
				assert.Equal(t, original.Name, unmarshaled.Name)
				assert.Equal(t, original.Password, unmarshaled.Password)
			},
		},
		{
			name: "UserCreateRequest with special characters",
			req: models.UserCreateRequest{
				Email:    "user+test@example.com",
				Name:     "José da Silva",
				Password: "p@ssw0rd!@#",
			},
			validate: func(t *testing.T, original, unmarshaled models.UserCreateRequest) {
				assert.Equal(t, original.Email, unmarshaled.Email)
				assert.Equal(t, original.Name, unmarshaled.Name)
				assert.Equal(t, original.Password, unmarshaled.Password)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.req)
			require.NoError(t, err)

			// Test JSON unmarshaling
			var unmarshaledReq models.UserCreateRequest
			err = json.Unmarshal(jsonData, &unmarshaledReq)
			require.NoError(t, err)

			// Validate fields
			tt.validate(t, tt.req, unmarshaledReq)
		})
	}
}

// TestUserLoginRequestSchema validates the UserLoginRequest structure
func TestUserLoginRequestSchema(t *testing.T) {
	tests := []struct {
		name     string
		req      models.UserLoginRequest
		validate func(t *testing.T, original, unmarshaled models.UserLoginRequest)
	}{
		{
			name: "Standard login request",
			req: models.UserLoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			validate: func(t *testing.T, original, unmarshaled models.UserLoginRequest) {
				assert.Equal(t, original.Email, unmarshaled.Email)
				assert.Equal(t, original.Password, unmarshaled.Password)
			},
		},
		{
			name: "Login with special characters",
			req: models.UserLoginRequest{
				Email:    "user+test@example.com",
				Password: "p@ssw0rd!@#",
			},
			validate: func(t *testing.T, original, unmarshaled models.UserLoginRequest) {
				assert.Equal(t, original.Email, unmarshaled.Email)
				assert.Equal(t, original.Password, unmarshaled.Password)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.req)
			require.NoError(t, err)

			// Test JSON unmarshaling
			var unmarshaledReq models.UserLoginRequest
			err = json.Unmarshal(jsonData, &unmarshaledReq)
			require.NoError(t, err)

			// Validate fields
			tt.validate(t, tt.req, unmarshaledReq)
		})
	}
}

// TestUserResponseSchema validates the UserResponse structure
func TestUserResponseSchema(t *testing.T) {
	response := models.UserResponse{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Name:      "Test User",
		Avatar:    stringPtr("avatar.jpg"),
		CreatedAt: time.Now(),
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(response)
	require.NoError(t, err)

	// Test JSON unmarshaling
	var unmarshaledResponse models.UserResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	require.NoError(t, err)

	// Validate fields
	assert.Equal(t, response.ID, unmarshaledResponse.ID)
	assert.Equal(t, response.Email, unmarshaledResponse.Email)
	assert.Equal(t, response.Name, unmarshaledResponse.Name)
	assert.Equal(t, response.Avatar, unmarshaledResponse.Avatar)
	assert.Equal(t, response.CreatedAt.Unix(), unmarshaledResponse.CreatedAt.Unix())
}

// TestRetrospectiveSchema validates the Retrospective model structure
func TestRetrospectiveSchema(t *testing.T) {
	now := time.Now()
	retro := models.Retrospective{
		ID:          uuid.New(),
		TeamID:      uuid.New(),
		Title:       "Test Retrospective",
		Description: stringPtr("A test retrospective"),
		Template:    models.TemplateStartStopContinue,
		Status:      models.RetroStatusPlanned,
		ScheduledAt: &now,
		StartedAt:   &now,
		EndedAt:     &now,
		CreatedBy:   uuid.New(),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(retro)
	require.NoError(t, err)

	// Test JSON unmarshaling
	var unmarshaledRetro models.Retrospective
	err = json.Unmarshal(jsonData, &unmarshaledRetro)
	require.NoError(t, err)

	// Validate fields
	assert.Equal(t, retro.ID, unmarshaledRetro.ID)
	assert.Equal(t, retro.TeamID, unmarshaledRetro.TeamID)
	assert.Equal(t, retro.Title, unmarshaledRetro.Title)
	assert.Equal(t, retro.Description, unmarshaledRetro.Description)
	assert.Equal(t, retro.Template, unmarshaledRetro.Template)
	assert.Equal(t, retro.Status, unmarshaledRetro.Status)
	assert.Equal(t, retro.CreatedBy, unmarshaledRetro.CreatedBy)
	assert.Equal(t, retro.CreatedAt.Unix(), unmarshaledRetro.CreatedAt.Unix())
	assert.Equal(t, retro.UpdatedAt.Unix(), unmarshaledRetro.UpdatedAt.Unix())
}

// TestRetrospectiveCreateRequestSchema validates the RetrospectiveCreateRequest structure
func TestRetrospectiveCreateRequestSchema(t *testing.T) {
	req := models.RetrospectiveCreateRequest{
		Title:       "Test Retrospective",
		Description: "A test retrospective",
		Template:    models.TemplateStartStopContinue,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(req)
	require.NoError(t, err)

	// Test JSON unmarshaling
	var unmarshaledReq models.RetrospectiveCreateRequest
	err = json.Unmarshal(jsonData, &unmarshaledReq)
	require.NoError(t, err)

	// Validate fields
	assert.Equal(t, req.Title, unmarshaledReq.Title)
	assert.Equal(t, req.Description, unmarshaledReq.Description)
	assert.Equal(t, req.Template, unmarshaledReq.Template)
}

// TestRetrospectiveStatusEnum validates all retrospective status values
func TestRetrospectiveStatusEnum(t *testing.T) {
	tests := []struct {
		name   string
		status models.RetrospectiveStatus
	}{
		{"Planned status", models.RetroStatusPlanned},
		{"Active status", models.RetroStatusActive},
		{"Collecting status", models.RetroStatusCollecting},
		{"Voting status", models.RetroStatusVoting},
		{"Discussing status", models.RetroStatusDiscussing},
		{"Closed status", models.RetroStatusClosed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.status)
			require.NoError(t, err)

			// Test JSON unmarshaling
			var unmarshaledStatus models.RetrospectiveStatus
			err = json.Unmarshal(jsonData, &unmarshaledStatus)
			require.NoError(t, err)

			assert.Equal(t, tt.status, unmarshaledStatus)
		})
	}
}

// TestRetrospectiveTemplateEnum validates all retrospective template values
func TestRetrospectiveTemplateEnum(t *testing.T) {
	tests := []struct {
		name     string
		template models.RetrospectiveTemplate
	}{
		{"Start Stop Continue template", models.TemplateStartStopContinue},
		{"4Ls template", models.Template4Ls},
		{"Mad Sad Glad template", models.TemplateMadSadGlad},
		{"Sailboat template", models.TemplateSailboat},
		{"Went Well To Improve template", models.TemplateWentWellToImprove},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.template)
			require.NoError(t, err)

			// Test JSON unmarshaling
			var unmarshaledTemplate models.RetrospectiveTemplate
			err = json.Unmarshal(jsonData, &unmarshaledTemplate)
			require.NoError(t, err)

			assert.Equal(t, tt.template, unmarshaledTemplate)
		})
	}
}

// TestRetrospectiveItemSchema validates the RetrospectiveItem model structure
func TestRetrospectiveItemSchema(t *testing.T) {
	item := models.RetrospectiveItem{
		ID:              uuid.New(),
		RetrospectiveID: uuid.New(),
		Category:        "Start",
		Content:         "Test content",
		AuthorID:        uuidPtr(uuid.New()),
		IsAnonymous:     true,
		Votes:           5,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(item)
	require.NoError(t, err)

	// Test JSON unmarshaling
	var unmarshaledItem models.RetrospectiveItem
	err = json.Unmarshal(jsonData, &unmarshaledItem)
	require.NoError(t, err)

	// Validate fields
	assert.Equal(t, item.ID, unmarshaledItem.ID)
	assert.Equal(t, item.RetrospectiveID, unmarshaledItem.RetrospectiveID)
	assert.Equal(t, item.Category, unmarshaledItem.Category)
	assert.Equal(t, item.Content, unmarshaledItem.Content)
	assert.Equal(t, item.AuthorID, unmarshaledItem.AuthorID)
	assert.Equal(t, item.IsAnonymous, unmarshaledItem.IsAnonymous)
	assert.Equal(t, item.Votes, unmarshaledItem.Votes)
	assert.Equal(t, item.CreatedAt.Unix(), unmarshaledItem.CreatedAt.Unix())
	assert.Equal(t, item.UpdatedAt.Unix(), unmarshaledItem.UpdatedAt.Unix())
}

// TestActionItemSchema validates the ActionItem model structure
func TestActionItemSchema(t *testing.T) {
	now := time.Now()
	actionItem := models.ActionItem{
		ID:              uuid.New(),
		RetrospectiveID: uuid.New(),
		ItemID:          uuidPtr(uuid.New()),
		Title:           "Test Action Item",
		Description:     stringPtr("A test action item"),
		AssignedTo:      uuidPtr(uuid.New()),
		Status:          "todo",
		DueDate:         &now,
		CompletedAt:     &now,
		CreatedBy:       uuid.New(),
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(actionItem)
	require.NoError(t, err)

	// Test JSON unmarshaling
	var unmarshaledActionItem models.ActionItem
	err = json.Unmarshal(jsonData, &unmarshaledActionItem)
	require.NoError(t, err)

	// Validate fields
	assert.Equal(t, actionItem.ID, unmarshaledActionItem.ID)
	assert.Equal(t, actionItem.RetrospectiveID, unmarshaledActionItem.RetrospectiveID)
	assert.Equal(t, actionItem.ItemID, unmarshaledActionItem.ItemID)
	assert.Equal(t, actionItem.Title, unmarshaledActionItem.Title)
	assert.Equal(t, actionItem.Description, unmarshaledActionItem.Description)
	assert.Equal(t, actionItem.AssignedTo, unmarshaledActionItem.AssignedTo)
	assert.Equal(t, actionItem.Status, unmarshaledActionItem.Status)
	assert.Equal(t, actionItem.CreatedBy, unmarshaledActionItem.CreatedBy)
	assert.Equal(t, actionItem.CreatedAt.Unix(), unmarshaledActionItem.CreatedAt.Unix())
	assert.Equal(t, actionItem.UpdatedAt.Unix(), unmarshaledActionItem.UpdatedAt.Unix())
}

// TestRetrospectiveParticipantSchema validates the RetrospectiveParticipant model structure
func TestRetrospectiveParticipantSchema(t *testing.T) {
	participant := models.RetrospectiveParticipant{
		ID:              uuid.New(),
		RetrospectiveID: uuid.New(),
		UserID:          uuid.New(),
		JoinedAt:        time.Now(),
		LastSeen:        time.Now(),
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(participant)
	require.NoError(t, err)

	// Test JSON unmarshaling
	var unmarshaledParticipant models.RetrospectiveParticipant
	err = json.Unmarshal(jsonData, &unmarshaledParticipant)
	require.NoError(t, err)

	// Validate fields
	assert.Equal(t, participant.ID, unmarshaledParticipant.ID)
	assert.Equal(t, participant.RetrospectiveID, unmarshaledParticipant.RetrospectiveID)
	assert.Equal(t, participant.UserID, unmarshaledParticipant.UserID)
	assert.Equal(t, participant.JoinedAt.Unix(), unmarshaledParticipant.JoinedAt.Unix())
	assert.Equal(t, participant.LastSeen.Unix(), unmarshaledParticipant.LastSeen.Unix())
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func uuidPtr(u uuid.UUID) *uuid.UUID {
	return &u
}

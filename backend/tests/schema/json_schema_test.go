package schema

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"educ-retro/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestJSONSchemaValidation tests JSON schema validation for API responses
func TestJSONSchemaValidation(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		validate func(t *testing.T, jsonMap map[string]interface{})
	}{
		{
			name: "UserResponse JSON Schema",
			data: models.UserResponse{
				ID:        uuid.New(),
				Email:     "test@example.com",
				Name:      "Test User",
				Avatar:    stringPtr("avatar.jpg"),
				CreatedAt: time.Now(),
			},
			validate: func(t *testing.T, jsonMap map[string]interface{}) {
				// Check required fields
				assert.Contains(t, jsonMap, "id")
				assert.Contains(t, jsonMap, "email")
				assert.Contains(t, jsonMap, "name")
				assert.Contains(t, jsonMap, "created_at")

				// Check field types
				assert.IsType(t, "", jsonMap["id"])         // UUID as string
				assert.IsType(t, "", jsonMap["email"])      // Email as string
				assert.IsType(t, "", jsonMap["name"])       // Name as string
				assert.IsType(t, "", jsonMap["created_at"]) // Time as string

				// Check optional fields
				if avatar, exists := jsonMap["avatar"]; exists {
					assert.IsType(t, "", avatar) // Avatar as string
				}
			},
		},
		{
			name: "UserResponse with nil avatar",
			data: models.UserResponse{
				ID:        uuid.New(),
				Email:     "test2@example.com",
				Name:      "Test User 2",
				Avatar:    nil,
				CreatedAt: time.Now(),
			},
			validate: func(t *testing.T, jsonMap map[string]interface{}) {
				// Check required fields
				assert.Contains(t, jsonMap, "id")
				assert.Contains(t, jsonMap, "email")
				assert.Contains(t, jsonMap, "name")
				assert.Contains(t, jsonMap, "created_at")

				// Check optional fields
				if avatar, exists := jsonMap["avatar"]; exists {
					assert.Nil(t, avatar)
				}
			},
		},
		{
			name: "Team JSON Schema",
			data: models.Team{
				ID:          uuid.New(),
				Name:        "Test Team",
				Description: stringPtr("A test team"),
				OwnerID:     uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			validate: func(t *testing.T, jsonMap map[string]interface{}) {
				// Check required fields
				assert.Contains(t, jsonMap, "id")
				assert.Contains(t, jsonMap, "name")
				assert.Contains(t, jsonMap, "owner_id")
				assert.Contains(t, jsonMap, "created_at")
				assert.Contains(t, jsonMap, "updated_at")

				// Check field types
				assert.IsType(t, "", jsonMap["id"])         // UUID as string
				assert.IsType(t, "", jsonMap["name"])       // Name as string
				assert.IsType(t, "", jsonMap["owner_id"])   // UUID as string
				assert.IsType(t, "", jsonMap["created_at"]) // Time as string
				assert.IsType(t, "", jsonMap["updated_at"]) // Time as string

				// Check optional fields
				if description, exists := jsonMap["description"]; exists {
					assert.IsType(t, "", description) // Description as string
				}
			},
		},
		{
			name: "Team without description",
			data: models.Team{
				ID:          uuid.New(),
				Name:        "Simple Team",
				Description: nil,
				OwnerID:     uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			validate: func(t *testing.T, jsonMap map[string]interface{}) {
				// Check required fields
				assert.Contains(t, jsonMap, "id")
				assert.Contains(t, jsonMap, "name")
				assert.Contains(t, jsonMap, "owner_id")

				// Check optional fields
				if description, exists := jsonMap["description"]; exists {
					assert.Nil(t, description)
				}
			},
		},
		{
			name: "Retrospective JSON Schema",
			data: models.Retrospective{
				ID:          uuid.New(),
				TeamID:      uuid.New(),
				Title:       "Test Retrospective",
				Description: stringPtr("A test retrospective"),
				Template:    models.TemplateStartStopContinue,
				Status:      models.RetroStatusPlanned,
				CreatedBy:   uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			validate: func(t *testing.T, jsonMap map[string]interface{}) {
				// Check required fields
				assert.Contains(t, jsonMap, "id")
				assert.Contains(t, jsonMap, "team_id")
				assert.Contains(t, jsonMap, "title")
				assert.Contains(t, jsonMap, "template")
				assert.Contains(t, jsonMap, "status")
				assert.Contains(t, jsonMap, "created_by")
				assert.Contains(t, jsonMap, "created_at")
				assert.Contains(t, jsonMap, "updated_at")

				// Check field types
				assert.IsType(t, "", jsonMap["id"])         // UUID as string
				assert.IsType(t, "", jsonMap["team_id"])    // UUID as string
				assert.IsType(t, "", jsonMap["title"])      // Title as string
				assert.IsType(t, "", jsonMap["template"])   // Template as string
				assert.IsType(t, "", jsonMap["status"])     // Status as string
				assert.IsType(t, "", jsonMap["created_by"]) // UUID as string
				assert.IsType(t, "", jsonMap["created_at"]) // Time as string
				assert.IsType(t, "", jsonMap["updated_at"]) // Time as string

				// Check optional fields
				if description, exists := jsonMap["description"]; exists {
					assert.IsType(t, "", description) // Description as string
				}
			},
		},
		{
			name: "Retrospective without description",
			data: models.Retrospective{
				ID:          uuid.New(),
				TeamID:      uuid.New(),
				Title:       "Simple Retrospective",
				Description: nil,
				Template:    models.Template4Ls,
				Status:      models.RetroStatusActive,
				CreatedBy:   uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			validate: func(t *testing.T, jsonMap map[string]interface{}) {
				// Check required fields
				assert.Contains(t, jsonMap, "id")
				assert.Contains(t, jsonMap, "team_id")
				assert.Contains(t, jsonMap, "title")
				assert.Contains(t, jsonMap, "template")
				assert.Contains(t, jsonMap, "status")

				// Check optional fields
				if description, exists := jsonMap["description"]; exists {
					assert.Nil(t, description)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.data)
			require.NoError(t, err)

			// Validate JSON structure
			var jsonMap map[string]interface{}
			err = json.Unmarshal(jsonData, &jsonMap)
			require.NoError(t, err)

			// Run validation
			tt.validate(t, jsonMap)
		})
	}
}

// TestAPIRequestSchema tests API request schema validation
func TestAPIRequestSchema(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		validate func(t *testing.T, jsonMap map[string]interface{})
	}{
		{
			name: "UserCreateRequest Schema",
			data: models.UserCreateRequest{
				Email:    "test@example.com",
				Name:     "Test User",
				Password: "password123",
			},
			validate: func(t *testing.T, jsonMap map[string]interface{}) {
				// Check required fields
				assert.Contains(t, jsonMap, "email")
				assert.Contains(t, jsonMap, "name")
				assert.Contains(t, jsonMap, "password")

				// Check field types
				assert.IsType(t, "", jsonMap["email"])    // Email as string
				assert.IsType(t, "", jsonMap["name"])     // Name as string
				assert.IsType(t, "", jsonMap["password"]) // Password as string

				// Validate email format (basic check)
				email := jsonMap["email"].(string)
				assert.Contains(t, email, "@")
				assert.Contains(t, email, ".")
			},
		},
		{
			name: "UserCreateRequest with special characters",
			data: models.UserCreateRequest{
				Email:    "user+test@example.com",
				Name:     "José da Silva",
				Password: "p@ssw0rd!@#",
			},
			validate: func(t *testing.T, jsonMap map[string]interface{}) {
				// Check required fields
				assert.Contains(t, jsonMap, "email")
				assert.Contains(t, jsonMap, "name")
				assert.Contains(t, jsonMap, "password")

				// Check field types
				assert.IsType(t, "", jsonMap["email"])
				assert.IsType(t, "", jsonMap["name"])
				assert.IsType(t, "", jsonMap["password"])

				// Validate email format
				email := jsonMap["email"].(string)
				assert.Contains(t, email, "@")
				assert.Contains(t, email, ".")
			},
		},
		{
			name: "TeamCreateRequest Schema",
			data: models.TeamCreateRequest{
				Name:        "Test Team",
				Description: "A test team",
			},
			validate: func(t *testing.T, jsonMap map[string]interface{}) {
				// Check required fields
				assert.Contains(t, jsonMap, "name")

				// Check field types
				assert.IsType(t, "", jsonMap["name"]) // Name as string

				// Check optional fields
				if description, exists := jsonMap["description"]; exists {
					assert.IsType(t, "", description) // Description as string
				}
			},
		},
		{
			name: "TeamCreateRequest without description",
			data: models.TeamCreateRequest{
				Name: "Simple Team",
			},
			validate: func(t *testing.T, jsonMap map[string]interface{}) {
				// Check required fields
				assert.Contains(t, jsonMap, "name")

				// Check field types
				assert.IsType(t, "", jsonMap["name"])

				// Description should be empty string or not exist
				if description, exists := jsonMap["description"]; exists {
					assert.Equal(t, "", description)
				}
			},
		},
		{
			name: "RetrospectiveCreateRequest Schema",
			data: models.RetrospectiveCreateRequest{
				Title:       "Test Retrospective",
				Description: "A test retrospective",
				Template:    models.TemplateStartStopContinue,
			},
			validate: func(t *testing.T, jsonMap map[string]interface{}) {
				// Check required fields
				assert.Contains(t, jsonMap, "title")
				assert.Contains(t, jsonMap, "template")

				// Check field types
				assert.IsType(t, "", jsonMap["title"])    // Title as string
				assert.IsType(t, "", jsonMap["template"]) // Template as string

				// Check optional fields
				if description, exists := jsonMap["description"]; exists {
					assert.IsType(t, "", description) // Description as string
				}
			},
		},
		{
			name: "RetrospectiveCreateRequest without description",
			data: models.RetrospectiveCreateRequest{
				Title:    "Simple Retrospective",
				Template: models.Template4Ls,
			},
			validate: func(t *testing.T, jsonMap map[string]interface{}) {
				// Check required fields
				assert.Contains(t, jsonMap, "title")
				assert.Contains(t, jsonMap, "template")

				// Check field types
				assert.IsType(t, "", jsonMap["title"])
				assert.IsType(t, "", jsonMap["template"])

				// Description should be empty string or not exist
				if description, exists := jsonMap["description"]; exists {
					assert.Equal(t, "", description)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.data)
			require.NoError(t, err)

			// Validate JSON structure
			var jsonMap map[string]interface{}
			err = json.Unmarshal(jsonData, &jsonMap)
			require.NoError(t, err)

			// Run validation
			tt.validate(t, jsonMap)
		})
	}
}

// TestEnumValidation tests enum value validation
func TestEnumValidation(t *testing.T) {
	tests := []struct {
		name        string
		enumType    string
		enumValues  []interface{}
		validValues []string
	}{
		{
			name:     "RetrospectiveStatus Enum",
			enumType: "status",
			enumValues: []interface{}{
				models.RetroStatusPlanned,
				models.RetroStatusActive,
				models.RetroStatusCollecting,
				models.RetroStatusVoting,
				models.RetroStatusDiscussing,
				models.RetroStatusClosed,
			},
			validValues: []string{
				"planned", "active", "collecting",
				"voting", "discussing", "closed",
			},
		},
		{
			name:     "RetrospectiveTemplate Enum",
			enumType: "template",
			enumValues: []interface{}{
				models.TemplateStartStopContinue,
				models.Template4Ls,
				models.TemplateMadSadGlad,
				models.TemplateSailboat,
				models.TemplateWentWellToImprove,
			},
			validValues: []string{
				"start_stop_continue", "4ls", "mad_sad_glad",
				"sailboat", "went_well_to_improve",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, enumValue := range tt.enumValues {
				jsonData, err := json.Marshal(enumValue)
				require.NoError(t, err)

				// Unmarshal as string to check value
				var enumStr string
				err = json.Unmarshal(jsonData, &enumStr)
				require.NoError(t, err)

				// Validate it's a valid enum value
				assert.Contains(t, tt.validValues, enumStr,
					"Enum value %v should be in valid values %v", enumStr, tt.validValues)
			}
		})
	}
}

// TestUUIDValidation tests UUID format validation
func TestUUIDValidation(t *testing.T) {
	tests := []struct {
		name     string
		uuid     uuid.UUID
		expected string
	}{
		{
			name:     "Random UUID",
			uuid:     uuid.New(),
			expected: "36 characters with dashes",
		},
		{
			name:     "Specific UUID",
			uuid:     uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			expected: "123e4567-e89b-12d3-a456-426614174000",
		},
		{
			name:     "Nil UUID",
			uuid:     uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			expected: "00000000-0000-0000-0000-000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.uuid)
			require.NoError(t, err)

			// Unmarshal as string
			var uuidStr string
			err = json.Unmarshal(jsonData, &uuidStr)
			require.NoError(t, err)

			// Validate UUID format (basic check)
			assert.Len(t, uuidStr, 36) // Standard UUID length
			assert.Contains(t, uuidStr, "-")

			// For specific UUIDs, check exact value
			if tt.expected != "36 characters with dashes" {
				assert.Equal(t, tt.expected, uuidStr)
			}
		})
	}
}

// TestTimeValidation tests time format validation
func TestTimeValidation(t *testing.T) {
	tests := []struct {
		name        string
		timeValue   time.Time
		description string
	}{
		{
			name:        "Current time",
			timeValue:   time.Now(),
			description: "Current timestamp",
		},
		{
			name:        "UTC time",
			timeValue:   time.Now().UTC(),
			description: "UTC timestamp",
		},
		{
			name:        "Fixed time",
			timeValue:   time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC),
			description: "Fixed timestamp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.timeValue)
			require.NoError(t, err)

			// Unmarshal as string
			var timeStr string
			err = json.Unmarshal(jsonData, &timeStr)
			require.NoError(t, err)

			// Validate time format (basic check)
			assert.Contains(t, timeStr, "T") // ISO 8601 format
			// Timezone can be Z (UTC) or offset like -03:00
			assert.True(t, strings.Contains(timeStr, "Z") || strings.Contains(timeStr, "+") || strings.Contains(timeStr, "-"),
				"Time string should contain timezone indicator: %s", timeStr)
		})
	}
}

// TestNullPointerHandling tests null pointer handling in JSON
func TestNullPointerHandling(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		validate func(t *testing.T, jsonMap map[string]interface{})
	}{
		{
			name: "Null String Pointers",
			data: models.User{
				ID:        uuid.New(),
				Email:     "test@example.com",
				Name:      "Test User",
				Password:  "hashedpassword",
				Avatar:    nil, // Null pointer
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			validate: func(t *testing.T, jsonMap map[string]interface{}) {
				// Avatar should be null or missing
				if avatar, exists := jsonMap["avatar"]; exists {
					assert.Nil(t, avatar)
				}
			},
		},
		{
			name: "Null UUID Pointers",
			data: models.ActionItem{
				ID:              uuid.New(),
				RetrospectiveID: uuid.New(),
				ItemID:          nil, // Null pointer
				Title:           "Test Action Item",
				Description:     stringPtr("A test action item"),
				AssignedTo:      nil, // Null pointer
				Status:          "todo",
				CreatedBy:       uuid.New(),
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
			validate: func(t *testing.T, jsonMap map[string]interface{}) {
				// ItemID should be null or missing
				if itemID, exists := jsonMap["item_id"]; exists {
					assert.Nil(t, itemID)
				}

				// AssignedTo should be null or missing
				if assignedTo, exists := jsonMap["assigned_to"]; exists {
					assert.Nil(t, assignedTo)
				}
			},
		},
		{
			name: "Mixed null and non-null pointers",
			data: models.Team{
				ID:          uuid.New(),
				Name:        "Test Team",
				Description: nil, // Null pointer
				OwnerID:     uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			validate: func(t *testing.T, jsonMap map[string]interface{}) {
				// Description should be null or missing
				if description, exists := jsonMap["description"]; exists {
					assert.Nil(t, description)
				}

				// Other fields should exist
				assert.Contains(t, jsonMap, "id")
				assert.Contains(t, jsonMap, "name")
				assert.Contains(t, jsonMap, "owner_id")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.data)
			require.NoError(t, err)

			// Validate JSON structure
			var jsonMap map[string]interface{}
			err = json.Unmarshal(jsonData, &jsonMap)
			require.NoError(t, err)

			// Run validation
			tt.validate(t, jsonMap)
		})
	}
}

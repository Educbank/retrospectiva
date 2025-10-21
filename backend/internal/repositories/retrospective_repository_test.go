package repositories

import (
	"database/sql"
	"testing"
	"time"

	"educ-retro/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRetrospectiveRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRetrospectiveRepository(db)
	description := "Test Description"
	retrospective := &models.Retrospective{
		Title:       "Test Retrospective",
		Description: &description,
		Template:    "Went Well_to_improve",
		Status:      models.RetroStatusPlanned,
		CreatedBy:   uuid.New(),
		TeamID:      uuid.New(),
	}

	mock.ExpectQuery(`INSERT INTO retrospectives`).
		WithArgs(sqlmock.AnyArg(), retrospective.Title, retrospective.Description, retrospective.Template, retrospective.Status, retrospective.CreatedBy).
		WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at"}).
			AddRow(time.Now(), time.Now()))

	err = repo.Create(retrospective)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, retrospective.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetrospectiveRepository_Create_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRetrospectiveRepository(db)
	description := "Test Description"
	retrospective := &models.Retrospective{
		Title:       "Test Retrospective",
		Description: &description,
		Template:    "Went Well_to_improve",
		Status:      models.RetroStatusPlanned,
		CreatedBy:   uuid.New(),
		TeamID:      uuid.New(),
	}

	mock.ExpectQuery(`INSERT INTO retrospectives`).
		WithArgs(sqlmock.AnyArg(), retrospective.Title, retrospective.Description, retrospective.Template, retrospective.Status, retrospective.CreatedBy).
		WillReturnError(sql.ErrConnDone)

	err = repo.Create(retrospective)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetrospectiveRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRetrospectiveRepository(db)
	retrospectiveID := uuid.New()
	teamID := uuid.New()
	description := "Test Description"
	expectedRetrospective := &models.Retrospective{
		ID:          retrospectiveID,
		Title:       "Test Retrospective",
		Description: &description,
		Template:    "Went Well_to_improve",
		Status:      models.RetroStatusActive,
		CreatedBy:   uuid.New(),
		TeamID:      teamID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mock.ExpectQuery(`SELECT.*FROM retrospectives WHERE id`).
		WithArgs(retrospectiveID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "team_id", "title", "description", "template", "status", "scheduled_at", "started_at", "ended_at", "created_by", "created_at", "updated_at"}).
			AddRow(expectedRetrospective.ID, expectedRetrospective.TeamID, expectedRetrospective.Title, expectedRetrospective.Description, expectedRetrospective.Template, expectedRetrospective.Status, nil, nil, nil, expectedRetrospective.CreatedBy, expectedRetrospective.CreatedAt, expectedRetrospective.UpdatedAt))

	retrospective, err := repo.GetByID(retrospectiveID)

	assert.NoError(t, err)
	assert.Equal(t, expectedRetrospective.ID, retrospective.ID)
	assert.Equal(t, expectedRetrospective.Title, retrospective.Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetrospectiveRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRetrospectiveRepository(db)
	retrospectiveID := uuid.New()

	mock.ExpectQuery(`SELECT.*FROM retrospectives WHERE id`).
		WithArgs(retrospectiveID).
		WillReturnError(sql.ErrNoRows)

	retrospective, err := repo.GetByID(retrospectiveID)

	assert.Error(t, err)
	assert.Nil(t, retrospective)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetrospectiveRepository_GetByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRetrospectiveRepository(db)
	userID := uuid.New()
	teamID := uuid.New()
	description := "Test Description"
	expectedRetrospective := models.Retrospective{
		ID:          uuid.New(),
		Title:       "Test Retrospective",
		Description: &description,
		Template:    "Went Well_to_improve",
		Status:      models.RetroStatusActive,
		CreatedBy:   userID,
		TeamID:      teamID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mock.ExpectQuery(`SELECT.*FROM retrospectives`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "team_id", "title", "description", "template", "status", "scheduled_at", "started_at", "ended_at", "created_by", "created_at", "updated_at"}).
			AddRow(expectedRetrospective.ID, expectedRetrospective.TeamID, expectedRetrospective.Title, expectedRetrospective.Description, expectedRetrospective.Template, expectedRetrospective.Status, nil, nil, nil, expectedRetrospective.CreatedBy, expectedRetrospective.CreatedAt, expectedRetrospective.UpdatedAt))

	retrospectives, err := repo.GetByUserID(userID)

	assert.NoError(t, err)
	assert.Len(t, retrospectives, 1)
	assert.Equal(t, expectedRetrospective.Title, retrospectives[0].Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetrospectiveRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRetrospectiveRepository(db)
	retrospectiveID := uuid.New()

	mock.ExpectExec(`DELETE FROM retrospectives WHERE id`).
		WithArgs(retrospectiveID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(retrospectiveID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetrospectiveRepository_Delete_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRetrospectiveRepository(db)
	retrospectiveID := uuid.New()

	mock.ExpectExec(`DELETE FROM retrospectives WHERE id`).
		WithArgs(retrospectiveID).
		WillReturnError(sql.ErrConnDone)

	err = repo.Delete(retrospectiveID)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetrospectiveRepository_ReopenRetrospective(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRetrospectiveRepository(db)
	retrospectiveID := uuid.New()

	mock.ExpectExec(`UPDATE retrospectives`).
		WithArgs(retrospectiveID, models.RetroStatusActive).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.ReopenRetrospective(retrospectiveID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetrospectiveRepository_ReopenRetrospective_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRetrospectiveRepository(db)
	retrospectiveID := uuid.New()

	mock.ExpectExec(`UPDATE retrospectives`).
		WithArgs(retrospectiveID, models.RetroStatusActive).
		WillReturnError(sql.ErrConnDone)

	err = repo.ReopenRetrospective(retrospectiveID)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

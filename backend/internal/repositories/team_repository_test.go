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

func TestTeamRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTeamRepository(db)
	description := "Test Description"
	team := &models.Team{
		Name:        "Test Team",
		Description: &description,
		OwnerID:     uuid.New(),
	}

	mock.ExpectQuery(`INSERT INTO teams`).
		WithArgs(sqlmock.AnyArg(), team.Name, team.Description, team.OwnerID).
		WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at"}).
			AddRow(time.Now(), time.Now()))

	mock.ExpectExec(`INSERT INTO team_members`).
		WithArgs(sqlmock.AnyArg(), team.OwnerID, "owner").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(team)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, team.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_Create_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTeamRepository(db)
	description := "Test Description"
	team := &models.Team{
		Name:        "Test Team",
		Description: &description,
		OwnerID:     uuid.New(),
	}

	mock.ExpectQuery(`INSERT INTO teams`).
		WithArgs(sqlmock.AnyArg(), team.Name, team.Description, team.OwnerID).
		WillReturnError(sql.ErrConnDone)

	err = repo.Create(team)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTeamRepository(db)
	teamID := uuid.New()
	description := "Test Description"
	expectedTeam := &models.Team{
		ID:          teamID,
		Name:        "Test Team",
		Description: &description,
		OwnerID:     uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mock.ExpectQuery(`SELECT.*FROM teams WHERE id`).
		WithArgs(teamID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "owner_id", "created_at", "updated_at"}).
			AddRow(expectedTeam.ID, expectedTeam.Name, expectedTeam.Description, expectedTeam.OwnerID, expectedTeam.CreatedAt, expectedTeam.UpdatedAt))

	team, err := repo.GetByID(teamID)

	assert.NoError(t, err)
	assert.Equal(t, expectedTeam.ID, team.ID)
	assert.Equal(t, expectedTeam.Name, team.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTeamRepository(db)
	teamID := uuid.New()

	mock.ExpectQuery(`SELECT.*FROM teams WHERE id`).
		WithArgs(teamID).
		WillReturnError(sql.ErrNoRows)

	team, err := repo.GetByID(teamID)

	assert.Error(t, err)
	assert.Nil(t, team)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_GetByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTeamRepository(db)
	userID := uuid.New()
	description := "Test Description"
	expectedTeam := models.TeamWithCounts{
		Team: models.Team{
			ID:          uuid.New(),
			Name:        "Test Team",
			Description: &description,
			OwnerID:     userID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		MemberCount:        5,
		RetrospectiveCount: 3,
	}

	mock.ExpectQuery(`SELECT.*FROM teams t`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "owner_id", "created_at", "updated_at", "member_count", "retrospective_count"}).
			AddRow(expectedTeam.ID, expectedTeam.Name, expectedTeam.Description, expectedTeam.OwnerID, expectedTeam.CreatedAt, expectedTeam.UpdatedAt, expectedTeam.MemberCount, expectedTeam.RetrospectiveCount))

	teams, err := repo.GetByUserID(userID)

	assert.NoError(t, err)
	assert.Len(t, teams, 1)
	assert.Equal(t, expectedTeam.Name, teams[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_GetByUserID_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTeamRepository(db)
	userID := uuid.New()

	mock.ExpectQuery(`SELECT.*FROM teams t`).
		WithArgs(userID).
		WillReturnError(sql.ErrConnDone)

	teams, err := repo.GetByUserID(userID)

	assert.Error(t, err)
	assert.Nil(t, teams)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTeamRepository(db)
	description := "Updated Description"
	team := &models.Team{
		ID:          uuid.New(),
		Name:        "Updated Team",
		Description: &description,
	}

	mock.ExpectQuery(`UPDATE teams`).
		WithArgs(team.ID, team.Name, team.Description).
		WillReturnRows(sqlmock.NewRows([]string{"updated_at"}).
			AddRow(time.Now()))

	err = repo.Update(team)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_Update_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTeamRepository(db)
	description := "Updated Description"
	team := &models.Team{
		ID:          uuid.New(),
		Name:        "Updated Team",
		Description: &description,
	}

	mock.ExpectQuery(`UPDATE teams`).
		WithArgs(team.ID, team.Name, team.Description).
		WillReturnError(sql.ErrConnDone)

	err = repo.Update(team)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTeamRepository(db)
	teamID := uuid.New()

	mock.ExpectExec(`DELETE FROM teams WHERE id`).
		WithArgs(teamID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(teamID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_Delete_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTeamRepository(db)
	teamID := uuid.New()

	mock.ExpectExec(`DELETE FROM teams WHERE id`).
		WithArgs(teamID).
		WillReturnError(sql.ErrConnDone)

	err = repo.Delete(teamID)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_AddMember(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTeamRepository(db)
	teamID := uuid.New()
	userID := uuid.New()
	role := "member"

	mock.ExpectExec(`INSERT INTO team_members`).
		WithArgs(teamID, userID, role).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.AddMember(teamID, userID, role)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_AddMember_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTeamRepository(db)
	teamID := uuid.New()
	userID := uuid.New()
	role := "member"

	mock.ExpectExec(`INSERT INTO team_members`).
		WithArgs(teamID, userID, role).
		WillReturnError(sql.ErrConnDone)

	err = repo.AddMember(teamID, userID, role)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_GetMembers(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTeamRepository(db)
	teamID := uuid.New()
	expectedMember := models.TeamMemberWithUser{
		ID:        uuid.New(),
		TeamID:    teamID,
		UserID:    uuid.New(),
		UserName:  "Test User",
		UserEmail: "test@example.com",
		Role:      "member",
		JoinedAt:  time.Now(),
	}

	mock.ExpectQuery(`SELECT.*FROM team_members tm`).
		WithArgs(teamID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "team_id", "user_id", "user_name", "user_email", "role", "joined_at"}).
			AddRow(expectedMember.ID, expectedMember.TeamID, expectedMember.UserID, expectedMember.UserName, expectedMember.UserEmail, expectedMember.Role, expectedMember.JoinedAt))

	members, err := repo.GetMembers(teamID)

	assert.NoError(t, err)
	assert.Len(t, members, 1)
	assert.Equal(t, expectedMember.UserName, members[0].UserName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_GetMembers_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTeamRepository(db)
	teamID := uuid.New()

	mock.ExpectQuery(`SELECT.*FROM team_members tm`).
		WithArgs(teamID).
		WillReturnError(sql.ErrConnDone)

	members, err := repo.GetMembers(teamID)

	assert.Error(t, err)
	assert.Nil(t, members)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_RemoveMember(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTeamRepository(db)
	teamID := uuid.New()
	userID := uuid.New()

	mock.ExpectExec(`DELETE FROM team_members WHERE team_id`).
		WithArgs(teamID, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.RemoveMember(teamID, userID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_RemoveMember_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTeamRepository(db)
	teamID := uuid.New()
	userID := uuid.New()

	mock.ExpectExec(`DELETE FROM team_members WHERE team_id`).
		WithArgs(teamID, userID).
		WillReturnError(sql.ErrConnDone)

	err = repo.RemoveMember(teamID, userID)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_UpdateMemberRole(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTeamRepository(db)
	teamID := uuid.New()
	userID := uuid.New()
	role := "admin"

	mock.ExpectExec(`UPDATE team_members`).
		WithArgs(teamID, userID, role).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateMemberRole(teamID, userID, role)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_UpdateMemberRole_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTeamRepository(db)
	teamID := uuid.New()
	userID := uuid.New()
	role := "admin"

	mock.ExpectExec(`UPDATE team_members`).
		WithArgs(teamID, userID, role).
		WillReturnError(sql.ErrConnDone)

	err = repo.UpdateMemberRole(teamID, userID, role)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_IsMember(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTeamRepository(db)
	teamID := uuid.New()
	userID := uuid.New()

	mock.ExpectQuery(`SELECT COUNT.*FROM team_members`).
		WithArgs(teamID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	isMember, err := repo.IsMember(teamID, userID)

	assert.NoError(t, err)
	assert.True(t, isMember)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_IsMember_NotMember(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTeamRepository(db)
	teamID := uuid.New()
	userID := uuid.New()

	mock.ExpectQuery(`SELECT COUNT.*FROM team_members`).
		WithArgs(teamID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	isMember, err := repo.IsMember(teamID, userID)

	assert.NoError(t, err)
	assert.False(t, isMember)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_IsMember_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTeamRepository(db)
	teamID := uuid.New()
	userID := uuid.New()

	mock.ExpectQuery(`SELECT COUNT.*FROM team_members`).
		WithArgs(teamID, userID).
		WillReturnError(sql.ErrConnDone)

	isMember, err := repo.IsMember(teamID, userID)

	assert.Error(t, err)
	assert.False(t, isMember)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_GetMemberRole(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTeamRepository(db)
	teamID := uuid.New()
	userID := uuid.New()
	expectedRole := "admin"

	mock.ExpectQuery(`SELECT role FROM team_members`).
		WithArgs(teamID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow(expectedRole))

	role, err := repo.GetMemberRole(teamID, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedRole, role)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_GetMemberRole_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTeamRepository(db)
	teamID := uuid.New()
	userID := uuid.New()

	mock.ExpectQuery(`SELECT role FROM team_members`).
		WithArgs(teamID, userID).
		WillReturnError(sql.ErrNoRows)

	role, err := repo.GetMemberRole(teamID, userID)

	assert.Error(t, err)
	assert.Empty(t, role)
	assert.NoError(t, mock.ExpectationsWereMet())
}

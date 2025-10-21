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

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	avatar := "avatar.jpg"
	user := &models.User{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "hashedpassword",
		Avatar:   &avatar,
	}

	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(sqlmock.AnyArg(), user.Email, user.Name, user.Password, user.Avatar).
		WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at"}).
			AddRow(time.Now(), time.Now()))

	err = repo.Create(user)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, user.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Create_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	avatar := "avatar.jpg"
	user := &models.User{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "hashedpassword",
		Avatar:   &avatar,
	}

	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(sqlmock.AnyArg(), user.Email, user.Name, user.Password, user.Avatar).
		WillReturnError(sql.ErrConnDone)

	err = repo.Create(user)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	userID := uuid.New()
	avatar := "avatar.jpg"
	expectedUser := &models.User{
		ID:        userID,
		Email:     "test@example.com",
		Name:      "Test User",
		Password:  "hashedpassword",
		Avatar:    &avatar,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectQuery(`SELECT.*FROM users WHERE id`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name", "password", "avatar", "created_at", "updated_at"}).
			AddRow(expectedUser.ID, expectedUser.Email, expectedUser.Name, expectedUser.Password, expectedUser.Avatar, expectedUser.CreatedAt, expectedUser.UpdatedAt))

	user, err := repo.GetByID(userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	userID := uuid.New()

	mock.ExpectQuery(`SELECT.*FROM users WHERE id`).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetByID(userID)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	email := "test@example.com"
	avatar := "avatar.jpg"
	expectedUser := &models.User{
		ID:        uuid.New(),
		Email:     email,
		Name:      "Test User",
		Password:  "hashedpassword",
		Avatar:    &avatar,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectQuery(`SELECT.*FROM users WHERE email`).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name", "password", "avatar", "created_at", "updated_at"}).
			AddRow(expectedUser.ID, expectedUser.Email, expectedUser.Name, expectedUser.Password, expectedUser.Avatar, expectedUser.CreatedAt, expectedUser.UpdatedAt))

	user, err := repo.GetByEmail(email)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	email := "test@example.com"

	mock.ExpectQuery(`SELECT.*FROM users WHERE email`).
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetByEmail(email)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	avatar := "new-avatar.jpg"
	user := &models.User{
		ID:     uuid.New(),
		Email:  "updated@example.com",
		Name:   "Updated User",
		Avatar: &avatar,
	}

	mock.ExpectQuery(`UPDATE users`).
		WithArgs(user.ID, user.Email, user.Name, user.Avatar).
		WillReturnRows(sqlmock.NewRows([]string{"updated_at"}).
			AddRow(time.Now()))

	err = repo.Update(user)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Update_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	avatar := "new-avatar.jpg"
	user := &models.User{
		ID:     uuid.New(),
		Email:  "updated@example.com",
		Name:   "Updated User",
		Avatar: &avatar,
	}

	mock.ExpectQuery(`UPDATE users`).
		WithArgs(user.ID, user.Email, user.Name, user.Avatar).
		WillReturnError(sql.ErrConnDone)

	err = repo.Update(user)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	userID := uuid.New()

	mock.ExpectExec(`DELETE FROM users WHERE id`).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(userID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Delete_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	userID := uuid.New()

	mock.ExpectExec(`DELETE FROM users WHERE id`).
		WithArgs(userID).
		WillReturnError(sql.ErrConnDone)

	err = repo.Delete(userID)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUsersByIDs(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	userIDs := []uuid.UUID{uuid.New(), uuid.New()}
	expectedUsers := []models.User{
		{ID: userIDs[0], Email: "user1@example.com", Name: "User 1"},
		{ID: userIDs[1], Email: "user2@example.com", Name: "User 2"},
	}

	mock.ExpectQuery(`SELECT.*FROM users WHERE id IN`).
		WithArgs(userIDs[0], userIDs[1]).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name", "avatar", "created_at", "updated_at"}).
			AddRow(expectedUsers[0].ID, expectedUsers[0].Email, expectedUsers[0].Name, "", time.Now(), time.Now()).
			AddRow(expectedUsers[1].ID, expectedUsers[1].Email, expectedUsers[1].Name, "", time.Now(), time.Now()))

	users, err := repo.GetUsersByIDs(userIDs)

	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUsersByIDs_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	users, err := repo.GetUsersByIDs([]uuid.UUID{})

	assert.NoError(t, err)
	assert.Empty(t, users)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUsersByIDs_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	userIDs := []uuid.UUID{uuid.New()}

	mock.ExpectQuery(`SELECT.*FROM users WHERE id IN`).
		WithArgs(userIDs[0]).
		WillReturnError(sql.ErrConnDone)

	users, err := repo.GetUsersByIDs(userIDs)

	assert.Error(t, err)
	assert.Nil(t, users)
	assert.NoError(t, mock.ExpectationsWereMet())
}

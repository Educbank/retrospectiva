package repositories

import (
	"database/sql"

	"educ-retro/internal/models"

	"github.com/google/uuid"
)

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Create(team *models.Team) error {
	query := `
		INSERT INTO teams (id, name, description, owner_id)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at, updated_at
	`

	team.ID = uuid.New()
	err := r.db.QueryRow(query, team.ID, team.Name, team.Description, team.OwnerID).
		Scan(&team.CreatedAt, &team.UpdatedAt)

	if err != nil {
		return err
	}

	// Add owner as team member
	return r.AddMember(team.ID, team.OwnerID, "owner")
}

func (r *TeamRepository) GetByID(id uuid.UUID) (*models.Team, error) {
	query := `
		SELECT id, name, description, owner_id, created_at, updated_at
		FROM teams WHERE id = $1
	`

	team := &models.Team{}
	err := r.db.QueryRow(query, id).Scan(
		&team.ID, &team.Name, &team.Description, &team.OwnerID,
		&team.CreatedAt, &team.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return team, nil
}

func (r *TeamRepository) GetByUserID(userID uuid.UUID) ([]models.TeamWithCounts, error) {
	query := `
		SELECT 
			t.id, 
			t.name, 
			t.description, 
			t.owner_id, 
			t.created_at, 
			t.updated_at,
			COALESCE(member_count.count, 0) as member_count,
			COALESCE(retro_count.count, 0) as retrospective_count
		FROM teams t
		INNER JOIN team_members tm ON t.id = tm.team_id
		LEFT JOIN (
			SELECT team_id, COUNT(*) as count 
			FROM team_members 
			GROUP BY team_id
		) member_count ON t.id = member_count.team_id
		LEFT JOIN (
			SELECT team_id, COUNT(*) as count 
			FROM retrospectives 
			GROUP BY team_id
		) retro_count ON t.id = retro_count.team_id
		WHERE tm.user_id = $1
		ORDER BY t.created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []models.TeamWithCounts
	for rows.Next() {
		var team models.TeamWithCounts
		err := rows.Scan(
			&team.ID, &team.Name, &team.Description, &team.OwnerID,
			&team.CreatedAt, &team.UpdatedAt,
			&team.MemberCount, &team.RetrospectiveCount,
		)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	return teams, nil
}

func (r *TeamRepository) Update(team *models.Team) error {
	query := `
		UPDATE teams 
		SET name = $2, description = $3, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRow(query, team.ID, team.Name, team.Description).
		Scan(&team.UpdatedAt)

	return err
}

func (r *TeamRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM teams WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *TeamRepository) AddMember(teamID, userID uuid.UUID, role string) error {
	query := `
		INSERT INTO team_members (team_id, user_id, role)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(query, teamID, userID, role)
	return err
}

func (r *TeamRepository) GetMembers(teamID uuid.UUID) ([]models.TeamMemberWithUser, error) {
	query := `
		SELECT 
			tm.id, 
			tm.team_id, 
			tm.user_id, 
			u.name as user_name,
			u.email as user_email,
			tm.role, 
			tm.joined_at
		FROM team_members tm
		INNER JOIN users u ON tm.user_id = u.id
		WHERE tm.team_id = $1
		ORDER BY tm.joined_at ASC
	`

	rows, err := r.db.Query(query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.TeamMemberWithUser
	for rows.Next() {
		var member models.TeamMemberWithUser
		err := rows.Scan(
			&member.ID, &member.TeamID, &member.UserID,
			&member.UserName, &member.UserEmail,
			&member.Role, &member.JoinedAt,
		)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	return members, nil
}

func (r *TeamRepository) RemoveMember(teamID, userID uuid.UUID) error {
	query := `DELETE FROM team_members WHERE team_id = $1 AND user_id = $2`
	_, err := r.db.Exec(query, teamID, userID)
	return err
}

func (r *TeamRepository) UpdateMemberRole(teamID, userID uuid.UUID, role string) error {
	query := `
		UPDATE team_members 
		SET role = $3
		WHERE team_id = $1 AND user_id = $2
	`
	_, err := r.db.Exec(query, teamID, userID, role)
	return err
}

func (r *TeamRepository) IsMember(teamID, userID uuid.UUID) (bool, error) {
	query := `
		SELECT COUNT(*) FROM team_members
		WHERE team_id = $1 AND user_id = $2
	`

	var count int
	err := r.db.QueryRow(query, teamID, userID).Scan(&count)
	return count > 0, err
}

func (r *TeamRepository) GetMemberRole(teamID, userID uuid.UUID) (string, error) {
	query := `
		SELECT role FROM team_members
		WHERE team_id = $1 AND user_id = $2
	`

	var role string
	err := r.db.QueryRow(query, teamID, userID).Scan(&role)
	return role, err
}

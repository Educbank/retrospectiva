package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"educ-retro/internal/models"

	"github.com/google/uuid"
)

type RetrospectiveRepository struct {
	db *sql.DB
}

func NewRetrospectiveRepository(db *sql.DB) *RetrospectiveRepository {
	return &RetrospectiveRepository{db: db}
}

func (r *RetrospectiveRepository) Create(retrospective *models.Retrospective) error {
	query := `
		INSERT INTO retrospectives (id, title, description, template, status, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`

	retrospective.ID = uuid.New()
	err := r.db.QueryRow(query,
		retrospective.ID,
		retrospective.Title,
		retrospective.Description,
		retrospective.Template,
		retrospective.Status,
		retrospective.CreatedBy,
	).Scan(&retrospective.CreatedAt, &retrospective.UpdatedAt)

	return err
}

func (r *RetrospectiveRepository) GetByID(id uuid.UUID) (*models.Retrospective, error) {
	query := `
		SELECT id, team_id, title, description, template, status, scheduled_at, started_at, ended_at, 
		       timer_duration, timer_started_at, timer_paused_at, timer_elapsed_time,
		       created_by, created_at, updated_at
		FROM retrospectives WHERE id = $1
	`

	retrospective := &models.Retrospective{}
	err := r.db.QueryRow(query, id).Scan(
		&retrospective.ID,
		&retrospective.TeamID,
		&retrospective.Title,
		&retrospective.Description,
		&retrospective.Template,
		&retrospective.Status,
		&retrospective.ScheduledAt,
		&retrospective.StartedAt,
		&retrospective.EndedAt,
		&retrospective.TimerDuration,
		&retrospective.TimerStartedAt,
		&retrospective.TimerPausedAt,
		&retrospective.TimerElapsedTime,
		&retrospective.CreatedBy,
		&retrospective.CreatedAt,
		&retrospective.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return retrospective, nil
}

func (r *RetrospectiveRepository) GetByTeamID(teamID uuid.UUID) ([]models.Retrospective, error) {
	query := `
		SELECT id, team_id, title, description, template, status, scheduled_at, started_at, ended_at, 
		       timer_duration, timer_started_at, timer_paused_at, timer_elapsed_time,
		       created_by, created_at, updated_at
		FROM retrospectives 
		WHERE team_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var retrospectives []models.Retrospective
	for rows.Next() {
		var retrospective models.Retrospective
		err := rows.Scan(
			&retrospective.ID,
			&retrospective.TeamID,
			&retrospective.Title,
			&retrospective.Description,
			&retrospective.Template,
			&retrospective.Status,
			&retrospective.ScheduledAt,
			&retrospective.StartedAt,
			&retrospective.EndedAt,
			&retrospective.TimerDuration,
			&retrospective.TimerStartedAt,
			&retrospective.TimerPausedAt,
			&retrospective.TimerElapsedTime,
			&retrospective.CreatedBy,
			&retrospective.CreatedAt,
			&retrospective.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		retrospectives = append(retrospectives, retrospective)
	}

	return retrospectives, nil
}

func (r *RetrospectiveRepository) GetByUserID(userID uuid.UUID) ([]models.Retrospective, error) {
	query := `
		SELECT id, team_id, title, description, template, status, scheduled_at, started_at, ended_at, 
		       timer_duration, timer_started_at, timer_paused_at, timer_elapsed_time,
		       created_by, created_at, updated_at
		FROM retrospectives
		WHERE created_by = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var retrospectives []models.Retrospective
	for rows.Next() {
		var retrospective models.Retrospective
		err := rows.Scan(
			&retrospective.ID,
			&retrospective.TeamID,
			&retrospective.Title,
			&retrospective.Description,
			&retrospective.Template,
			&retrospective.Status,
			&retrospective.ScheduledAt,
			&retrospective.StartedAt,
			&retrospective.EndedAt,
			&retrospective.TimerDuration,
			&retrospective.TimerStartedAt,
			&retrospective.TimerPausedAt,
			&retrospective.TimerElapsedTime,
			&retrospective.CreatedBy,
			&retrospective.CreatedAt,
			&retrospective.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		retrospectives = append(retrospectives, retrospective)
	}

	return retrospectives, nil
}

func (r *RetrospectiveRepository) GetAllRetrospectives() ([]models.Retrospective, error) {
	query := `
		SELECT id, team_id, title, description, template, status, scheduled_at, started_at, ended_at, 
		       timer_duration, timer_started_at, timer_paused_at, timer_elapsed_time,
		       created_by, created_at, updated_at
		FROM retrospectives
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var retrospectives []models.Retrospective
	for rows.Next() {
		var retrospective models.Retrospective
		err := rows.Scan(
			&retrospective.ID,
			&retrospective.TeamID,
			&retrospective.Title,
			&retrospective.Description,
			&retrospective.Template,
			&retrospective.Status,
			&retrospective.ScheduledAt,
			&retrospective.StartedAt,
			&retrospective.EndedAt,
			&retrospective.TimerDuration,
			&retrospective.TimerStartedAt,
			&retrospective.TimerPausedAt,
			&retrospective.TimerElapsedTime,
			&retrospective.CreatedBy,
			&retrospective.CreatedAt,
			&retrospective.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		retrospectives = append(retrospectives, retrospective)
	}

	return retrospectives, nil
}

func (r *RetrospectiveRepository) Update(retrospective *models.Retrospective) error {
	query := `
		UPDATE retrospectives 
		SET title = $2, description = $3, template = $4, status = $5, started_at = $6, ended_at = $7, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRow(query,
		retrospective.ID,
		retrospective.Title,
		retrospective.Description,
		retrospective.Template,
		retrospective.Status,
		retrospective.StartedAt,
		retrospective.EndedAt,
	).Scan(&retrospective.UpdatedAt)

	return err
}

func (r *RetrospectiveRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM retrospectives WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *RetrospectiveRepository) UpdateStatus(id uuid.UUID, status models.RetrospectiveStatus) error {
	var query string
	var args []interface{}

	switch status {
	case models.RetroStatusActive:
		query = `
			UPDATE retrospectives 
			SET status = $2, started_at = NOW(), updated_at = NOW()
			WHERE id = $1
		`
		args = []interface{}{id, status}
	case models.RetroStatusClosed:
		query = `
			UPDATE retrospectives 
			SET status = $2, ended_at = NOW(), updated_at = NOW()
			WHERE id = $1
		`
		args = []interface{}{id, status}
	default:
		query = `
			UPDATE retrospectives 
			SET status = $2, updated_at = NOW()
			WHERE id = $1
		`
		args = []interface{}{id, status}
	}

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *RetrospectiveRepository) GetRetrospectiveCount(userID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM retrospectives r
		INNER JOIN team_members tm ON r.team_id = tm.team_id
		WHERE tm.user_id = $1
	`

	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	return count, err
}

func (r *RetrospectiveRepository) GetActionItemCount(userID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM action_items ai
		INNER JOIN retrospectives r ON ai.retrospective_id = r.id
		WHERE r.created_by = $1
	`

	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	return count, err
}

func (r *RetrospectiveRepository) AddItem(item *models.RetrospectiveItem) error {
	query := `
		INSERT INTO retrospective_items (id, retrospective_id, category, content, author_id, is_anonymous, votes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		item.ID,
		item.RetrospectiveID,
		item.Category,
		item.Content,
		item.AuthorID,
		item.IsAnonymous,
		item.Votes,
	).Scan(&item.CreatedAt, &item.UpdatedAt)

	return err
}

func (r *RetrospectiveRepository) VoteItem(itemID, userID uuid.UUID) error {
	// First, check if user already voted
	var count int
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM retrospective_votes WHERE item_id = $1 AND user_id = $2",
		itemID, userID,
	).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		// User already voted, remove the vote
		_, err = r.db.Exec(
			"DELETE FROM retrospective_votes WHERE item_id = $1 AND user_id = $2",
			itemID, userID,
		)
		if err != nil {
			return err
		}

		// Decrease vote count
		_, err = r.db.Exec(
			"UPDATE retrospective_items SET votes = votes - 1 WHERE id = $1",
			itemID,
		)
	} else {
		// Add vote
		_, err = r.db.Exec(
			"INSERT INTO retrospective_votes (id, item_id, user_id) VALUES ($1, $2, $3)",
			uuid.New(), itemID, userID,
		)
		if err != nil {
			return err
		}

		// Increase vote count
		_, err = r.db.Exec(
			"UPDATE retrospective_items SET votes = votes + 1 WHERE id = $1",
			itemID,
		)
	}

	return err
}

func (r *RetrospectiveRepository) AddActionItem(actionItem *models.ActionItem) error {
	query := `
		INSERT INTO action_items (id, retrospective_id, item_id, title, description, assigned_to, status, due_date, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		actionItem.ID,
		actionItem.RetrospectiveID,
		actionItem.ItemID,
		actionItem.Title,
		actionItem.Description,
		actionItem.AssignedTo,
		actionItem.Status,
		actionItem.DueDate,
		actionItem.CreatedBy,
	).Scan(&actionItem.CreatedAt, &actionItem.UpdatedAt)

	return err
}

func (r *RetrospectiveRepository) GetRetrospectiveWithDetails(retrospectiveID uuid.UUID) (*models.RetrospectiveWithDetails, error) {
	// Get retrospective
	retrospective, err := r.GetByID(retrospectiveID)
	if err != nil {
		return nil, err
	}

	// Get items
	items, err := r.GetItemsByRetrospectiveID(retrospectiveID)
	if err != nil {
		return nil, err
	}

	// Get action items
	actionItems, err := r.GetActionItemsByRetrospectiveID(retrospectiveID)
	if err != nil {
		return nil, err
	}

	// Get participants
	participants, err := r.GetParticipants(retrospectiveID)
	if err != nil {
		return nil, err
	}

	// Get groups
	groups, err := r.GetGroupsByRetrospectiveID(retrospectiveID)
	if err != nil {
		return nil, err
	}

	return &models.RetrospectiveWithDetails{
		Retrospective: *retrospective,
		Items:         items,
		ActionItems:   actionItems,
		Participants:  participants,
		Groups:        groups,
	}, nil
}

func (r *RetrospectiveRepository) GetItemsByRetrospectiveID(retrospectiveID uuid.UUID) ([]models.RetrospectiveItem, error) {
	query := `
		SELECT id, retrospective_id, category, content, author_id, is_anonymous, votes, created_at, updated_at
		FROM retrospective_items
		WHERE retrospective_id = $1
		ORDER BY votes DESC, created_at ASC
	`

	rows, err := r.db.Query(query, retrospectiveID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.RetrospectiveItem
	for rows.Next() {
		var item models.RetrospectiveItem
		err := rows.Scan(
			&item.ID,
			&item.RetrospectiveID,
			&item.Category,
			&item.Content,
			&item.AuthorID,
			&item.IsAnonymous,
			&item.Votes,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *RetrospectiveRepository) GetActionItemsByRetrospectiveID(retrospectiveID uuid.UUID) ([]models.ActionItem, error) {
	query := `
		SELECT id, retrospective_id, item_id, title, description, assigned_to, status, due_date, completed_at, created_by, created_at, updated_at
		FROM action_items
		WHERE retrospective_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query, retrospectiveID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actionItems []models.ActionItem
	for rows.Next() {
		var actionItem models.ActionItem
		err := rows.Scan(
			&actionItem.ID,
			&actionItem.RetrospectiveID,
			&actionItem.ItemID,
			&actionItem.Title,
			&actionItem.Description,
			&actionItem.AssignedTo,
			&actionItem.Status,
			&actionItem.DueDate,
			&actionItem.CompletedAt,
			&actionItem.CreatedBy,
			&actionItem.CreatedAt,
			&actionItem.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		actionItems = append(actionItems, actionItem)
	}

	return actionItems, nil
}

func (r *RetrospectiveRepository) RegisterParticipant(retrospectiveID, userID uuid.UUID) error {
	query := `
		INSERT INTO retrospective_participants (id, retrospective_id, user_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (retrospective_id, user_id)
		DO UPDATE SET last_seen = NOW()
	`

	_, err := r.db.Exec(query, uuid.New(), retrospectiveID, userID)
	return err
}

func (r *RetrospectiveRepository) GetParticipants(retrospectiveID uuid.UUID) ([]models.RetrospectiveParticipant, error) {
	query := `
		SELECT id, retrospective_id, user_id, joined_at, last_seen
		FROM retrospective_participants
		WHERE retrospective_id = $1
		ORDER BY joined_at ASC
	`

	rows, err := r.db.Query(query, retrospectiveID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []models.RetrospectiveParticipant
	for rows.Next() {
		var participant models.RetrospectiveParticipant
		err := rows.Scan(
			&participant.ID,
			&participant.RetrospectiveID,
			&participant.UserID,
			&participant.JoinedAt,
			&participant.LastSeen,
		)
		if err != nil {
			return nil, err
		}
		participants = append(participants, participant)
	}

	return participants, nil
}

func (r *RetrospectiveRepository) GetItemByID(itemID uuid.UUID) (*models.RetrospectiveItem, error) {
	query := `
		SELECT id, retrospective_id, category, content, author_id, is_anonymous, votes, created_at, updated_at
		FROM retrospective_items
		WHERE id = $1
	`

	item := &models.RetrospectiveItem{}
	err := r.db.QueryRow(query, itemID).Scan(
		&item.ID,
		&item.RetrospectiveID,
		&item.Category,
		&item.Content,
		&item.AuthorID,
		&item.IsAnonymous,
		&item.Votes,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *RetrospectiveRepository) DeleteItem(itemID uuid.UUID) error {
	query := `DELETE FROM retrospective_items WHERE id = $1`
	_, err := r.db.Exec(query, itemID)
	return err
}

// Group methods
func (r *RetrospectiveRepository) CreateGroup(group *models.RetrospectiveGroup, itemIDs []uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert group
	query := `
		INSERT INTO retrospective_groups (id, retrospective_id, name, description, created_by, votes)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`
	err = tx.QueryRow(query, group.ID, group.RetrospectiveID, group.Name, group.Description, group.CreatedBy, 0).
		Scan(&group.CreatedAt, &group.UpdatedAt)
	if err != nil {
		return err
	}

	// Insert group items
	for _, itemID := range itemIDs {
		itemQuery := `
			INSERT INTO retrospective_group_items (id, group_id, item_id)
			VALUES ($1, $2, $3)
		`
		_, err = tx.Exec(itemQuery, uuid.New(), group.ID, itemID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *RetrospectiveRepository) GetGroupsByRetrospectiveID(retrospectiveID uuid.UUID) ([]models.RetrospectiveGroup, error) {
	query := `
		SELECT id, retrospective_id, name, description, created_by, created_at, updated_at, votes
		FROM retrospective_groups
		WHERE retrospective_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.Query(query, retrospectiveID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []models.RetrospectiveGroup
	for rows.Next() {
		var group models.RetrospectiveGroup
		err := rows.Scan(&group.ID, &group.RetrospectiveID, &group.Name, &group.Description,
			&group.CreatedBy, &group.CreatedAt, &group.UpdatedAt, &group.Votes)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}

func (r *RetrospectiveRepository) VoteGroup(groupID, userID uuid.UUID) error {
	// Check if user already voted
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM retrospective_group_votes WHERE group_id = $1 AND user_id = $2)`
	err := r.db.QueryRow(checkQuery, groupID, userID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		// Remove vote
		_, err = r.db.Exec(`DELETE FROM retrospective_group_votes WHERE group_id = $1 AND user_id = $2`, groupID, userID)
		if err != nil {
			return err
		}
		// Decrease vote count
		_, err = r.db.Exec(`UPDATE retrospective_groups SET votes = votes - 1 WHERE id = $1`, groupID)
	} else {
		// Add vote
		_, err = r.db.Exec(`INSERT INTO retrospective_group_votes (id, group_id, user_id) VALUES ($1, $2, $3)`, uuid.New(), groupID, userID)
		if err != nil {
			return err
		}
		// Increase vote count
		_, err = r.db.Exec(`UPDATE retrospective_groups SET votes = votes + 1 WHERE id = $1`, groupID)
	}

	return err
}

func (r *RetrospectiveRepository) GetGroupByID(groupID uuid.UUID) (*models.RetrospectiveGroup, error) {
	query := `
		SELECT id, retrospective_id, name, description, votes, created_by, created_at, updated_at
		FROM retrospective_groups
		WHERE id = $1
	`
	var group models.RetrospectiveGroup
	err := r.db.QueryRow(query, groupID).Scan(
		&group.ID, &group.RetrospectiveID, &group.Name, &group.Description,
		&group.Votes, &group.CreatedBy, &group.CreatedAt, &group.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *RetrospectiveRepository) DeleteGroup(groupID uuid.UUID) error {
	query := `DELETE FROM retrospective_groups WHERE id = $1`
	_, err := r.db.Exec(query, groupID)
	return err
}

func (r *RetrospectiveRepository) MergeItems(sourceItemID, targetItemID uuid.UUID) (*models.RetrospectiveItem, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Get both items
	var sourceItem, targetItem models.RetrospectiveItem
	sourceQuery := `SELECT id, retrospective_id, category, content, author_id, is_anonymous, votes, created_at, updated_at FROM retrospective_items WHERE id = $1`
	targetQuery := `SELECT id, retrospective_id, category, content, author_id, is_anonymous, votes, created_at, updated_at FROM retrospective_items WHERE id = $1`

	err = tx.QueryRow(sourceQuery, sourceItemID).Scan(
		&sourceItem.ID, &sourceItem.RetrospectiveID, &sourceItem.Category, &sourceItem.Content,
		&sourceItem.AuthorID, &sourceItem.IsAnonymous, &sourceItem.Votes,
		&sourceItem.CreatedAt, &sourceItem.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(targetQuery, targetItemID).Scan(
		&targetItem.ID, &targetItem.RetrospectiveID, &targetItem.Category, &targetItem.Content,
		&targetItem.AuthorID, &targetItem.IsAnonymous, &targetItem.Votes,
		&targetItem.CreatedAt, &targetItem.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Check if items belong to same retrospective and category
	if sourceItem.RetrospectiveID != targetItem.RetrospectiveID {
		return nil, errors.New("items must belong to the same retrospective")
	}
	if sourceItem.Category != targetItem.Category {
		return nil, errors.New("items must belong to the same category")
	}

	// Merge content (combine both contents)
	mergedContent := targetItem.Content + " | " + sourceItem.Content

	// Discard votes when merging (reset to 0)
	mergedVotes := 0

	// Update target item with merged content and votes
	updateQuery := `UPDATE retrospective_items SET content = $1, votes = $2, updated_at = NOW() WHERE id = $3 RETURNING created_at, updated_at`
	err = tx.QueryRow(updateQuery, mergedContent, mergedVotes, targetItemID).Scan(&targetItem.CreatedAt, &targetItem.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Delete source item
	deleteQuery := `DELETE FROM retrospective_items WHERE id = $1`
	_, err = tx.Exec(deleteQuery, sourceItemID)
	if err != nil {
		return nil, err
	}

	// Update the target item content for return
	targetItem.Content = mergedContent
	targetItem.Votes = mergedVotes

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &targetItem, nil
}

// Action Item methods
func (r *RetrospectiveRepository) GetActionItemByID(actionItemID uuid.UUID) (*models.ActionItem, error) {
	query := `
		SELECT id, retrospective_id, item_id, title, description, status, assigned_to, due_date, completed_at, created_by, created_at, updated_at
		FROM action_items
		WHERE id = $1
	`

	var actionItem models.ActionItem
	err := r.db.QueryRow(query, actionItemID).Scan(
		&actionItem.ID,
		&actionItem.RetrospectiveID,
		&actionItem.ItemID,
		&actionItem.Title,
		&actionItem.Description,
		&actionItem.Status,
		&actionItem.AssignedTo,
		&actionItem.DueDate,
		&actionItem.CompletedAt,
		&actionItem.CreatedBy,
		&actionItem.CreatedAt,
		&actionItem.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &actionItem, nil
}

func (r *RetrospectiveRepository) UpdateActionItem(actionItemID uuid.UUID, req *models.ActionItemUpdateRequest) (*models.ActionItem, error) {
	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, *req.Title)
		argIndex++
	}

	if req.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *req.Description)
		argIndex++
	}

	if req.Status != nil {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *req.Status)
		argIndex++
	}

	if req.AssignedTo != nil {
		if *req.AssignedTo == "" {
			setParts = append(setParts, "assigned_to = NULL")
		} else {
			setParts = append(setParts, fmt.Sprintf("assigned_to = $%d", argIndex))
			args = append(args, *req.AssignedTo)
			argIndex++
		}
	}

	if req.DueDate != nil {
		if *req.DueDate == "" {
			setParts = append(setParts, "due_date = NULL")
		} else {
			setParts = append(setParts, fmt.Sprintf("due_date = $%d", argIndex))
			args = append(args, *req.DueDate)
			argIndex++
		}
	}

	if req.CompletedAt != nil {
		if *req.CompletedAt == "" {
			setParts = append(setParts, "completed_at = NULL")
		} else {
			setParts = append(setParts, fmt.Sprintf("completed_at = $%d", argIndex))
			args = append(args, *req.CompletedAt)
			argIndex++
		}
	}

	if len(setParts) == 0 {
		return nil, errors.New("no fields to update")
	}

	setParts = append(setParts, "updated_at = NOW()")
	args = append(args, actionItemID)

	query := fmt.Sprintf(`
		UPDATE action_items 
		SET %s 
		WHERE id = $%d
		RETURNING id, retrospective_id, item_id, title, description, status, assigned_to, due_date, completed_at, created_by, created_at, updated_at
	`, strings.Join(setParts, ", "), argIndex)

	var actionItem models.ActionItem
	err := r.db.QueryRow(query, args...).Scan(
		&actionItem.ID,
		&actionItem.RetrospectiveID,
		&actionItem.ItemID,
		&actionItem.Title,
		&actionItem.Description,
		&actionItem.Status,
		&actionItem.AssignedTo,
		&actionItem.DueDate,
		&actionItem.CompletedAt,
		&actionItem.CreatedBy,
		&actionItem.CreatedAt,
		&actionItem.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &actionItem, nil
}

func (r *RetrospectiveRepository) DeleteActionItem(actionItemID uuid.UUID) error {
	query := `DELETE FROM action_items WHERE id = $1`
	_, err := r.db.Exec(query, actionItemID)
	return err
}

// UpdateTimer updates the timer fields for a retrospective
func (r *RetrospectiveRepository) UpdateTimer(retrospectiveID uuid.UUID, req models.TimerUpdateRequest) error {
	query := "UPDATE retrospectives SET updated_at = NOW()"
	args := []interface{}{}
	argIndex := 1

	if req.Duration != nil {
		query += fmt.Sprintf(", timer_duration = $%d", argIndex)
		args = append(args, *req.Duration)
		argIndex++
	}

	if req.StartedAt != nil {
		query += fmt.Sprintf(", timer_started_at = $%d", argIndex)
		args = append(args, *req.StartedAt)
		argIndex++
	}

	if req.PausedAt != nil {
		query += fmt.Sprintf(", timer_paused_at = $%d", argIndex)
		args = append(args, *req.PausedAt)
		argIndex++
	}

	if req.ElapsedTime != nil {
		query += fmt.Sprintf(", timer_elapsed_time = $%d", argIndex)
		args = append(args, *req.ElapsedTime)
		argIndex++
	}

	query += fmt.Sprintf(" WHERE id = $%d", argIndex)
	args = append(args, retrospectiveID)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update timer: %w", err)
	}

	return nil
}

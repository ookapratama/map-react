package repository

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/ookapratama/go-todo-api/model"
)

// TaskRepository handles database operations for tasks.
type TaskRepository struct {
	db *sql.DB
}

// NewTaskRepository creates a new TaskRepository.
func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// Create inserts a new task into the database.
func (r *TaskRepository) Create(ctx context.Context, req *model.CreateTaskRequest) (*model.Task, error) {
	query := `
		INSERT INTO tasks (title, description, status, due_date)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, description, status, due_date, created_at, updated_at
	`

	task := &model.Task{}
	err := r.db.QueryRowContext(ctx, query,
		req.Title, req.Description, req.Status, req.DueDate,
	).Scan(
		&task.ID, &task.Title, &task.Description,
		&task.Status, &task.DueDate, &task.CreatedAt, &task.UpdatedAt,
	)
	if err != nil {
		log.Error().Err(err).Str("title", req.Title).Msg("Failed to create task")
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	// Normalize due_date to YYYY-MM-DD format
	task.DueDate = task.DueDate[:10]

	log.Info().Str("task_id", task.ID).Msg("Task created successfully")
	return task, nil
}

// GetAll retrieves tasks with filtering, searching, and pagination.
func (r *TaskRepository) GetAll(ctx context.Context, status, search string, page, limit int) (*model.TaskListResponse, error) {
	// Defaults
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Build WHERE conditions
	conditions := []string{}
	args := []interface{}{}
	argIdx := 1

	if status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, status)
		argIdx++
	}

	if search != "" {
		conditions = append(conditions, fmt.Sprintf(
			"(LOWER(title) LIKE $%d OR LOWER(description) LIKE $%d)", argIdx, argIdx+1,
		))
		searchTerm := "%" + strings.ToLower(search) + "%"
		args = append(args, searchTerm, searchTerm)
		argIdx += 2
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total tasks (concurrent with fetching tasks)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tasks %s", whereClause)

	var totalTasks int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalTasks)
	if err != nil {
		log.Error().Err(err).Msg("Failed to count tasks")
		return nil, fmt.Errorf("failed to count tasks: %w", err)
	}

	totalPages := int(math.Ceil(float64(totalTasks) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	// Fetch tasks
	offset := (page - 1) * limit
	dataQuery := fmt.Sprintf(
		"SELECT id, title, description, status, due_date, created_at, updated_at FROM tasks %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		whereClause, argIdx, argIdx+1,
	)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to query tasks")
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	tasks := []model.Task{}
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(
			&task.ID, &task.Title, &task.Description,
			&task.Status, &task.DueDate, &task.CreatedAt, &task.UpdatedAt,
		); err != nil {
			log.Error().Err(err).Msg("Failed to scan task row")
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		task.DueDate = task.DueDate[:10]
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error iterating task rows")
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return &model.TaskListResponse{
		Tasks: tasks,
		Pagination: model.Pagination{
			CurrentPage: page,
			TotalPages:  totalPages,
			TotalTasks:  totalTasks,
		},
	}, nil
}

// GetByID retrieves a single task by its UUID.
func (r *TaskRepository) GetByID(ctx context.Context, id string) (*model.Task, error) {
	query := `
		SELECT id, title, description, status, due_date, created_at, updated_at
		FROM tasks WHERE id = $1
	`

	task := &model.Task{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID, &task.Title, &task.Description,
		&task.Status, &task.DueDate, &task.CreatedAt, &task.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Str("task_id", id).Msg("Task not found")
			return nil, nil
		}
		log.Error().Err(err).Str("task_id", id).Msg("Failed to get task by ID")
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	task.DueDate = task.DueDate[:10]
	return task, nil
}

// Update modifies an existing task.
func (r *TaskRepository) Update(ctx context.Context, id string, req *model.UpdateTaskRequest) (*model.Task, error) {
	query := `
		UPDATE tasks
		SET title = $1, description = $2, status = $3, due_date = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING id, title, description, status, due_date, created_at, updated_at
	`

	task := &model.Task{}
	err := r.db.QueryRowContext(ctx, query,
		req.Title, req.Description, req.Status, req.DueDate, id,
	).Scan(
		&task.ID, &task.Title, &task.Description,
		&task.Status, &task.DueDate, &task.CreatedAt, &task.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Str("task_id", id).Msg("Task not found for update")
			return nil, nil
		}
		log.Error().Err(err).Str("task_id", id).Msg("Failed to update task")
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	task.DueDate = task.DueDate[:10]

	log.Info().Str("task_id", task.ID).Msg("Task updated successfully")
	return task, nil
}

// Delete removes a task by its UUID.
func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tasks WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Error().Err(err).Str("task_id", id).Msg("Failed to delete task")
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get rows affected")
		return fmt.Errorf("failed to check deletion: %w", err)
	}

	if rowsAffected == 0 {
		log.Warn().Str("task_id", id).Msg("Task not found for deletion")
		return sql.ErrNoRows
	}

	log.Info().Str("task_id", id).Msg("Task deleted successfully")
	return nil
}

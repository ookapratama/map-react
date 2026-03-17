package model

import (
	"time"
)

// TaskStatus represents the allowed status values for a task.
type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusCompleted TaskStatus = "completed"
)

// IsValid checks if the status is a valid value.
func (s TaskStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusCompleted:
		return true
	}
	return false
}

// Task represents a to-do task in the system.
type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	DueDate     string     `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CreateTaskRequest is the payload for creating a new task.
type CreateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	DueDate     string     `json:"due_date"`
}

// Validate checks that the request payload is valid.
func (r *CreateTaskRequest) Validate() map[string]string {
	errors := make(map[string]string)

	if r.Title == "" {
		errors["title"] = "Title is required"
	} else if len(r.Title) > 255 {
		errors["title"] = "Title must be at most 255 characters"
	}

	if r.Status == "" {
		errors["status"] = "Status is required"
	} else if !r.Status.IsValid() {
		errors["status"] = "Status must be 'pending' or 'completed'"
	}

	if r.DueDate == "" {
		errors["due_date"] = "Due date is required"
	} else {
		_, err := time.Parse("2006-01-02", r.DueDate)
		if err != nil {
			errors["due_date"] = "Due date must be in YYYY-MM-DD format"
		}
	}

	if len(r.Description) > 1000 {
		errors["description"] = "Description must be at most 1000 characters"
	}

	return errors
}

// UpdateTaskRequest is the payload for updating a task.
type UpdateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	DueDate     string     `json:"due_date"`
}

// Validate checks that the update payload is valid.
func (r *UpdateTaskRequest) Validate() map[string]string {
	errors := make(map[string]string)

	if r.Title == "" {
		errors["title"] = "Title is required"
	} else if len(r.Title) > 255 {
		errors["title"] = "Title must be at most 255 characters"
	}

	if r.Status == "" {
		errors["status"] = "Status is required"
	} else if !r.Status.IsValid() {
		errors["status"] = "Status must be 'pending' or 'completed'"
	}

	if r.DueDate == "" {
		errors["due_date"] = "Due date is required"
	} else {
		_, err := time.Parse("2006-01-02", r.DueDate)
		if err != nil {
			errors["due_date"] = "Due date must be in YYYY-MM-DD format"
		}
	}

	if len(r.Description) > 1000 {
		errors["description"] = "Description must be at most 1000 characters"
	}

	return errors
}

// TaskListResponse is the response format for listing tasks.
type TaskListResponse struct {
	Tasks      []Task     `json:"tasks"`
	Pagination Pagination `json:"pagination"`
}

// Pagination holds pagination metadata.
type Pagination struct {
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	TotalTasks  int `json:"total_tasks"`
}

// ErrorResponse is the standard error response format.
type ErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}

// SuccessResponse is a generic success response.
type SuccessResponse struct {
	Message string `json:"message"`
	Task    *Task  `json:"task,omitempty"`
}

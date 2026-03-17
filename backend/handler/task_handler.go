package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/ookapratama/go-todo-api/model"
	"github.com/ookapratama/go-todo-api/service"
)

// TaskHandler handles HTTP requests for tasks.
type TaskHandler struct {
	service *service.TaskService
}

// NewTaskHandler creates a new TaskHandler.
func NewTaskHandler(svc *service.TaskService) *TaskHandler {
	return &TaskHandler{service: svc}
}

// RegisterRoutes sets up the task-related routes on the given router.
func (h *TaskHandler) RegisterRoutes(r chi.Router) {
	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", h.CreateTask)
		r.Get("/", h.GetAllTasks)
		r.Post("/batch", h.CreateTasksBatch)
		r.Get("/{id}", h.GetTaskByID)
		r.Put("/{id}", h.UpdateTask)
		r.Delete("/{id}", h.DeleteTask)
	})
}

// CreateTask handles POST /tasks
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req model.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode create task request body")
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	task, validationErrors := h.service.CreateTask(r.Context(), &req)
	if validationErrors != nil {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse{
			Error:   "Validation failed",
			Details: validationErrors,
		})
		return
	}

	respondJSON(w, http.StatusCreated, model.SuccessResponse{
		Message: "Task created successfully",
		Task:    task,
	})
}

// GetAllTasks handles GET /tasks
func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	status := query.Get("status")
	search := query.Get("search")

	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit < 1 {
		limit = 10
	}

	// Validate status filter if provided
	if status != "" {
		ts := model.TaskStatus(status)
		if !ts.IsValid() {
			respondJSON(w, http.StatusBadRequest, model.ErrorResponse{
				Error: "Invalid status filter. Must be 'pending' or 'completed'",
			})
			return
		}
	}

	result, err := h.service.GetAllTasks(r.Context(), status, search, page, limit)
	if err != nil {
		log.Error().Err(err).Msg("Handler: failed to get all tasks")
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse{
			Error: "Failed to retrieve tasks",
		})
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// GetTaskByID handles GET /tasks/{id}
func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse{
			Error: "Task ID is required",
		})
		return
	}

	task, err := h.service.GetTaskByID(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("task_id", id).Msg("Handler: failed to get task")
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse{
			Error: "Failed to retrieve task",
		})
		return
	}

	if task == nil {
		respondJSON(w, http.StatusNotFound, model.ErrorResponse{
			Error: "Task not found",
		})
		return
	}

	respondJSON(w, http.StatusOK, task)
}

// UpdateTask handles PUT /tasks/{id}
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse{
			Error: "Task ID is required",
		})
		return
	}

	var req model.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode update task request body")
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	task, validationErrors := h.service.UpdateTask(r.Context(), id, &req)
	if validationErrors != nil {
		// Check if it's a server error
		if _, ok := validationErrors["server"]; ok {
			respondJSON(w, http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to update task",
			})
			return
		}
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse{
			Error:   "Validation failed",
			Details: validationErrors,
		})
		return
	}

	if task == nil {
		respondJSON(w, http.StatusNotFound, model.ErrorResponse{
			Error: "Task not found",
		})
		return
	}

	respondJSON(w, http.StatusOK, model.SuccessResponse{
		Message: "Task updated successfully",
		Task:    task,
	})
}

// DeleteTask handles DELETE /tasks/{id}
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse{
			Error: "Task ID is required",
		})
		return
	}

	err := h.service.DeleteTask(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			respondJSON(w, http.StatusNotFound, model.ErrorResponse{
				Error: "Task not found",
			})
			return
		}
		log.Error().Err(err).Str("task_id", id).Msg("Handler: failed to delete task")
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse{
			Error: "Failed to delete task",
		})
		return
	}

	respondJSON(w, http.StatusOK, model.SuccessResponse{
		Message: "Task deleted successfully",
	})
}

// CreateTasksBatch handles POST /tasks/batch - creates multiple tasks concurrently.
func (h *TaskHandler) CreateTasksBatch(w http.ResponseWriter, r *http.Request) {
	var requests []*model.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&requests); err != nil {
		log.Error().Err(err).Msg("Failed to decode batch create request body")
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse{
			Error: "Invalid request body. Expected an array of tasks.",
		})
		return
	}

	if len(requests) == 0 {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse{
			Error: "At least one task is required",
		})
		return
	}

	if len(requests) > 50 {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse{
			Error: "Maximum 50 tasks can be created at once",
		})
		return
	}

	tasks, errs := h.service.CreateTasksConcurrently(r.Context(), requests)

	response := map[string]interface{}{
		"message":       "Batch creation completed",
		"created_count": len(tasks),
		"error_count":   len(errs),
		"tasks":         tasks,
	}

	if len(errs) > 0 {
		errMessages := make([]string, len(errs))
		for i, e := range errs {
			errMessages[i] = e.Error()
		}
		response["errors"] = errMessages
	}

	statusCode := http.StatusCreated
	if len(errs) > 0 && len(tasks) == 0 {
		statusCode = http.StatusBadRequest
	} else if len(errs) > 0 {
		statusCode = http.StatusMultiStatus
	}

	respondJSON(w, statusCode, response)
}

// respondJSON writes a JSON response with the given status code.
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Error().Err(err).Msg("Failed to encode JSON response")
	}
}

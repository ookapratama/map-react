package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/ookapratama/go-todo-api/model"
	"github.com/ookapratama/go-todo-api/repository"
)

// TaskService handles business logic for tasks.
type TaskService struct {
	repo *repository.TaskRepository
}

// NewTaskService creates a new TaskService.
func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

// CreateTask validates and creates a new task.
func (s *TaskService) CreateTask(ctx context.Context, req *model.CreateTaskRequest) (*model.Task, map[string]string) {
	if errs := req.Validate(); len(errs) > 0 {
		return nil, errs
	}

	task, err := s.repo.Create(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("Service: failed to create task")
		return nil, map[string]string{"server": "Internal server error"}
	}

	return task, nil
}

// GetAllTasks retrieves tasks with concurrent count + data fetch.
// This demonstrates concurrent execution using goroutines and WaitGroup.
func (s *TaskService) GetAllTasks(ctx context.Context, status, search string, page, limit int) (*model.TaskListResponse, error) {
	// Use a channel-based concurrent approach for potential enhancements:
	// We fetch the data concurrently and also perform async logging of the access.
	type result struct {
		response *model.TaskListResponse
		err      error
	}

	ch := make(chan result, 1)

	// WaitGroup to coordinate goroutines
	var wg sync.WaitGroup
	wg.Add(1)

	// Goroutine 1: Fetch task data
	go func() {
		defer wg.Done()
		resp, err := s.repo.GetAll(ctx, status, search, page, limit)
		ch <- result{response: resp, err: err}
	}()

	// Goroutine 2: Concurrent async access logging
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info().
			Str("status_filter", status).
			Str("search_term", search).
			Int("page", page).
			Int("limit", limit).
			Msg("Task list accessed concurrently")
	}()

	// Wait for the data fetch result
	res := <-ch

	// Wait for all goroutines to finish
	wg.Wait()

	if res.err != nil {
		log.Error().Err(res.err).Msg("Service: failed to get all tasks")
		return nil, res.err
	}

	return res.response, nil
}

// GetTaskByID retrieves a single task.
func (s *TaskService) GetTaskByID(ctx context.Context, id string) (*model.Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("task_id", id).Msg("Service: failed to get task")
		return nil, err
	}
	return task, nil
}

// UpdateTask validates and updates a task.
func (s *TaskService) UpdateTask(ctx context.Context, id string, req *model.UpdateTaskRequest) (*model.Task, map[string]string) {
	if errs := req.Validate(); len(errs) > 0 {
		return nil, errs
	}

	task, err := s.repo.Update(ctx, id, req)
	if err != nil {
		log.Error().Err(err).Str("task_id", id).Msg("Service: failed to update task")
		return nil, map[string]string{"server": "Internal server error"}
	}

	return task, nil
}

// DeleteTask removes a task.
func (s *TaskService) DeleteTask(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// CreateTasksConcurrently creates multiple tasks concurrently.
// This demonstrates the use of goroutines and WaitGroup for concurrent execution.
func (s *TaskService) CreateTasksConcurrently(ctx context.Context, requests []*model.CreateTaskRequest) ([]*model.Task, []error) {
	var (
		mu     sync.Mutex
		wg     sync.WaitGroup
		tasks  []*model.Task
		errors []error
	)

	for _, req := range requests {
		wg.Add(1)
		go func(r *model.CreateTaskRequest) {
			defer wg.Done()

			if errs := r.Validate(); len(errs) > 0 {
				mu.Lock()
				errors = append(errors, fmt.Errorf("validation failed for task '%s': %v", r.Title, errs))
				mu.Unlock()
				return
			}

			task, err := s.repo.Create(ctx, r)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errors = append(errors, err)
				log.Error().Err(err).Str("title", r.Title).Msg("Concurrent create failed")
			} else {
				tasks = append(tasks, task)
				log.Info().Str("task_id", task.ID).Msg("Task created concurrently")
			}
		}(req)
	}

	wg.Wait()
	return tasks, errors
}

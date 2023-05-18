package service

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/repository"
)

type TaskService interface {
	FindAll(int, int, string) (*[]db.Task, error)
	FindById(int) (*db.Task, error)
	Create(*models.CreateTask) (*db.Task, error)
	Update(int, *models.UpdateTask) (*db.Task, error)
	Delete(int) error
}

type taskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo}
}

// Creates a task in the database
func (s *taskService) Create(task *models.CreateTask) (*db.Task, error) {
	// Create a new struct of type task
	taskToCreate := db.Task{
		TaskName: task.TaskName,
		Status:   task.Status,
		Type:     task.Type,
		Notes:    task.Notes,
	}

	// Create task in database
	createdTask, err := s.repo.Create(&taskToCreate)
	if err != nil {
		return nil, fmt.Errorf("failed creating task: %w", err)
	}

	return createdTask, nil
}

// Find a list of tasks in the database
func (s *taskService) FindAll(limit int, offset int, order string) (*[]db.Task, error) {
	tasks, err := s.repo.FindAll(limit, offset, order)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// Find task in database by ID
func (s *taskService) FindById(id int) (*db.Task, error) {
	// Find task by id
	task, err := s.repo.FindById(id)

	// If error detected
	if err != nil {
		return nil, err
	}
	// else
	return task, nil
}

// Delete task in database
func (s *taskService) Delete(id int) error {
	err := s.repo.Delete(id)
	// If error detected
	if err != nil {
		fmt.Println("error in deleting task: ", err)
		return err
	}
	// else
	return nil
}

// Updates task in database
func (s *taskService) Update(id int, task *models.UpdateTask) (*db.Task, error) {
	// Create db property type of incoming DTO
	taskToCreate := db.Task{
		TaskName:   task.TaskName,
		Assignment: task.Assignment,
		Status:     task.Status,
		Type:       task.Type,
		Notes:      task.Notes,
	}

	// Update using repo
	updatedTask, err := s.repo.Update(id, &taskToCreate)
	if err != nil {
		return nil, err
	}

	return updatedTask, nil
}

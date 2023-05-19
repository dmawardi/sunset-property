package service

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/repository"
)

type TaskLogService interface {
	FindAll(int, int, string) (*[]db.TaskLog, error)
	FindById(int) (*db.TaskLog, error)
	Create(*models.CreateTaskLog) (*db.TaskLog, error)
	Update(int, *models.UpdateTaskLog) (*db.TaskLog, error)
	Delete(int) error
}

type taskLogService struct {
	repo repository.TaskLogRepository
}

func NewTaskLogService(repo repository.TaskLogRepository) TaskLogService {
	return &taskLogService{repo}
}

// Creates a task log message
func (s *taskLogService) Create(log *models.CreateTaskLog) (*db.TaskLog, error) {
	// Create a new property of type db User
	logMessageToCreate := db.TaskLog{
		User:       log.User,
		TaskID:     log.Task.ID,
		LogMessage: log.LogMessage,
		Type:       log.Type,
	}

	// Create task in database
	createdMsg, err := s.repo.Create(&logMessageToCreate)
	if err != nil {
		return nil, fmt.Errorf("failed creating task log message: %w", err)
	}

	return createdMsg, nil
}

// Find a list of task log messages
func (s *taskLogService) FindAll(limit int, offset int, order string) (*[]db.TaskLog, error) {
	logMessages, err := s.repo.FindAll(limit, offset, order)
	if err != nil {
		return nil, err
	}
	return logMessages, nil
}

// Find task log message in database by ID
func (s *taskLogService) FindById(id int) (*db.TaskLog, error) {
	// Find log message by id
	logMessage, err := s.repo.FindById(id)
	// If error detected
	if err != nil {
		return nil, err
	}
	// else
	return logMessage, nil
}

// Delete task log message in database
func (s *taskLogService) Delete(id int) error {
	err := s.repo.Delete(id)
	// If error detected
	if err != nil {
		fmt.Println("error in deleting task log message: ", err)
		return err
	}
	// else
	return nil
}

// Updates task log message in database (Only log message can be updated)
func (s *taskLogService) Update(id int, log *models.UpdateTaskLog) (*db.TaskLog, error) {
	// Create db task Log message type from DTO
	logMessage := db.TaskLog{
		LogMessage: log.LogMessage,
	}

	// Update using repo
	updatedLogMessage, err := s.repo.Update(id, &logMessage)
	if err != nil {
		return nil, err
	}

	return updatedLogMessage, nil
}

// Generates a log message based on task update DTO
// func generateLogTaskUpdate(taskUpdate *models.UpdateTask) string {
// 	// Iterate through all fields in property update
// 	var logMessage string

// 	values := reflect.ValueOf(taskUpdate)
// 	types := values.Type()
// 	for i := 0; i < values.NumField(); i++ {
// 		fmt.Println(types.Field(i).Index[0], types.Field(i).Name, values.Field(i))
// 	}
// 	return logMessage
// }

package service

import (
	"fmt"
	"reflect"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/repository"
)

type PropertyLogService interface {
	FindAll(int, int, string) (*[]db.PropertyLog, error)
	FindById(int) (*db.PropertyLog, error)
	Create(*models.CreatePropertyLog) (*db.PropertyLog, error)
	Update(int, *models.UpdatePropertyLog) (*db.PropertyLog, error)
	Delete(int) error
}

type propertyLogService struct {
	repo repository.PropertyLogRepository
}

func NewPropertyLogService(repo repository.PropertyLogRepository) PropertyLogService {
	return &propertyLogService{repo}
}

// Creates a property log message in the database
func (s *propertyLogService) Create(log *models.CreatePropertyLog) (*db.PropertyLog, error) {
	// Create a new property of type db User
	logMessageToCreate := db.PropertyLog{
		User:       log.User,
		Property:   log.Property,
		LogMessage: log.LogMessage,
	}

	// Create above user in database
	createdMsg, err := s.repo.Create(&logMessageToCreate)
	if err != nil {
		return nil, fmt.Errorf("failed creating property log message: %w", err)
	}

	return createdMsg, nil
}

// Find a list of property log messages in the database
func (s *propertyLogService) FindAll(limit int, offset int, order string) (*[]db.PropertyLog, error) {
	logMessages, err := s.repo.FindAll(limit, offset, order)
	if err != nil {
		return nil, err
	}
	return logMessages, nil
}

// Find property log message in database by ID
func (s *propertyLogService) FindById(id int) (*db.PropertyLog, error) {
	// Find log message by id
	logMessage, err := s.repo.FindById(id)
	// If error detected
	if err != nil {
		return nil, err
	}
	// else
	return logMessage, nil
}

// Delete property log message in database
func (s *propertyLogService) Delete(id int) error {
	err := s.repo.Delete(id)
	// If error detected
	if err != nil {
		fmt.Println("error in deleting property feature: ", err)
		return err
	}
	// else
	return nil
}

// Updates property log message in database (Only log message can be updated)
func (s *propertyLogService) Update(id int, log *models.UpdatePropertyLog) (*db.PropertyLog, error) {
	// Create db Property Log message type from DTO
	logMessage := db.PropertyLog{
		LogMessage: log.LogMessage,
	}

	// Update using repo
	updatedLogMessage, err := s.repo.Update(id, &logMessage)
	if err != nil {
		return nil, err
	}

	return updatedLogMessage, nil
}

// Generates a log message based on property update DTO
func generateLogPropertyUpdate(propUpdate *models.UpdateProperty) string {
	// Iterate through all fields in property update
	var logMessage string

	values := reflect.ValueOf(propUpdate)
	types := values.Type()
	for i := 0; i < values.NumField(); i++ {
		fmt.Println(types.Field(i).Index[0], types.Field(i).Name, values.Field(i))
	}
	return logMessage
}

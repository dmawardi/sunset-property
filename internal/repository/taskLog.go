package repository

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"gorm.io/gorm"
)

type TaskLogRepository interface {
	FindAll(int, int, string) (*[]db.TaskLog, error)
	FindById(int) (*db.TaskLog, error)
	Create(*db.TaskLog) (*db.TaskLog, error)
	Update(int, *db.TaskLog) (*db.TaskLog, error)
	Delete(int) error
}

type taskLogRepository struct {
	DB *gorm.DB
}

func NewTaskLogRepository(db *gorm.DB) TaskLogRepository {
	return &taskLogRepository{db}
}

// Creates a task log message in the database
func (r *taskLogRepository) Create(logMessage *db.TaskLog) (*db.TaskLog, error) {
	// Create new log message in database
	result := r.DB.Create(&logMessage)
	if result.Error != nil {
		return nil, fmt.Errorf("failed creating task log message: %w", result.Error)
	}

	return logMessage, nil
}

// Find a list of log messages in the database
func (r *taskLogRepository) FindAll(limit int, offset int, order string) (*[]db.TaskLog, error) {
	// Query all log messages based on the received parameters
	logMessages, err := QueryAllTaskLogsBasedOnParams(limit, offset, order, r.DB)
	if err != nil {
		fmt.Printf("Error querying db for list of task log Messages: %s", err)
		return nil, err
	}

	return &logMessages, nil
}

// Find a task log message in database by ID
func (r *taskLogRepository) FindById(id int) (*db.TaskLog, error) {
	// Create an empty ref object of type property log
	logMessage := db.TaskLog{}
	// Grab log message from db if exists
	result := r.DB.Preload("User").Preload("Task").First(&logMessage, id)

	// If error detected
	if result.Error != nil {
		return nil, result.Error
	}
	// else
	return &logMessage, nil
}

// Delete task log message in database
func (r *taskLogRepository) Delete(id int) error {
	// Create an empty ref object of type property log message
	logMessage := db.TaskLog{}
	// Delete log message from db if exists
	result := r.DB.Delete(&logMessage, id)

	// If error detected
	if result.Error != nil {
		fmt.Println("error in deleting task log message: ", result.Error)
		return result.Error
	}
	// else
	return nil
}

// Updates task log message in database
func (r *taskLogRepository) Update(id int, feature *db.TaskLog) (*db.TaskLog, error) {
	// Init
	var err error
	// Find task log message by id to ensure it exists
	foundLogMessage, err := r.FindById(id)
	if err != nil {
		fmt.Println("Task log message to update not found: ", err)
		return nil, err
	}

	// Update found log message
	updateResult := r.DB.Model(&foundLogMessage).Updates(feature)
	if updateResult.Error != nil {
		fmt.Println("Task log message update failed: ", updateResult.Error)
		return nil, updateResult.Error
	}

	// Retrieve updated task log message by id
	updatedLogMessage, err := r.FindById(id)
	if err != nil {
		fmt.Println("Updated Task log message not found: ", err)
		return nil, err
	}
	return updatedLogMessage, nil
}

// Takes limit, offset, and order parameters, builds a query and executes returning a list of log messages
func QueryAllTaskLogsBasedOnParams(limit int, offset int, order string, dbClient *gorm.DB) ([]db.TaskLog, error) {
	// Build model to query database
	log := []db.TaskLog{}
	// Build base query for property log messages table
	query := dbClient.Model(&log)

	// Add parameters into query as needed
	if limit != 0 {
		query.Limit(limit)
	}
	if offset != 0 {
		query.Offset(offset)
	}
	// order format should be "column_name ASC/DESC" eg. "created_at ASC"
	if order != "" {
		query.Order(order)
	} else {
		query.Order("created_at DESC")
	}
	// Query database
	result := query.Find(&log)
	if result.Error != nil {
		return nil, result.Error
	}
	// Return if no errors with result
	return log, nil
}

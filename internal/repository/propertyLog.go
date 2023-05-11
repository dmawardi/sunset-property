package repository

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"gorm.io/gorm"
)

type PropertyLogRepository interface {
	FindAll(int, int, string) (*[]db.PropertyLog, error)
	FindById(int) (*db.PropertyLog, error)
	Create(*db.PropertyLog) (*db.PropertyLog, error)
	Update(int, *db.PropertyLog) (*db.PropertyLog, error)
	Delete(int) error
}

type propertyLogRepository struct {
	DB *gorm.DB
}

func NewPropertyLogRepository(db *gorm.DB) PropertyLogRepository {
	return &propertyLogRepository{db}
}

// Creates a Property log message in the database
func (r *propertyLogRepository) Create(logMessage *db.PropertyLog) (*db.PropertyLog, error) {
	// Create new log message in database
	result := r.DB.Create(&logMessage)
	if result.Error != nil {
		return nil, fmt.Errorf("failed creating property log message: %w", result.Error)
	}

	return logMessage, nil
}

// Find a list of log messages in the database
func (r *propertyLogRepository) FindAll(limit int, offset int, order string) (*[]db.PropertyLog, error) {
	// Query all log messages based on the received parameters
	logMessages, err := QueryAllPropertyLogsBasedOnParams(limit, offset, order, r.DB)
	if err != nil {
		fmt.Printf("Error querying db for list of logMessages: %s", err)
		return nil, err
	}

	return &logMessages, nil
}

// Find a property log message in database by ID
func (r *propertyLogRepository) FindById(id int) (*db.PropertyLog, error) {
	// Create an empty ref object of type property log
	logMessage := db.PropertyLog{}
	// Grab log message from db if exists
	result := r.DB.First(&logMessage, id)

	// If error detected
	if result.Error != nil {
		return nil, result.Error
	}
	// else
	return &logMessage, nil
}

// Delete property log message in database
func (r *propertyLogRepository) Delete(id int) error {
	// Create an empty ref object of type property log message
	logMessage := db.PropertyLog{}
	// Delete log message from db if exists
	result := r.DB.Delete(&logMessage, id)

	// If error detected
	if result.Error != nil {
		fmt.Println("error in deleting property log message: ", result.Error)
		return result.Error
	}
	// else
	return nil
}

// Updates property log message in database
func (r *propertyLogRepository) Update(id int, feature *db.PropertyLog) (*db.PropertyLog, error) {
	// Init
	var err error
	// Find property log message by id to ensure it exists
	foundLogMessage, err := r.FindById(id)
	if err != nil {
		fmt.Println("Property log message to update not found: ", err)
		return nil, err
	}

	// Update found log message
	updateResult := r.DB.Model(&foundLogMessage).Updates(feature)
	if updateResult.Error != nil {
		fmt.Println("Property log message update failed: ", updateResult.Error)
		return nil, updateResult.Error
	}

	// Retrieve updated property log message by id
	updatedLogMessage, err := r.FindById(id)
	if err != nil {
		fmt.Println("Updated property log message not found: ", err)
		return nil, err
	}
	return updatedLogMessage, nil
}

// Takes limit, offset, and order parameters, builds a query and executes returning a list of log messages
func QueryAllPropertyLogsBasedOnParams(limit int, offset int, order string, dbClient *gorm.DB) ([]db.PropertyLog, error) {
	// Build model to query database
	log := []db.PropertyLog{}
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
	}
	// Query database
	result := query.Find(&log)
	if result.Error != nil {
		return nil, result.Error
	}
	// Return if no errors with result
	return log, nil
}

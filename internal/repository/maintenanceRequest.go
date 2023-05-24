package repository

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"gorm.io/gorm"
)

type MaintenanceRequestRepository interface {
	FindAll(int, int, string) (*[]db.MaintenanceRequest, error)
	FindById(int) (*db.MaintenanceRequest, error)
	Create(*db.MaintenanceRequest) (*db.MaintenanceRequest, error)
	Update(int, *db.MaintenanceRequest) (*db.MaintenanceRequest, error)
	Delete(int) error
}

type maintenanceRequestRepository struct {
	DB *gorm.DB
}

func NewMaintenanceRequestRepository(db *gorm.DB) MaintenanceRequestRepository {
	return &maintenanceRequestRepository{db}
}

// Creates a maintenance request in the database
func (r *maintenanceRequestRepository) Create(request *db.MaintenanceRequest) (*db.MaintenanceRequest, error) {
	// Create new log message in database
	result := r.DB.Create(&request)
	if result.Error != nil {
		return nil, fmt.Errorf("failed creating maintenance request: %w", result.Error)
	}

	return request, nil
}

// Find a list of maintenance requests in the database
func (r *maintenanceRequestRepository) FindAll(limit int, offset int, order string) (*[]db.MaintenanceRequest, error) {
	// Query all maintenance requests based on the received parameters
	requests, err := QueryAllMaintenanceRequestsBasedOnParams(limit, offset, order, r.DB)
	if err != nil {
		fmt.Printf("Error querying db for list of maintenance requests: %s", err)
		return nil, err
	}

	return &requests, nil
}

// Find a maintenance request in database by ID
func (r *maintenanceRequestRepository) FindById(id int) (*db.MaintenanceRequest, error) {
	// Create an empty ref object of type maintenance request
	request := db.MaintenanceRequest{}
	// Grab maint. request from db if exists
	result := r.DB.Preload("Property").First(&request, id)

	// If error detected
	if result.Error != nil {
		return nil, result.Error
	}
	// else
	return &request, nil
}

// Delete maintenance request in database
func (r *maintenanceRequestRepository) Delete(id int) error {
	// Create an empty ref object of type maintenance request
	request := db.MaintenanceRequest{}
	// Delete maint. request from db if exists
	result := r.DB.Delete(&request, id)

	// If error detected
	if result.Error != nil {
		fmt.Println("error in deleting maintenance request: ", result.Error)
		return result.Error
	}
	// else
	return nil
}

// Updates maintenance request in database
func (r *maintenanceRequestRepository) Update(id int, request *db.MaintenanceRequest) (*db.MaintenanceRequest, error) {
	// Init
	var err error
	// Find maint. request by id to ensure it exists
	foundRequest, err := r.FindById(id)
	if err != nil {
		fmt.Println("Maintenance request to update not found: ", err)
		return nil, err
	}

	// Update found maint. request with incoming details
	updateResult := r.DB.Model(&foundRequest).Updates(request)
	if updateResult.Error != nil {
		fmt.Println("Maintenance request update failed: ", updateResult.Error)
		return nil, updateResult.Error
	}

	// Build associations
	assResult := r.DB.Model(&foundRequest).Association("Property").Append(request.Property)
	// Check if association update failed
	if assResult != nil {
		fmt.Println("Property association update failed: ", assResult)
		return nil, assResult
	}

	// Retrieve updated maintenance request by id
	updatedRequest, err := r.FindById(id)
	if err != nil {
		fmt.Println("Updated maintenance request not found: ", err)
		return nil, err
	}
	return updatedRequest, nil
}

// Takes limit, offset, and order parameters, builds a query and executes returning a list of maintenance requests
func QueryAllMaintenanceRequestsBasedOnParams(limit int, offset int, order string, dbClient *gorm.DB) ([]db.MaintenanceRequest, error) {
	// Build model to query database
	maintenanceRequest := []db.MaintenanceRequest{}
	// Build base query for property log messages table
	query := dbClient.Model(&maintenanceRequest)

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
	result := query.Find(&maintenanceRequest)
	if result.Error != nil {
		return nil, result.Error
	}
	// Return if no errors with result
	return maintenanceRequest, nil
}

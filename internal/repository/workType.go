package repository

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"gorm.io/gorm"
)

type WorkTypeRepository interface {
	FindAll(int, int, string) (*[]db.WorkType, error)
	FindById(int) (*db.WorkType, error)
	Create(*db.WorkType) (*db.WorkType, error)
	Update(int, *db.WorkType) (*db.WorkType, error)
	Delete(int) error
}

type workTypeRepository struct {
	DB *gorm.DB
}

func NewWorkTypeRepository(db *gorm.DB) WorkTypeRepository {
	return &workTypeRepository{db}
}

// Creates a work type in the database
func (r *workTypeRepository) Create(workType *db.WorkType) (*db.WorkType, error) {
	// Create in database
	result := r.DB.Create(&workType)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create work type: %w", result.Error)
	}

	return workType, nil
}

// Find a list of work types in the database
func (r *workTypeRepository) FindAll(limit int, offset int, order string) (*[]db.WorkType, error) {
	// Query all work types based on the received parameters
	workTypes, err := QueryAllWorkTypesBasedOnParams(limit, offset, order, r.DB)
	if err != nil {
		fmt.Printf("Error querying db for list of work types: %s", err)
		return nil, err
	}

	return &workTypes, nil
}

// Find a work type in database by ID
func (r *workTypeRepository) FindById(id int) (*db.WorkType, error) {
	// Create an empty ref object of type work type
	workType := db.WorkType{}
	// Grab work type from db if exists
	result := r.DB.First(&workType, id)

	// If error detected
	if result.Error != nil {
		return nil, result.Error
	}
	// else
	return &workType, nil
}

// Delete work type in database
func (r *workTypeRepository) Delete(id int) error {
	// Create an empty ref object of type work type
	request := db.WorkType{}
	// Delete work type from db if exists
	result := r.DB.Delete(&request, id)

	// If error detected
	if result.Error != nil {
		fmt.Println("error in deleting maintenance request: ", result.Error)
		return result.Error
	}
	// else
	return nil
}

// Updates work type in database
func (r *workTypeRepository) Update(id int, workType *db.WorkType) (*db.WorkType, error) {
	// Init
	var err error
	// Find work type by id to ensure it exists
	foundWorkType, err := r.FindById(id)
	if err != nil {
		fmt.Println("Work type to update not found: ", err)
		return nil, err
	}

	// Update found work type with incoming details
	updateResult := r.DB.Model(&foundWorkType).Updates(workType)
	if updateResult.Error != nil {
		fmt.Println("Work type update failed: ", updateResult.Error)
		return nil, updateResult.Error
	}

	// Retrieve updated work type by id
	updatedWorkType, err := r.FindById(id)
	if err != nil {
		fmt.Println("Updated work type not found: ", err)
		return nil, err
	}
	return updatedWorkType, nil
}

// Takes limit, offset, and order parameters, builds a query and executes returning a list of work types
func QueryAllWorkTypesBasedOnParams(limit int, offset int, order string, dbClient *gorm.DB) ([]db.WorkType, error) {
	// Build model to query database
	workType := []db.WorkType{}
	// Build base query for property log messages table
	query := dbClient.Model(&workType)

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
	result := query.Find(&workType)
	if result.Error != nil {
		return nil, result.Error
	}
	// Return if no errors with result
	return workType, nil
}

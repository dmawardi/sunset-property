package repository

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"gorm.io/gorm"
)

type PropertyRepository interface {
	FindAll(int, int, string) (*[]db.Property, error)
	FindById(int) (*db.Property, error)
	Create(property *db.Property) (*db.Property, error)
	Update(int, *db.Property) (*db.Property, error)
	Delete(int) error
}

type propertyRepository struct {
	DB *gorm.DB
}

func NewPropertyRepository(db *gorm.DB) PropertyRepository {
	return &propertyRepository{db}
}

// Creates a property in the database
func (r *propertyRepository) Create(property *db.Property) (*db.Property, error) {
	// Create above property in database
	result := r.DB.Create(&property)
	if result.Error != nil {
		return nil, fmt.Errorf("failed creating property: %w", result.Error)
	}

	// Build associations
	assResult := r.DB.Model(&property).Association("Features").Append(property.Features)
	// Check if association update failed
	if assResult != nil {
		fmt.Println("Property association update failed: ", assResult)
		return nil, assResult
	}

	return property, nil
}

// Find a list of properties in the database
func (r *propertyRepository) FindAll(limit int, offset int, order string) (*[]db.Property, error) {
	// Query all properties based on the received parameters
	properties, err := QueryAllPropertiesBasedOnParams(limit, offset, order, r.DB)
	if err != nil {
		fmt.Printf("Error querying db for list of properties: %s", err)
		return nil, err
	}

	return &properties, nil
}

// Find property in database by ID
func (r *propertyRepository) FindById(id int) (*db.Property, error) {
	// Create an empty ref object of type property
	property := db.Property{}
	// Check if property exists in db
	result := r.DB.Preload("Features").Preload("PropertyLogs").Preload("Contacts").First(&property, id)

	// Extract error result
	err := result.Error
	fmt.Printf("err found searching for id (%v): %v\n", id, err)
	// If error detected
	if err != nil {
		return nil, err
	}
	// else
	return &property, nil
}

// Delete property in database
func (r *propertyRepository) Delete(id int) error {
	// Create an empty ref object of type property
	property := db.Property{}
	// Check if property exists in db
	result := r.DB.Delete(&property, id)

	// If error detected
	if result.Error != nil {
		fmt.Println("error in deleting property: ", result.Error)
		return result.Error
	}
	// else
	return nil
}

// Updates property in database
func (r *propertyRepository) Update(id int, property *db.Property) (*db.Property, error) {
	// Init
	var err error
	// Find property by id
	foundProperty, err := r.FindById(id)
	if err != nil {
		fmt.Println("Property to update not found: ", err)
		return nil, err
	}

	// Update found property using new property
	updateResult := r.DB.Model(&foundProperty).Updates(property)

	// Extract error
	err = updateResult.Error
	// Check if update failed
	if err != nil {
		fmt.Println("Property update failed: ", err)
		return nil, err
	}

	// Init associate error
	var assResult error
	// Depending on if features already exist on property
	if len(foundProperty.Features) > 0 {
		// Replace if existent
		assResult = r.DB.Model(&foundProperty).Association("Features").Replace(property.Features)
	} else {
		// Append if non existent
		assResult = r.DB.Model(&foundProperty).Association("Features").Append(property.Features)
	}
	// Check if association update failed
	if assResult != nil {
		fmt.Println("Property association update failed: ", assResult)
		return nil, assResult
	}

	// Depending on if contacts already exist on property
	if len(property.Contacts) > 0 {
		// Replace
		assResult = r.DB.Model(&foundProperty).Association("Contacts").Replace(property.Contacts)
	}
	// Check if association update failed
	if assResult != nil {
		fmt.Println("Property association update failed: ", assResult)
		return nil, assResult
	}

	// Retrieve updated property by id
	updatedProperty, err := r.FindById(id)
	if err != nil {
		fmt.Println("Updated property not found: ", err)
		return nil, assResult
	}
	return updatedProperty, nil
}

// Takes limit, offset, and order parameters, builds a query and executes returning a list of properties
func QueryAllPropertiesBasedOnParams(limit int, offset int, order string, dbClient *gorm.DB) ([]db.Property, error) {
	// Build model to query database
	properties := []db.Property{}
	// Build base query for properties table
	query := dbClient.Model(&properties).Preload("Features")

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
	result := query.Find(&properties)
	if result.Error != nil {
		return nil, result.Error
	}
	// Return if no errors with result
	return properties, nil
}

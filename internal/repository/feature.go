package repository

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"gorm.io/gorm"
)

type FeatureRepository interface {
	FindAll(int, int, string) (*[]db.Feature, error)
	FindById(int) (*db.Feature, error)
	Create(feature *db.Feature) (*db.Feature, error)
	Update(int, *db.Feature) (*db.Feature, error)
	Delete(int) error
}

type featureRepository struct {
	DB *gorm.DB
}

func NewFeatureRepository(db *gorm.DB) FeatureRepository {
	return &featureRepository{db}
}

// Creates a property feature in the database
func (r *featureRepository) Create(feature *db.Feature) (*db.Feature, error) {
	// Create above property in database
	result := r.DB.Create(&feature)
	if result.Error != nil {
		return nil, fmt.Errorf("failed creating property feature: %w", result.Error)
	}

	return feature, nil
}

// Find a list of property features in the database
func (r *featureRepository) FindAll(limit int, offset int, order string) (*[]db.Feature, error) {
	// Query all property features based on the received parameters
	features, err := QueryAllFeaturesBasedOnParams(limit, offset, order, r.DB)
	if err != nil {
		fmt.Printf("Error querying db for list of features: %s", err)
		return nil, err
	}

	return &features, nil
}

// Find property feature in database by ID
func (r *featureRepository) FindById(id int) (*db.Feature, error) {
	// Create an empty ref object of type feature
	feature := db.Feature{}
	// Check if property exists in db
	result := r.DB.First(&feature, id)

	// If error detected
	if result.Error != nil {
		return nil, result.Error
	}
	// else
	return &feature, nil
}

// Delete property feature in database
func (r *featureRepository) Delete(id int) error {
	// Create an empty ref object of type property feature
	feature := db.Feature{}
	// Check if property exists in db
	result := r.DB.Delete(&feature, id)

	// If error detected
	if result.Error != nil {
		fmt.Println("error in deleting property feature: ", result.Error)
		return result.Error
	}
	// else
	return nil
}

// Updates property feature in database
func (r *featureRepository) Update(id int, feature *db.Feature) (*db.Feature, error) {
	// Init
	var err error
	// Find property feature by id
	foundFeature, err := r.FindById(id)
	if err != nil {
		fmt.Println("Property feature to update not found: ", err)
		return nil, err
	}

	// Update found feature using new feature
	updateResult := r.DB.Model(&foundFeature).Updates(feature)
	if updateResult.Error != nil {
		fmt.Println("Property feature update failed: ", updateResult.Error)
		return nil, updateResult.Error
	}

	// Retrieve updated property feature by id
	updatedFeature, err := r.FindById(id)
	if err != nil {
		fmt.Println("Updated property feature not found: ", err)
		return nil, err
	}
	return updatedFeature, nil
}

// Takes limit, offset, and order parameters, builds a query and executes returning a list of properties
func QueryAllFeaturesBasedOnParams(limit int, offset int, order string, dbClient *gorm.DB) ([]db.Feature, error) {
	// Build model to query database
	features := []db.Feature{}
	// Build base query for properties table
	query := dbClient.Model(&features)

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
	result := query.Find(&features)
	if result.Error != nil {
		return nil, result.Error
	}
	// Return if no errors with result
	return features, nil
}

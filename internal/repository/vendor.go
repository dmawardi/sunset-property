package repository

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"gorm.io/gorm"
)

type VendorRepository interface {
	FindAll(int, int, string) (*[]db.Vendor, error)
	FindById(int) (*db.Vendor, error)
	Create(*db.Vendor) (*db.Vendor, error)
	Update(int, *db.Vendor) (*db.Vendor, error)
	Delete(int) error
}

type vendorRepository struct {
	DB *gorm.DB
}

func NewVendorRepository(db *gorm.DB) VendorRepository {
	return &vendorRepository{db}
}

// Creates a vendor in the database
func (r *vendorRepository) Create(vendor *db.Vendor) (*db.Vendor, error) {
	// Create in database
	result := r.DB.Create(&vendor)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create vendor: %w", result.Error)
	}

	return vendor, nil
}

// Find a list of vendors in the database
func (r *vendorRepository) FindAll(limit int, offset int, order string) (*[]db.Vendor, error) {
	// Query all vendors based on the received parameters
	vendors, err := QueryAllVendorsBasedOnParams(limit, offset, order, r.DB)
	if err != nil {
		fmt.Printf("Error querying db for list of vendors: %s", err)
		return nil, err
	}

	return &vendors, nil
}

// Find a vendor in database by ID
func (r *vendorRepository) FindById(id int) (*db.Vendor, error) {
	// Create an empty ref object of type vendor
	vendor := db.Vendor{}
	// Grab vendor from db if exists
	result := r.DB.First(&vendor, id)

	// If error detected
	if result.Error != nil {
		return nil, result.Error
	}
	// else
	return &vendor, nil
}

// Delete vendor in database
func (r *vendorRepository) Delete(id int) error {
	// Create an empty ref object of type vendor
	vendor := db.Vendor{}
	// Delete work type from db if exists
	result := r.DB.Delete(&vendor, id)

	// If error detected
	if result.Error != nil {
		fmt.Println("error in deleting vendor: ", result.Error)
		return result.Error
	}
	// else
	return nil
}

// Updates vendor in database
func (r *vendorRepository) Update(id int, vendor *db.Vendor) (*db.Vendor, error) {
	// Init
	var err error
	// Find vendor by id to ensure it exists
	foundVendor, err := r.FindById(id)
	if err != nil {
		fmt.Println("Vendor to update not found: ", err)
		return nil, err
	}

	// Update found vendor with incoming details
	updateResult := r.DB.Model(&foundVendor).Updates(vendor)
	if updateResult.Error != nil {
		fmt.Println("Vendor update failed: ", updateResult.Error)
		return nil, updateResult.Error
	}

	// Retrieve updated vendor by id
	updatedVendor, err := r.FindById(id)
	if err != nil {
		fmt.Println("Updated vendor not found: ", err)
		return nil, err
	}
	return updatedVendor, nil
}

// Takes limit, offset, and order parameters, builds a query and executes returning a list of vendors
func QueryAllVendorsBasedOnParams(limit int, offset int, order string, dbClient *gorm.DB) ([]db.Vendor, error) {
	// Build model to query database
	vendors := []db.Vendor{}
	// Build base query for vendors table
	query := dbClient.Model(&vendors)

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
	result := query.Find(&vendors)
	if result.Error != nil {
		return nil, result.Error
	}
	// Return if no errors with result
	return vendors, nil
}

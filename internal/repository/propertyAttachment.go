package repository

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"gorm.io/gorm"
)

type PropertyAttachmentRepository interface {
	FindAll(int, int, string) (*[]db.PropertyAttachment, error)
	FindById(int) (*db.PropertyAttachment, error)
	Create(*db.PropertyAttachment) (*db.PropertyAttachment, error)
	Update(int, *db.PropertyAttachment) (*db.PropertyAttachment, error)
	Delete(int) error
}

type propertyAttachmentRepository struct {
	DB *gorm.DB
	// objectStorage db.ObjectRepository
}

func NewPropertyAttachmentRepository(db *gorm.DB) PropertyAttachmentRepository {
	return &propertyAttachmentRepository{db}
}

// Creates a Property attachment in the database
func (r *propertyAttachmentRepository) Create(attachment *db.PropertyAttachment) (*db.PropertyAttachment, error) {
	// Create new attachment in database
	result := r.DB.Create(&attachment)
	if result.Error != nil {
		return nil, fmt.Errorf("failed creating property attachment: %w", result.Error)
	}

	return attachment, nil
}

// Find a list of attachments in the database
func (r *propertyAttachmentRepository) FindAll(limit int, offset int, order string) (*[]db.PropertyAttachment, error) {
	// Query all log messages based on the received parameters
	attachments, err := QueryAllPropertyAttachmentsBasedOnParams(limit, offset, order, r.DB)
	if err != nil {
		fmt.Printf("Error querying db for list of attachments: %s", err)
		return nil, err
	}

	return &attachments, nil
}

// Find a property attachment in database by ID
func (r *propertyAttachmentRepository) FindById(id int) (*db.PropertyAttachment, error) {
	// Create an empty ref object of type property attachment
	attachment := db.PropertyAttachment{}
	// Grab log message from db if exists
	result := r.DB.Preload("Property").First(&attachment, id)

	// If error detected
	if result.Error != nil {
		return nil, result.Error
	}
	// else
	return &attachment, nil
}

// Delete property attachment in database
func (r *propertyAttachmentRepository) Delete(id int) error {
	// Create an empty ref object of type property attachment
	attachment := db.PropertyAttachment{}
	// Delete log message from db if exists
	result := r.DB.Delete(&attachment, id)

	// If error detected
	if result.Error != nil {
		fmt.Println("error in deleting property attachment: ", result.Error)
		return result.Error
	}
	// else
	return nil
}

// Updates property attachment in database
func (r *propertyAttachmentRepository) Update(id int, attachment *db.PropertyAttachment) (*db.PropertyAttachment, error) {
	// Init
	var err error
	// Find property attachment by id to ensure it exists
	foundAttachment, err := r.FindById(id)
	if err != nil {
		fmt.Println("Property attachment to update not found: ", err)
		return nil, err
	}

	// Update found attachment
	updateResult := r.DB.Model(&foundAttachment).Updates(attachment)
	if updateResult.Error != nil {
		fmt.Println("Property attachment update failed: ", updateResult.Error)
		return nil, updateResult.Error
	}

	// Retrieve updated property attachment by id
	updatedAttachment, err := r.FindById(id)
	if err != nil {
		fmt.Println("Updated property attachment not found: ", err)
		return nil, err
	}
	return updatedAttachment, nil
}

// Takes limit, offset, and order parameters, builds a query and executes returning a list of prop attachments
func QueryAllPropertyAttachmentsBasedOnParams(limit int, offset int, order string, dbClient *gorm.DB) ([]db.PropertyAttachment, error) {
	// Build model to query database
	propAttachments := []db.PropertyAttachment{}
	// Build base query for property attachments table
	query := dbClient.Model(&propAttachments)

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
	result := query.Find(&propAttachments)
	if result.Error != nil {
		return nil, result.Error
	}
	// Return if no errors with result
	return propAttachments, nil
}

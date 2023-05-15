package repository

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"gorm.io/gorm"
)

type ContactRepository interface {
	FindAll(int, int, string) (*[]db.Contact, error)
	FindById(int) (*db.Contact, error)
	Create(*db.Contact) (*db.Contact, error)
	Update(int, *db.Contact) (*db.Contact, error)
	Delete(int) error
}

type contactRepository struct {
	DB *gorm.DB
}

func NewContactRepository(db *gorm.DB) ContactRepository {
	return &contactRepository{db}
}

// Creates a contact in the database
func (r *contactRepository) Create(contact *db.Contact) (*db.Contact, error) {
	// Create in database
	result := r.DB.Create(&contact)
	if result.Error != nil {
		return nil, fmt.Errorf("failed creating contact: %w", result.Error)
	}

	return contact, nil
}

// Find a list of contacts in the database
func (r *contactRepository) FindAll(limit int, offset int, order string) (*[]db.Contact, error) {
	// Query all based on the received parameters
	contacts, err := QueryAllContactsBasedOnParams(limit, offset, order, r.DB)
	if err != nil {
		return nil, err
	}

	return &contacts, nil
}

// Find contact in database by ID
func (r *contactRepository) FindById(id int) (*db.Contact, error) {
	// Create an empty ref object of required type
	contact := db.Contact{}
	// Grab first item with matching id
	result := r.DB.First(&contact, id)

	// If error detected
	if result.Error != nil {
		return nil, result.Error
	}
	// else return contact
	return &contact, nil
}

// Delete contact in database
func (r *contactRepository) Delete(id int) error {
	// Create an empty ref object of required type
	contact := db.Contact{}
	// Delete first item with matching id
	result := r.DB.Delete(&contact, id)

	// If error detected
	if result.Error != nil {
		return result.Error
	}
	// else return no error
	return nil
}

// Updates contact in database
func (r *contactRepository) Update(id int, contact *db.Contact) (*db.Contact, error) {
	// Init
	var err error
	// Find contact by id
	foundContact, err := r.FindById(id)
	if err != nil {
		return nil, err
	}

	// Update found contact using new struct
	updateResult := r.DB.Model(&foundContact).Updates(contact)
	if updateResult.Error != nil {
		return nil, updateResult.Error
	}

	// Retrieve updated contact by id
	updatedContact, err := r.FindById(id)
	if err != nil {
		return nil, err
	}

	return updatedContact, nil
}

// Takes limit, offset, and order parameters, builds a query and executes returning a list of contacts
func QueryAllContactsBasedOnParams(limit int, offset int, order string, dbClient *gorm.DB) ([]db.Contact, error) {
	// Build model to query database
	contacts := []db.Contact{}
	// Build base query for contacts table
	query := dbClient.Model(&contacts)

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
		// Else default to order by created_at DESC
		query.Order("created_at DESC")
	}
	// Query database
	result := query.Find(&contacts)
	if result.Error != nil {
		return nil, result.Error
	}
	// Return if no errors with result
	return contacts, nil
}

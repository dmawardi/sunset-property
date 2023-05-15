package service

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/repository"
)

type ContactService interface {
	FindAll(int, int, string) (*[]db.Contact, error)
	FindById(int) (*db.Contact, error)
	Create(*models.CreateContact) (*db.Contact, error)
	Update(int, *models.UpdateContact) (*db.Contact, error)
	Delete(int) error
}

type contactService struct {
	repo repository.ContactRepository
}

func NewContactService(repo repository.ContactRepository) ContactService {
	return &contactService{repo}
}

// Creates a contact in the database
func (s *contactService) Create(c *models.CreateContact) (*db.Contact, error) {
	// Create a new contact
	contactToCreate := db.Contact{
		FirstName:    c.FirstName,
		LastName:     c.LastName,
		Email:        c.Email,
		Phone:        c.Phone,
		Mobile:       c.Mobile,
		ContactType:  c.ContactType,
		ContactNotes: c.ContactNotes,
	}

	// Create contact in database
	createdContact, err := s.repo.Create(&contactToCreate)
	if err != nil {
		return nil, fmt.Errorf("failed creating contact: %w", err)
	}

	return createdContact, nil
}

// Find a list of contacts in the database
func (s *contactService) FindAll(limit int, offset int, order string) (*[]db.Contact, error) {

	contacts, err := s.repo.FindAll(limit, offset, order)
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

// Find contact in database by ID
func (s *contactService) FindById(id int) (*db.Contact, error) {
	// Find by id
	contact, err := s.repo.FindById(id)
	// If error detected
	if err != nil {
		return nil, err
	}
	// else
	return contact, nil
}

// Delete contact in database
func (s *contactService) Delete(id int) error {
	// Delete using id
	err := s.repo.Delete(id)
	// If error detected
	if err != nil {
		fmt.Println("error in deleting property feature: ", err)
		return err
	}
	// else
	return nil
}

// Updates contact in database
func (s *contactService) Update(id int, c *models.UpdateContact) (*db.Contact, error) {
	// Create db type from incoming DTO
	contactToUpdate := &db.Contact{
		FirstName:    c.FirstName,
		LastName:     c.LastName,
		Email:        c.Email,
		Phone:        c.Phone,
		Mobile:       c.Mobile,
		ContactType:  c.ContactType,
		ContactNotes: c.ContactNotes,
	}

	// Update using repo
	updatedContact, err := s.repo.Update(id, contactToUpdate)
	if err != nil {
		return nil, err
	}

	return updatedContact, nil
}

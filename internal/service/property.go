package service

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/repository"
)

type PropertyService interface {
	FindAll(int, int, string) (*[]db.Property, error)
	FindById(int) (*db.Property, error)
	Create(*models.CreateProperty) (*db.Property, error)
	Update(int, *models.UpdateProperty) (*db.Property, error)
	Delete(int) error
}

type propertyService struct {
	repo repository.PropertyRepository
}

func NewPropertyService(repo repository.PropertyRepository) PropertyService {
	return &propertyService{repo}
}

// Creates a property in the database
func (s *propertyService) Create(prop *models.CreateProperty) (*db.Property, error) {
	// Create a new property of type db User
	propToCreate := db.Property{
		Postcode:         prop.Postcode,
		Property_Name:    prop.Property_Name,
		Suburb:           prop.Suburb,
		City:             prop.City,
		Street_Address_1: prop.Street_Address_1,
		Street_Address_2: prop.Street_Address_2,
		Bedrooms:         prop.Bedrooms,
		Bathrooms:        prop.Bathrooms,
		Land_Area:        prop.Land_Area,
		Land_Metric:      prop.Land_Metric,
		Description:      prop.Description,
		Notes:            prop.Notes,
		Features:         prop.Features,
	}

	// Create above user in database
	createdProp, err := s.repo.Create(&propToCreate)
	if err != nil {
		return nil, fmt.Errorf("failed creating property: %w", err)
	}

	return createdProp, nil
}

// Find a list of properties in the database
func (s *propertyService) FindAll(limit int, offset int, order string) (*[]db.Property, error) {

	properties, err := s.repo.FindAll(limit, offset, order)
	if err != nil {
		return nil, err
	}
	return properties, nil
}

// Find property in database by ID
func (s *propertyService) FindById(id int) (*db.Property, error) {
	fmt.Printf("Finding property with id: %v\n", id)
	// Find user by id
	prop, err := s.repo.FindById(id)

	// If error detected
	if err != nil {
		return nil, err
	}
	// else
	return prop, nil
}

// Delete property in database
func (s *propertyService) Delete(id int) error {
	err := s.repo.Delete(id)
	// If error detected
	if err != nil {
		fmt.Println("error in deleting property: ", err)
		return err
	}
	// else
	return nil
}

// Updates property in database
func (s *propertyService) Update(id int, prop *models.UpdateProperty) (*db.Property, error) {
	// Create db property type of incoming DTO
	dbProp := &db.Property{
		Postcode:         prop.Postcode,
		Property_Name:    prop.Property_Name,
		Suburb:           prop.Suburb,
		City:             prop.City,
		Street_Address_1: prop.Street_Address_1,
		Street_Address_2: prop.Street_Address_2,
		Bedrooms:         prop.Bedrooms,
		Bathrooms:        prop.Bathrooms,
		Land_Area:        prop.Land_Area,
		Land_Metric:      prop.Land_Metric,
		Description:      prop.Description,
		Notes:            prop.Notes,
		Features:         prop.Features,
		Contacts:         prop.Contacts,
	}

	// Update using repo
	updatedProperty, err := s.repo.Update(id, dbProp)
	if err != nil {
		return nil, err
	}

	return updatedProperty, nil
}

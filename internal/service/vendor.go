package service

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/repository"
)

type VendorService interface {
	FindAll(int, int, string) (*[]db.Vendor, error)
	FindById(int) (*db.Vendor, error)
	Create(*models.CreateVendor) (*db.Vendor, error)
	Update(int, *models.UpdateVendor) (*db.Vendor, error)
	Delete(int) error
}

type vendorService struct {
	repo repository.VendorRepository
}

func NewVendorService(repo repository.VendorRepository) VendorService {
	return &vendorService{repo}
}

// Creates a vendor
func (s *vendorService) Create(vendor *models.CreateVendor) (*db.Vendor, error) {
	// Create a new vendor from DTO
	vendorToCreate := db.Vendor{
		CompanyName:      vendor.CompanyName,
		NPWP:             vendor.NPWP,
		Email:            vendor.Email,
		Phone:            vendor.Phone,
		NIB:              vendor.NIB,
		Street_Address_1: vendor.Street_Address_1,
		Street_Address_2: vendor.Street_Address_2,
		City:             vendor.City,
		Province:         vendor.Province,
		Postal_Code:      vendor.Postal_Code,
		Suburb:           vendor.Suburb,
		WorkTypes:        vendor.WorkTypes,
	}

	// Create vendor in database
	createdVendor, err := s.repo.Create(&vendorToCreate)
	if err != nil {
		return nil, fmt.Errorf("failed creating vendor: %w", err)
	}

	return createdVendor, nil
}

// Find a list of vendors
func (s *vendorService) FindAll(limit int, offset int, order string) (*[]db.Vendor, error) {
	vendors, err := s.repo.FindAll(limit, offset, order)
	if err != nil {
		return nil, err
	}
	return vendors, nil
}

// Find vendor in database by ID
func (s *vendorService) FindById(id int) (*db.Vendor, error) {
	// Find by id
	vendor, err := s.repo.FindById(id)
	// If error detected
	if err != nil {
		fmt.Println("error in finding vendor: ", err)
		return nil, err
	}
	// else
	return vendor, nil
}

// Delete vendor in database
func (s *vendorService) Delete(id int) error {
	err := s.repo.Delete(id)
	// If error detected
	if err != nil {
		fmt.Println("error in deleting vendor: ", err)
		return err
	}
	// else
	return nil
}

// Updates vendor in database
func (s *vendorService) Update(id int, vendor *models.UpdateVendor) (*db.Vendor, error) {
	// Create a new vendor from incoming DTO
	vendorToUpdate := &db.Vendor{
		CompanyName:      vendor.CompanyName,
		NPWP:             vendor.NPWP,
		Email:            vendor.Email,
		Phone:            vendor.Phone,
		NIB:              vendor.NIB,
		Street_Address_1: vendor.Street_Address_1,
		Street_Address_2: vendor.Street_Address_2,
		City:             vendor.City,
		Province:         vendor.Province,
		Postal_Code:      vendor.Postal_Code,
		Suburb:           vendor.Suburb,
		WorkTypes:        vendor.WorkTypes,
	}

	// Update using repo
	updatedVendor, err := s.repo.Update(id, vendorToUpdate)
	if err != nil {
		return nil, err
	}

	return updatedVendor, nil
}

package service

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/repository"
)

type MaintenanceRequestService interface {
	FindAll(int, int, string) (*[]db.MaintenanceRequest, error)
	FindById(int) (*db.MaintenanceRequest, error)
	Create(*models.CreateMaintenanceRequest) (*db.MaintenanceRequest, error)
	Update(int, *models.UpdateMaintenanceRequest) (*db.MaintenanceRequest, error)
	Delete(int) error
}

type maintenanceRequestService struct {
	repo repository.MaintenanceRequestRepository
}

func NewMaintenanceRequestService(repo repository.MaintenanceRequestRepository) MaintenanceRequestService {
	return &maintenanceRequestService{repo}
}

// Creates a maintenance request
func (s *maintenanceRequestService) Create(request *models.CreateMaintenanceRequest) (*db.MaintenanceRequest, error) {
	// Create a new maintenance request from DTO
	requestToCreate := db.MaintenanceRequest{
		WorkDefinition: request.WorkDefinition,
		Type:           request.Type,
		Notes:          request.Notes,
		Scale:          request.Scale,
		Tax:            request.Tax,
		TotalCost:      request.TotalCost,
		Property:       request.Property,
		TaskID:         request.Task.ID,
	}

	// Create request in database
	createdRequest, err := s.repo.Create(&requestToCreate)
	if err != nil {
		return nil, fmt.Errorf("failed creating maintenance request: %w", err)
	}

	return createdRequest, nil
}

// Find a list of maintenance requests
func (s *maintenanceRequestService) FindAll(limit int, offset int, order string) (*[]db.MaintenanceRequest, error) {
	requests, err := s.repo.FindAll(limit, offset, order)
	if err != nil {
		return nil, err
	}
	return requests, nil
}

// Find maintenance request in database by ID
func (s *maintenanceRequestService) FindById(id int) (*db.MaintenanceRequest, error) {
	// Find by id
	request, err := s.repo.FindById(id)
	// If error detected
	if err != nil {
		fmt.Println("error in finding maintenance request: ", err)
		return nil, err
	}
	// else
	return request, nil
}

// Delete maintenance request in database
func (s *maintenanceRequestService) Delete(id int) error {
	err := s.repo.Delete(id)
	// If error detected
	if err != nil {
		fmt.Println("error in deleting maintenance request: ", err)
		return err
	}
	// else
	return nil
}

// Updates maintenance request in database
func (s *maintenanceRequestService) Update(id int, request *models.UpdateMaintenanceRequest) (*db.MaintenanceRequest, error) {
	// Create a new maintenance request from DTO
	requestToUpdate := db.MaintenanceRequest{
		WorkDefinition: request.WorkDefinition,
		Type:           request.Type,
		Notes:          request.Notes,
		Scale:          request.Scale,
		Tax:            request.Tax,
		TotalCost:      request.TotalCost,
		Property:       request.Property,
	}

	// Update using repo
	updatedRequest, err := s.repo.Update(id, &requestToUpdate)
	if err != nil {
		return nil, err
	}

	return updatedRequest, nil
}

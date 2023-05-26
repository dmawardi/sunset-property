package service

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/repository"
)

type WorkTypeService interface {
	FindAll(int, int, string) (*[]db.WorkType, error)
	FindById(int) (*db.WorkType, error)
	Create(*models.CreateWorkType) (*db.WorkType, error)
	Update(int, *models.UpdateWorkType) (*db.WorkType, error)
	Delete(int) error
}

type workTypeService struct {
	repo repository.WorkTypeRepository
}

func NewWorkTypeService(repo repository.WorkTypeRepository) WorkTypeService {
	return &workTypeService{repo}
}

// Creates a work type
func (s *workTypeService) Create(workType *models.CreateWorkType) (*db.WorkType, error) {
	// Create a new work type from DTO
	workTypeToCreate := db.WorkType{
		Name: workType.Name,
	}

	// Create work type in database
	createdWorkType, err := s.repo.Create(&workTypeToCreate)
	if err != nil {
		return nil, fmt.Errorf("failed creating work type: %w", err)
	}

	return createdWorkType, nil
}

// Find a list of work types
func (s *workTypeService) FindAll(limit int, offset int, order string) (*[]db.WorkType, error) {
	workTypes, err := s.repo.FindAll(limit, offset, order)
	if err != nil {
		return nil, err
	}
	return workTypes, nil
}

// Find work type in database by ID
func (s *workTypeService) FindById(id int) (*db.WorkType, error) {
	// Find by id
	workType, err := s.repo.FindById(id)
	// If error detected
	if err != nil {
		fmt.Println("error in finding work type: ", err)
		return nil, err
	}
	// else
	return workType, nil
}

// Delete work type in database
func (s *workTypeService) Delete(id int) error {
	err := s.repo.Delete(id)
	// If error detected
	if err != nil {
		fmt.Println("error in deleting work type: ", err)
		return err
	}
	// else
	return nil
}

// Updates work type in database
func (s *workTypeService) Update(id int, workType *models.UpdateWorkType) (*db.WorkType, error) {
	// Create a new maintenance request from DTO
	workTypeToUpdate := &db.WorkType{
		Name: workType.Name,
	}

	// Update using repo
	updatedWorkType, err := s.repo.Update(id, workTypeToUpdate)
	if err != nil {
		return nil, err
	}

	return updatedWorkType, nil
}

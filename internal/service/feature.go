package service

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/repository"
)

type FeatureService interface {
	FindAll(int, int, string) (*[]db.Feature, error)
	FindById(int) (*db.Feature, error)
	Create(*models.CreateFeature) (*db.Feature, error)
	Update(int, *models.UpdateFeature) (*db.Feature, error)
	Delete(int) error
}

type featureService struct {
	repo repository.FeatureRepository
}

func NewFeatureService(repo repository.FeatureRepository) FeatureService {
	return &featureService{repo}
}

// Creates a property feature in the database
func (s *featureService) Create(feat *models.CreateFeature) (*db.Feature, error) {
	// Create a new property of type db User
	featToCreate := db.Feature{
		Feature_Name: feat.Feature_Name,
	}

	// Create above user in database
	createdFeat, err := s.repo.Create(&featToCreate)
	if err != nil {
		return nil, fmt.Errorf("failed creating property feature: %w", err)
	}

	return createdFeat, nil
}

// Find a list of property features in the database
func (s *featureService) FindAll(limit int, offset int, order string) (*[]db.Feature, error) {

	features, err := s.repo.FindAll(limit, offset, order)
	if err != nil {
		return nil, err
	}
	return features, nil
}

// Find property feature in database by ID
func (s *featureService) FindById(featId int) (*db.Feature, error) {
	fmt.Printf("Finding property feature with id: %v\n", featId)
	// Find user by id
	feature, err := s.repo.FindById(featId)
	// If error detected
	if err != nil {
		return nil, err
	}
	// else
	return feature, nil
}

// Delete property feature in database
func (s *featureService) Delete(id int) error {
	err := s.repo.Delete(id)
	// If error detected
	if err != nil {
		fmt.Println("error in deleting property feature: ", err)
		return err
	}
	// else
	return nil
}

// Updates property feature in database
func (s *featureService) Update(id int, feat *models.UpdateFeature) (*db.Feature, error) {
	// Create db property type of incoming DTO
	dbProp := &db.Feature{
		Feature_Name: feat.Feature_Name,
	}

	// Update using repo
	updatedFeature, err := s.repo.Update(id, dbProp)
	if err != nil {
		return nil, err
	}

	return updatedFeature, nil
}

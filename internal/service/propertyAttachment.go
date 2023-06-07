package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/repository"
)

type PropertyAttachmentService interface {
	FindAll(int, int, string) (*[]db.PropertyAttachment, error)
	FindById(int) (*db.PropertyAttachment, error)
	Create(*models.CreatePropertyAttachment) (*db.PropertyAttachment, error)
	AttachToProperty(propertyId uint, userUpload *http.Request) (*db.PropertyAttachment, error)
	Update(int, *models.UpdatePropertyAttachment) (*db.PropertyAttachment, error)
	Delete(int) error
}

type propertyAttachmentService struct {
	repo          repository.PropertyAttachmentRepository
	objectStorage db.ObjectRepository
}

func NewPropertyAttachmentService(repo repository.PropertyAttachmentRepository, objStorage db.ObjectRepository) PropertyAttachmentService {
	return &propertyAttachmentService{repo, objStorage}
}

// Creates a property attachment in the database
func (s *propertyAttachmentService) AttachToProperty(propertyId uint, r *http.Request) (*db.PropertyAttachment, error) {
	// Extract file from request
	file, handler, err := helpers.ExtractFileFromResponse(r)

	if err != nil {
		fmt.Println("error in extracting file from request: ", err)
		return nil, fmt.Errorf(`failed extracting file from request: %w`, err)
	}

	// Save a copy of the file on the server
	err = helpers.SaveACopyOfTheFileOnTheServer(file, handler, "./temp")
	if err != nil {
		fmt.Println("error in saving a copy of the file on the server: ", err)
		return nil, fmt.Errorf(`failed saving a copy of the file on the server: %w`, err)
	}

	// Build the required details for saving
	fileName := handler.Filename
	fileExtension := strings.Split(fileName, ".")[1]
	tempFilePath := "./temp/" + fileName

	// Build file key path based on prop attachment requirements
	fileKeyPath := fmt.Sprintf("property/%v/attachments/%s", propertyId, fileName)

	// Upload file to object storage
	eTag, fileSize, err := s.objectStorage.UploadFile(tempFilePath, fileKeyPath, false)
	if err != nil {
		fmt.Println("error in uploading file to object storage: ", err)
		return nil, fmt.Errorf(`failed uploading file to object storage: %w`, err)
	}

	fmt.Println("Successfully uploaded file to object storage")

	// Create a new property attachment
	attachmentToCreate := db.PropertyAttachment{
		Label:     fileName,
		FileName:  fileName,
		FileSize:  fileSize,
		FileType:  fileExtension,
		ObjectKey: fileKeyPath,
		ETag:      eTag,
		Property: db.Property{
			ID: propertyId,
		},
	}

	// Create attachment
	createdAttachment, err := s.repo.Create(&attachmentToCreate)
	if err != nil {
		return nil, fmt.Errorf("failed creating property attachment: %w", err)
	}

	// Delete temp file
	err = helpers.DeleteFile(tempFilePath)
	if err != nil {
		fmt.Println("error in deleting temp file: ", err)
		return nil, fmt.Errorf(`failed deleting temp file: %w`, err)
	}

	return createdAttachment, nil
}

// Creates a property attachment
func (s *propertyAttachmentService) Create(attachment *models.CreatePropertyAttachment) (*db.PropertyAttachment, error) {
	// Create a new attachment from DTO
	attachmentToCreate := db.PropertyAttachment{
		Label:     attachment.Label,
		FileName:  attachment.FileName,
		FileSize:  attachment.FileSize,
		FileType:  attachment.FileType,
		ETag:      attachment.ETag,
		ObjectKey: attachment.ObjectKey,
		Property:  attachment.Property,
	}

	// Create property attachment in database
	createdAttachment, err := s.repo.Create(&attachmentToCreate)
	if err != nil {
		return nil, fmt.Errorf("failed creating property attachment: %w", err)
	}

	return createdAttachment, nil
}

// Find a list of property attachments in the database
func (s *propertyAttachmentService) FindAll(limit int, offset int, order string) (*[]db.PropertyAttachment, error) {
	attachments, err := s.repo.FindAll(limit, offset, order)
	if err != nil {
		return nil, err
	}
	return attachments, nil
}

// Find property attachment in database by ID
func (s *propertyAttachmentService) FindById(id int) (*db.PropertyAttachment, error) {
	// Find attachment by id
	attachment, err := s.repo.FindById(id)
	// If error detected
	if err != nil {
		return nil, err
	}
	// else
	return attachment, nil
}

// Delete property attachment in database
func (s *propertyAttachmentService) Delete(id int) error {
	err := s.repo.Delete(id)
	// If error detected
	if err != nil {
		fmt.Println("error in deleting property attachment: ", err)
		return err
	}
	// else
	return nil
}

// Updates property attachment in database (only label can be updated)
func (s *propertyAttachmentService) Update(id int, attachment *models.UpdatePropertyAttachment) (*db.PropertyAttachment, error) {
	// Create db Property attachment type from DTO
	attachToCreate := db.PropertyAttachment{
		Label: attachment.Label,
	}

	// Update using repo
	updatedAttachment, err := s.repo.Update(id, &attachToCreate)
	if err != nil {
		return nil, err
	}

	return updatedAttachment, nil
}

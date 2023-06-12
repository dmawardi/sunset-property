package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/service"
	"github.com/go-chi/chi"
)

type PropertyAttachmentController interface {
	Upload(w http.ResponseWriter, r *http.Request)
	Download(w http.ResponseWriter, r *http.Request)
	// HandleUpload(w http.ResponseWriter, r *http.Request)
	FindAll(w http.ResponseWriter, r *http.Request)
	Find(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type propertyAttachmentController struct {
	service     service.PropertyAttachmentService
	propService service.PropertyService
}

func NewPropertyAttachmentController(service service.PropertyAttachmentService, propService service.PropertyService) PropertyAttachmentController {
	return &propertyAttachmentController{service, propService}
}

// API/PROPERTY-ATTACH/{propertyId}
// Attaches a file to a property
// @Summary      Attaches a file to a property
// @Description  Accepts a propertyId parameter and a file delivered by form-data with a key of "file"
// @Tags         Property Attachments
// @Accept       json
// @Produce      json
// @Param        propertyId   path      int  true  "Property ID"
// @Success      200 {object} string "Property attachment upload successful!"
// @Failure      400 {string} string "Can't find property with ID: {id}"
// @Failure      400 {string} string "Invalid property ID"
// @Router       /property-attach [post]
// @Security BearerToken
func (c propertyAttachmentController) Upload(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "propertyId")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		http.Error(w, "Invalid property ID", http.StatusBadRequest)
		return
	}
	// Query database for property attachment using ID
	_, err = c.propService.FindById(idParameter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find property with ID: %v\n", idParameter), http.StatusBadRequest)
		return
	}
	// if no error, proceed to upload file
	// Get the file from the request
	_, createErr := c.service.AttachToProperty(uint(idParameter), r)

	if createErr != nil {
		fmt.Printf("Issue with prop attachment creation: %v\n", createErr)
		http.Error(w, "Property attachment creation failed.", http.StatusBadRequest)
		return
	}

	// Set status to created
	w.WriteHeader(http.StatusCreated)
	// Send user success message in body
	w.Write([]byte("Property attachment upload successful!"))
}

func (c propertyAttachmentController) Download(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	// Query database for property attachment using ID and download if found
	downloadedFilePath, err := c.service.DownloadPropertyAttachment(idParameter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find property attachment with ID: %v\n", idParameter), http.StatusBadRequest)
		return
	}
	// read downloaded file
	file, err := helpers.ReadFile(downloadedFilePath)
	if err != nil {
		http.Error(w, "Failed to download file", http.StatusInternalServerError)
		return
	}
	// Copy the file contents to the response writer
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Failed to download file", http.StatusInternalServerError)
		return
	}
	// Set status to OK
	w.WriteHeader(http.StatusOK)

	// Delete the file from the server
	err = helpers.DeleteFile(downloadedFilePath)
	if err != nil {
		fmt.Printf("Error deleting temporary file: %v\n", err)
		return
	}

}

// func (c propertyAttachmentController) HandleUpload(w http.ResponseWriter, r *http.Request) {
// 	// Get the file from the request
// 	file, handler, err := helpers.ExtractFileFromResponse(r)
// 	if err != nil {
// 		http.Error(w, "Failed to extract file content", http.StatusInternalServerError)
// 		return
// 	}

// 	// Save a temporary copy of the file in temp folder
// 	err = helpers.SaveACopyOfTheFileOnTheServer(file, handler, "./temp/")
// 	if err != nil {
// 		fmt.Printf("Error saving file: %v\n", err)
// 		http.Error(w, "Failed to save file", http.StatusInternalServerError)
// 		return
// 	}

// 	fmt.Fprintf(w, "File uploaded successfully!")
// }

// API/PROPERTY-ATTACHMENTS
// Find a list of Property attachments
// @Summary      Find a list of property attachments
// @Description  Accepts limit, offset, and order params and returns list of attachments
// @Tags         Property Attachments
// @Accept       json
// @Produce      json
// @Param        limit   path      int  true  "limit"
// @Param        offset   path      int  true  "offset"
// @Param        order   path      int  true  "order by"
// @Success      200 {object} []db.PropertyAttachment
// @Failure      400 {string} string "Can't find property attachments"
// @Failure      400 {string} string "Must include limit parameter with a max value of 50"
// @Router       /property-attachments [get]
// @Security BearerToken
func (c propertyAttachmentController) FindAll(w http.ResponseWriter, r *http.Request) {
	// Grab URL query parameters
	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")
	orderBy := r.URL.Query().Get("order")

	// Convert to int
	limit, _ := strconv.Atoi(limitParam)
	offset, _ := strconv.Atoi(offsetParam)

	// Check that limit is present as requirement
	if (limit == 0) || (limit >= 50) {
		http.Error(w, "Must include limit parameter with a max value of 50", http.StatusBadRequest)
		return
	}

	// Query database for all attachments using query params
	found, err := c.service.FindAll(limit, offset, orderBy)
	if err != nil {
		http.Error(w, "Can't find property attachments", http.StatusBadRequest)
		return
	}
	// Write found attachments to response
	err = helpers.WriteAsJSON(w, found)
	if err != nil {
		fmt.Println("error writing property attachments to response: ", err)
		http.Error(w, "Can't find property attachments", http.StatusBadRequest)
		return
	}
}

// Find a created property attachment by ID
// @Summary      Find property attachment by ID
// @Description  Find a property attachment by ID
// @Tags         Property Attachments
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Property Attachment ID"
// @Success      200 {object} db.PropertyAttachment
// @Failure      400 {string} string "Can't find property attachment with ID:"
// @Failure      400 {string} string "Invalid ID"
// @Router       /property-attachment/{id} [get]
// @Security BearerToken
func (c propertyAttachmentController) Find(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	// Query database for property attachment using ID
	found, err := c.service.FindById(idParameter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find property attachment with ID: %v\n", idParameter), http.StatusBadRequest)
		return
	}
	// Write found property attachment to response
	err = helpers.WriteAsJSON(w, found)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find property attachment with ID: %v\n", idParameter), http.StatusBadRequest)
		return
	}
}

// Create a new property attachment
// @Summary      Create property attachment
// @Description  Creates a new property attachment
// @Tags         Property Attachments
// @Accept       json
// @Produce      json
// @Param        feature body models.CreatePropertyAttachment true "New Property Attachment Json"
// @Success      201 {string} string "Property attachment creation successful!"
// @Failure      400 {string} string "Property attachment creation failed."
// @Router       /property-attachment [post]
func (c propertyAttachmentController) Create(w http.ResponseWriter, r *http.Request) {
	// Init
	var attachment models.CreatePropertyAttachment
	// Decode request body as JSON and store
	err := json.NewDecoder(r.Body).Decode(&attachment)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&attachment)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		helpers.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through
	// Create property attachment
	_, createErr := c.service.Create(&attachment)
	if createErr != nil {
		fmt.Printf("Issue with prop attachment creation: %v\n", createErr)
		http.Error(w, "Property attachment creation failed.", http.StatusBadRequest)
		return
	}

	// Set status to created
	w.WriteHeader(http.StatusCreated)
	// Send user success message in body
	w.Write([]byte("Property attachment creation successful!"))
}

// Update a property attachment (using URL parameter id)
// @Summary      Update property attachment
// @Description  Updates an existing property attachment
// @Tags         Property Attachments
// @Accept       json
// @Produce      json
// @Param        feature body models.UpdatePropertyAttachment true "Update Property Attachment Json"
// @Param        id   path      int  true  "Property Attachment ID"
// @Success      200 {object} db.PropertyAttachment
// @Failure      400 {string} string "Failed property attachment update"
// @Failure      400 {string} string "Invalid ID"
// @Failure      403 {string} string "Authentication Token not detected"
// @Router       /property-attachments/{id} [put]
// @Security BearerToken
func (c propertyAttachmentController) Update(w http.ResponseWriter, r *http.Request) {
	// Grab ID URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Init
	var attachment models.UpdatePropertyAttachment
	// Decode request body as JSON and store
	err = json.NewDecoder(r.Body).Decode(&attachment)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&attachment)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		helpers.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	// Update property attachment
	updatedAttachment, createErr := c.service.Update(idParameter, &attachment)
	if createErr != nil {
		http.Error(w, fmt.Sprintf("Failed property attachment update: %s", createErr), http.StatusBadRequest)
		return
	}
	// Write property attachment to output
	err = helpers.WriteAsJSON(w, updatedAttachment)
	if err != nil {
		fmt.Printf("Error encountered when writing to JSON. Err: %s", err)
	}
}

// Delete property attachment (using URL parameter id)
// @Summary      Delete Property Attachment
// @Description  Deletes an existing property attachment
// @Tags         Property Attachments
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Property attachment ID"
// @Success      200 {string} string "Deletion successful!"
// @Failure      400 {string} string "Failed property attachment deletion"
// @Failure      400 {string} string "Invalid ID"
// @Router       /property-attachments/{id} [delete]
// @Security BearerToken
func (c propertyAttachmentController) Delete(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Attampt to delete property attachment using id
	err = c.service.Delete(idParameter)

	// If error detected
	if err != nil {
		http.Error(w, "Failed property attachment deletion", http.StatusBadRequest)
		return
	}
	// Else write success
	w.Write([]byte("Deletion successful!"))
}

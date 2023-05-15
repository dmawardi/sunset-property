package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/service"
	"github.com/go-chi/chi"
)

type ContactController interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	Find(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type contactController struct {
	service service.ContactService
}

func NewContactController(service service.ContactService) ContactController {
	return &contactController{service}
}

// API/CONTACTS
// Find a list of contacts
// @Summary      Find List of Contacts
// @Description  Accepts limit, offset, and order params and returns list of contacts
// @Tags         Contacts
// @Accept       json
// @Produce      json
// @Param        limit   path      int  true  "limit"
// @Param        offset   path      int  true  "offset"
// @Param        order   path      int  true  "order by"
// @Success      200 {object} []db.Contact
// @Failure      400 {string} string "Can't find contacts"
// @Failure      400 {string} string "Must include limit parameter with a max value of 50"
// @Router       /contacts [get]
// @Security BearerToken
func (c contactController) FindAll(w http.ResponseWriter, r *http.Request) {
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

	// Query database for list of contacts using query params
	foundFeatures, err := c.service.FindAll(limit, offset, orderBy)
	if err != nil {
		http.Error(w, "Can't find contacts", http.StatusBadRequest)
		return
	}

	// Write response
	err = helpers.WriteAsJSON(w, foundFeatures)
	if err != nil {
		http.Error(w, "Can't find contacts", http.StatusBadRequest)
		return
	}
}

// Find a created contact
// @Summary      Find Contact
// @Description  Find a contact by ID
// @Tags         Contacts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Contact ID"
// @Success      200 {object} db.Feature
// @Failure      400 {string} string "Can't find contact with ID:"
// @Router       /contacts/{id} [get]
// @Security BearerToken
func (c contactController) Find(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Query database for contact using id
	foundContact, err := c.service.FindById(idParameter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find contact with ID: %v\n", idParameter), http.StatusBadRequest)
		return
	}
	// Write response
	err = helpers.WriteAsJSON(w, foundContact)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find contact with ID: %v\n", idParameter), http.StatusBadRequest)
		return
	}
}

// Create a new contact
// @Summary      Create Contact
// @Description  Creates a new contact
// @Tags         Contacts
// @Accept       json
// @Produce      json
// @Param        contact body models.CreateContact true "New Contact Json"
// @Success      201 {string} string "Contact creation successful!"
// @Failure      400 {string} string "Contact creation failed."
// @Router       /contacts [post]
func (c contactController) Create(w http.ResponseWriter, r *http.Request) {
	// Init
	var contact models.CreateContact
	// Decode request body as JSON and store in contact
	err := json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&contact)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		helpers.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	// Create contact
	_, createErr := c.service.Create(&contact)
	if createErr != nil {
		http.Error(w, "Contact creation failed.", http.StatusBadRequest)
		return
	}

	// Set status to created
	w.WriteHeader(http.StatusCreated)
	// Send user success message in body
	w.Write([]byte("Contact creation successful!"))
}

// Update a contact (using URL parameter id)
// @Summary      Update Contact
// @Description  Updates an existing contact
// @Tags         Contacts
// @Accept       json
// @Produce      json
// @Param        contact body models.UpdateContact true "Update Contact Json"
// @Param        id   path      int  true  "Contact ID"
// @Success      200 {object} db.Contact
// @Failure      400 {string} string "Failed contact update"
// @Failure      403 {string} string "Authentication Token not detected"
// @Router       /contacts/{id} [put]
// @Security BearerToken
func (c contactController) Update(w http.ResponseWriter, r *http.Request) {
	// grab id parameter
	var contact models.UpdateContact
	// Decode request body as JSON and store in contact
	err := json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&contact)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		helpers.WriteAsJSON(w, valErrors)
		return
	}

	// else, validation passes and allow through
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, _ := strconv.Atoi(stringParameter)

	// Update contact
	updatedContact, createErr := c.service.Update(idParameter, &contact)
	if createErr != nil {
		http.Error(w, fmt.Sprintf("Failed contact update: %s", createErr), http.StatusBadRequest)
		return
	}
	// Write property feature to output
	err = helpers.WriteAsJSON(w, updatedContact)
	if err != nil {
		fmt.Printf("Error encountered when writing to JSON. Err: %s", err)
	}
}

// Delete contact (using URL parameter id)
// @Summary      Delete Contact
// @Description  Deletes an existing contact
// @Tags         Contacts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Contact ID"
// @Success      200 {string} string "Deletion successful!"
// @Failure      400 {string} string "Failed contact deletion"
// @Router       /contacts/{id} [delete]
// @Security BearerToken
func (c contactController) Delete(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, _ := strconv.Atoi(stringParameter)

	// Delete using id
	err := c.service.Delete(idParameter)

	// If error detected
	if err != nil {
		http.Error(w, "Failed property feature deletion", http.StatusBadRequest)
		return
	}
	// Else write success
	w.Write([]byte("Deletion successful!"))
	return
}

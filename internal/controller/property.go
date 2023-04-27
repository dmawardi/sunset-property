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

type PropertyController interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	Find(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type propertyController struct {
	service service.PropertyService
}

func NewPropertyController(service service.PropertyService) PropertyController {
	return &propertyController{service}
}

// API/PROPERTIES
// Find a list of properties
// @Summary      Find a list of properties
// @Description  Accepts limit, offset, and order params and returns list of properties
// @Tags         Property
// @Accept       json
// @Produce      json
// @Param        limit   path      int  true  "limit"
// @Param        offset   path      int  true  "offset"
// @Param        order   path      int  true  "order by"
// @Success      200 {object} []db.Property
// @Failure      400 {string} string "Can't find properties"
// @Failure      400 {string} string "Must include limit parameter with a max value of 50"
// @Router       /properties [get]
// @Security BearerToken
func (c propertyController) FindAll(w http.ResponseWriter, r *http.Request) {
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

	// Query database for all properties using query params
	foundProperties, err := c.service.FindAll(limit, offset, orderBy)
	if err != nil {
		http.Error(w, "Can't find properties", http.StatusBadRequest)
		return
	}
	err = helpers.WriteAsJSON(w, foundProperties)
	if err != nil {
		http.Error(w, "Can't find properties", http.StatusBadRequest)
		fmt.Println("error writing properties to response: ", err)
		return
	}
}

// Find a created property
// @Summary      Find Property
// @Description  Find a property by ID
// @Tags         Property
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Property ID"
// @Success      200 {object} db.Property
// @Failure      400 {string} string "Can't find property"
// @Router       /properties/{id} [get]
// @Security BearerToken
func (c propertyController) Find(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	foundProperty, err := c.service.FindById(idParameter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find property with ID: %v\n", idParameter), http.StatusBadRequest)
		return
	}
	err = helpers.WriteAsJSON(w, foundProperty)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find property with ID: %v\n", idParameter), http.StatusBadRequest)
		return
	}
}

// Create a new property
// @Summary      Create Property
// @Description  Creates a new property
// @Tags         Property
// @Accept       json
// @Produce      json
// @Param        property body models.CreateProperty true "NewPropertyJson"
// @Success      201 {string} string "Property creation successful!"
// @Failure      400 {string} string "Property creation failed."
// @Router       /properties [post]
func (c propertyController) Create(w http.ResponseWriter, r *http.Request) {
	// Init
	var prop models.CreateProperty
	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&prop)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&prop)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		helpers.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	// Create property
	_, createErr := c.service.Create(&prop)
	if createErr != nil {
		http.Error(w, "Property creation failed.", http.StatusBadRequest)
		return
	}

	// Set status to created
	w.WriteHeader(http.StatusCreated)
	// Send user success message in body
	w.Write([]byte("Property creation successful!"))
}

// Update a property (using URL parameter id)
// @Summary      Update Property
// @Description  Updates an existing property
// @Tags         Property
// @Accept       json
// @Produce      json
// @Param        property body models.UpdateProperty true "Update Property Json"
// @Param        id   path      int  true  "Property ID"
// @Success      200 {object} db.Property
// @Failure      400 {string} string "Failed property update"
// @Failure      403 {string} string "Authentication Token not detected"
// @Router       /properties/{id} [put]
// @Security BearerToken
func (c propertyController) Update(w http.ResponseWriter, r *http.Request) {
	// grab id parameter
	var prop models.UpdateProperty
	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&prop)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&prop)
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

	// Update property
	updatedProperty, createErr := c.service.Update(idParameter, &prop)
	if createErr != nil {
		http.Error(w, fmt.Sprintf("Failed property update: %s", createErr), http.StatusBadRequest)
		return
	}
	// Write property to output
	err = helpers.WriteAsJSON(w, updatedProperty)
	if err != nil {
		fmt.Printf("Error encountered when writing to JSON. Err: %s", err)
	}
}

// Delete property (using URL parameter id)
// @Summary      Delete Property
// @Description  Deletes an existing property
// @Tags         Property
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Property ID"
// @Success      200 {string} string "Deletion successful!"
// @Failure      400 {string} string "Failed property deletion"
// @Router       /properties/{id} [delete]
// @Security BearerToken
func (c propertyController) Delete(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, _ := strconv.Atoi(stringParameter)

	// Attampt to delete property using id
	err := c.service.Delete(idParameter)

	// If error detected
	if err != nil {
		http.Error(w, "Failed property deletion", http.StatusBadRequest)
		return
	}
	// Else write success
	w.Write([]byte("Deletion successful!"))
	return
}

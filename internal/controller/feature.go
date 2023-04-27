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

type FeatureController interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	Find(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type featureController struct {
	service service.FeatureService
}

func NewFeatureController(service service.FeatureService) FeatureController {
	return &featureController{service}
}

// API/FEATURES
// Find a list of Property Features
// @Summary      Find a list of property features
// @Description  Accepts limit, offset, and order params and returns list of features
// @Tags         Property Feature
// @Accept       json
// @Produce      json
// @Param        limit   path      int  true  "limit"
// @Param        offset   path      int  true  "offset"
// @Param        order   path      int  true  "order by"
// @Success      200 {object} []db.Feature
// @Failure      400 {string} string "Can't find property features"
// @Failure      400 {string} string "Must include limit parameter with a max value of 50"
// @Router       /features [get]
// @Security BearerToken
func (c featureController) FindAll(w http.ResponseWriter, r *http.Request) {
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

	// Query database for all features using query params
	foundFeatures, err := c.service.FindAll(limit, offset, orderBy)
	if err != nil {
		http.Error(w, "Can't find property features", http.StatusBadRequest)
		return
	}
	err = helpers.WriteAsJSON(w, foundFeatures)
	if err != nil {
		http.Error(w, "Can't find property features", http.StatusBadRequest)
		fmt.Println("error writing features to response: ", err)
		return
	}
}

// Find a created property feature
// @Summary      Find Property Feature
// @Description  Find a property feature by ID
// @Tags         Property Feature
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Feature ID"
// @Success      200 {object} db.Feature
// @Failure      400 {string} string "Can't find property feature with ID:"
// @Router       /features/{id} [get]
// @Security BearerToken
func (c featureController) Find(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	fmt.Println("id parameter from request: ", stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	foundProperty, err := c.service.FindById(idParameter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find property feature with ID: %v\n", idParameter), http.StatusBadRequest)
		return
	}
	err = helpers.WriteAsJSON(w, foundProperty)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find property feature with ID: %v\n", idParameter), http.StatusBadRequest)
		return
	}
}

// Create a new property feature
// @Summary      Create Property Feature
// @Description  Creates a new property feature
// @Tags         Property Feature
// @Accept       json
// @Produce      json
// @Param        feature body models.CreateFeature true "New Feature Json"
// @Success      201 {string} string "Property feature creation successful!"
// @Failure      400 {string} string "Property feature creation failed."
// @Router       /features [post]
func (c featureController) Create(w http.ResponseWriter, r *http.Request) {
	// Init
	var feat models.CreateFeature
	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&feat)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&feat)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		helpers.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	// Create property feature
	_, createErr := c.service.Create(&feat)
	if createErr != nil {
		http.Error(w, "Property feature creation failed.", http.StatusBadRequest)
		return
	}

	// Set status to created
	w.WriteHeader(http.StatusCreated)
	// Send user success message in body
	w.Write([]byte("Property feature creation successful!"))
}

// Update a property feature (using URL parameter id)
// @Summary      Update Property Feature
// @Description  Updates an existing property feature
// @Tags         Property Feature
// @Accept       json
// @Produce      json
// @Param        feature body models.UpdateFeature true "Update Feature Json"
// @Param        id   path      int  true  "Feature ID"
// @Success      200 {object} db.Feature
// @Failure      400 {string} string "Failed property feature update"
// @Failure      403 {string} string "Authentication Token not detected"
// @Router       /features/{id} [put]
// @Security BearerToken
func (c featureController) Update(w http.ResponseWriter, r *http.Request) {
	// grab id parameter
	var feat models.UpdateFeature
	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&feat)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&feat)
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

	// Update property feature
	updatedFeat, createErr := c.service.Update(idParameter, &feat)
	if createErr != nil {
		http.Error(w, fmt.Sprintf("Failed property feature update: %s", createErr), http.StatusBadRequest)
		return
	}
	// Write property feature to output
	err = helpers.WriteAsJSON(w, updatedFeat)
	if err != nil {
		fmt.Printf("Error encountered when writing to JSON. Err: %s", err)
	}
}

// Delete property feature (using URL parameter id)
// @Summary      Delete Property Feature
// @Description  Deletes an existing property feature
// @Tags         Property Feature
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Feature ID"
// @Success      200 {string} string "Deletion successful!"
// @Failure      400 {string} string "Failed property feature deletion"
// @Router       /features/{id} [delete]
// @Security BearerToken
func (c featureController) Delete(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, _ := strconv.Atoi(stringParameter)

	// Attampt to delete user using id
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

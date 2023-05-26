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

type WorkTypeController interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	Find(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type workTypeController struct {
	service service.WorkTypeService
}

func NewWorkTypeController(service service.WorkTypeService) WorkTypeController {
	return &workTypeController{service}
}

// API/WORK-TYPES
// Find a list of work types
// @Summary      Find a list of work types
// @Description  Accepts limit, offset, and order params and returns list of work types
// @Tags         Work Types
// @Accept       json
// @Produce      json
// @Param        limit   path      int  true  "limit"
// @Param        offset   path      int  true  "offset"
// @Param        order   path      int  true  "order by"
// @Success      200 {object} []db.WorkType
// @Failure      400 {string} string "Can't find work types"
// @Failure      400 {string} string "Must include limit parameter with a max value of 50"
// @Router       /work-types [get]
// @Security BearerToken
func (c workTypeController) FindAll(w http.ResponseWriter, r *http.Request) {
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

	// Query database for all work types using query params
	found, err := c.service.FindAll(limit, offset, orderBy)
	if err != nil {
		http.Error(w, "Can't find work types", http.StatusBadRequest)
		return
	}
	// Write found work types to response
	err = helpers.WriteAsJSON(w, found)
	if err != nil {
		http.Error(w, "Can't find work types", http.StatusBadRequest)
		fmt.Println("error writing work types to response: ", err)
		return
	}
}

// Find a created work type by ID
// @Summary      Find work type by ID
// @Description  Find a work type by ID
// @Tags         Work Types
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Work type ID"
// @Success      200 {object} db.WorkType
// @Failure      400 {string} string "Can't find work type with ID:"
// @Failure      400 {string} string "Invalid ID"
// @Router       /work-types/{id} [get]
// @Security BearerToken
func (c workTypeController) Find(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	// Query database for work type using ID
	found, err := c.service.FindById(idParameter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find work type with ID: %v\n%v", idParameter, err), http.StatusBadRequest)
		return
	}
	// Write found item to response
	err = helpers.WriteAsJSON(w, found)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find work type with ID: %v\n", idParameter), http.StatusBadRequest)
		return
	}
}

// Create a new work type
// @Summary      Create work type
// @Description  Creates a new work type
// @Tags         Work Types
// @Accept       json
// @Produce      json
// @Param        request body models.CreateWorkType true "New Work Type Json"
// @Success      201 {string} string "Work type creation successful!"
// @Failure      400 {string} string "Work type creation failed."
// @Router       /work-types [post]
func (c workTypeController) Create(w http.ResponseWriter, r *http.Request) {
	// Init
	var workType models.CreateWorkType
	// Decode request body from JSON and store
	err := json.NewDecoder(r.Body).Decode(&workType)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&workType)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		helpers.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	// Create work type in db
	_, createErr := c.service.Create(&workType)
	if createErr != nil {
		http.Error(w, "Work type creation failed:."+createErr.Error(), http.StatusBadRequest)
		return
	}

	// Set status to created
	w.WriteHeader(http.StatusCreated)
	// Send user success message in body
	w.Write([]byte("Work type creation successful!"))
}

// Update a work type (using URL parameter id)
// @Summary      Update work type
// @Description  Updates an existing work type
// @Tags         Work Types
// @Accept       json
// @Produce      json
// @Param        request body models.UpdateWorkType true "Update Work Type Json"
// @Param        id   path      int  true  "Work Type ID"
// @Success      200 {object} db.MaintenanceRequest
// @Failure      400 {string} string "Failed work type update"
// @Failure      400 {string} string "Invalid ID"
// @Failure      403 {string} string "Authentication Token not detected"
// @Router       /work-types/{id} [put]
// @Security BearerToken
func (c workTypeController) Update(w http.ResponseWriter, r *http.Request) {
	// Grab ID URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Init
	var workType models.UpdateWorkType
	// Decode request body from JSON and store
	err = json.NewDecoder(r.Body).Decode(&workType)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&workType)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		helpers.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	// Update work type
	updatedRequest, createErr := c.service.Update(idParameter, &workType)
	if createErr != nil {
		http.Error(w, fmt.Sprintf("Failed work type update: %s", createErr), http.StatusBadRequest)
		return
	}
	// Write work type to output
	err = helpers.WriteAsJSON(w, updatedRequest)
	if err != nil {
		fmt.Printf("Error encountered when writing to JSON. Err: %s", err)
	}
}

// Delete work type (using URL parameter id)
// @Summary      Delete work type
// @Description  Deletes an existing work type
// @Tags         Work Types
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Work Type ID"
// @Success      200 {string} string "Deletion successful!"
// @Failure      400 {string} string "Failed work type deletion"
// @Failure      400 {string} string "Invalid ID"
// @Router       /work-types/{id} [delete]
// @Security BearerToken
func (c workTypeController) Delete(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Attampt to delete work type using id
	err = c.service.Delete(idParameter)

	// If error detected
	if err != nil {
		http.Error(w, "Failed work type deletion", http.StatusBadRequest)
		return
	}
	// Else write success
	w.Write([]byte("Deletion successful!"))
}

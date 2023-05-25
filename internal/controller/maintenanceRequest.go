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

type MaintenanceRequestController interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	Find(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type maintenanceRequestController struct {
	service service.MaintenanceRequestService
}

func NewMaintenanceRequestController(service service.MaintenanceRequestService) MaintenanceRequestController {
	return &maintenanceRequestController{service}
}

// API/MAINTENANCE
// Find a list of maintenance requests
// @Summary      Find a list of maintenance requests
// @Description  Accepts limit, offset, and order params and returns list of maintenance requests
// @Tags         Maintenance Requests
// @Accept       json
// @Produce      json
// @Param        limit   path      int  true  "limit"
// @Param        offset   path      int  true  "offset"
// @Param        order   path      int  true  "order by"
// @Success      200 {object} []db.MaintenanceRequest
// @Failure      400 {string} string "Can't find maintenance requests"
// @Failure      400 {string} string "Must include limit parameter with a max value of 50"
// @Router       /maintenance [get]
// @Security BearerToken
func (c maintenanceRequestController) FindAll(w http.ResponseWriter, r *http.Request) {
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

	// Query database for all maintenance requests using query params
	found, err := c.service.FindAll(limit, offset, orderBy)
	if err != nil {
		http.Error(w, "Can't find maintenance requests", http.StatusBadRequest)
		return
	}
	// Write found maintenance requests to response
	err = helpers.WriteAsJSON(w, found)
	if err != nil {
		http.Error(w, "Can't find maintenance requests", http.StatusBadRequest)
		fmt.Println("error writing maintenance requests to response: ", err)
		return
	}
}

// Find a created maintenance request by ID
// @Summary      Find maintenence request by ID
// @Description  Find a maintenance request by ID
// @Tags         Maintenance Requests
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Maintenance Request ID"
// @Success      200 {object} db.MaintenanceRequest
// @Failure      400 {string} string "Can't find maintenance request with ID:"
// @Failure      400 {string} string "Invalid ID"
// @Router       /maintenance/{id} [get]
// @Security BearerToken
func (c maintenanceRequestController) Find(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	// Query database for maintenance request using ID
	found, err := c.service.FindById(idParameter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find maintenance request with ID: %v\n%v", idParameter, err), http.StatusBadRequest)
		return
	}
	// Write found transaction to response
	err = helpers.WriteAsJSON(w, found)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find maintenance request with ID: %v\n", idParameter), http.StatusBadRequest)
		return
	}
}

// Create a new maintenance request
// @Summary      Create maintenance request
// @Description  Creates a new maintenance request
// @Tags         Maintenance Requests
// @Accept       json
// @Produce      json
// @Param        request body models.CreateMaintenanceRequest true "New Maintenance Request Json"
// @Success      201 {string} string "Maintenance request creation successful!"
// @Failure      400 {string} string "Maintenance request creation failed."
// @Router       /maintenance [post]
func (c maintenanceRequestController) Create(w http.ResponseWriter, r *http.Request) {
	// Init
	var request models.CreateMaintenanceRequest
	// Decode request body from JSON and store
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&request)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		helpers.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	// Create maintenance request in db
	_, createErr := c.service.Create(&request)
	if createErr != nil {
		http.Error(w, "Maintenance request creation failed:."+createErr.Error(), http.StatusBadRequest)
		return
	}

	// Set status to created
	w.WriteHeader(http.StatusCreated)
	// Send user success message in body
	w.Write([]byte("Maintenance request creation successful!"))
}

// Update a maintenance request (using URL parameter id)
// @Summary      Update maintenance request
// @Description  Updates an existing maintenance request
// @Tags         Maintenance Requests
// @Accept       json
// @Produce      json
// @Param        request body models.UpdateMaintenanceRequest true "Update Maintenance Request Json"
// @Param        id   path      int  true  "Maintenance Request ID"
// @Success      200 {object} db.MaintenanceRequest
// @Failure      400 {string} string "Failed maintenance request update"
// @Failure      400 {string} string "Invalid ID"
// @Failure      403 {string} string "Authentication Token not detected"
// @Router       /maintenance/{id} [put]
// @Security BearerToken
func (c maintenanceRequestController) Update(w http.ResponseWriter, r *http.Request) {
	// Grab ID URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Init
	var maintenanceRequest models.UpdateMaintenanceRequest
	// Decode request body from JSON and store
	err = json.NewDecoder(r.Body).Decode(&maintenanceRequest)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&maintenanceRequest)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		helpers.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	// Update maintenance request
	updatedRequest, createErr := c.service.Update(idParameter, &maintenanceRequest)
	if createErr != nil {
		http.Error(w, fmt.Sprintf("Failed maintenance request update: %s", createErr), http.StatusBadRequest)
		return
	}
	// Write maintenance request to output
	err = helpers.WriteAsJSON(w, updatedRequest)
	if err != nil {
		fmt.Printf("Error encountered when writing to JSON. Err: %s", err)
	}
}

// Delete maintenance request (using URL parameter id)
// @Summary      Delete maintenance request
// @Description  Deletes an existing maintenance request
// @Tags         Maintenance Requests
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Maintenance request ID"
// @Success      200 {string} string "Deletion successful!"
// @Failure      400 {string} string "Failed maintenance request deletion"
// @Failure      400 {string} string "Invalid ID"
// @Router       /maintenance/{id} [delete]
// @Security BearerToken
func (c maintenanceRequestController) Delete(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Attampt to delete transaction using id
	err = c.service.Delete(idParameter)

	// If error detected
	if err != nil {
		http.Error(w, "Failed maintenance request deletion", http.StatusBadRequest)
		return
	}
	// Else write success
	w.Write([]byte("Deletion successful!"))
}

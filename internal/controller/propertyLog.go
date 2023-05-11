package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/service"
	"github.com/go-chi/chi"
)

type PropertyLogController interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	Find(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type propertyLogController struct {
	service service.PropertyLogService
}

func NewPropertyLogController(service service.PropertyLogService) PropertyLogController {
	return &propertyLogController{service}
}

// API/PROPERTY-LOGS
// Find a list of Property log messages
// @Summary      Find a list of property log messages
// @Description  Accepts limit, offset, and order params and returns list of log messages
// @Tags         Property Log
// @Accept       json
// @Produce      json
// @Param        limit   path      int  true  "limit"
// @Param        offset   path      int  true  "offset"
// @Param        order   path      int  true  "order by"
// @Success      200 {object} []db.PropertyLog
// @Failure      400 {string} string "Can't find property log messages"
// @Failure      400 {string} string "Must include limit parameter with a max value of 50"
// @Router       /property-logs [get]
// @Security BearerToken
func (c propertyLogController) FindAll(w http.ResponseWriter, r *http.Request) {
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

	// Query database for all log messages using query params
	found, err := c.service.FindAll(limit, offset, orderBy)
	if err != nil {
		http.Error(w, "Can't find property log messages", http.StatusBadRequest)
		return
	}
	// Write found log messages to response
	err = helpers.WriteAsJSON(w, found)
	if err != nil {
		http.Error(w, "Can't find property log messages", http.StatusBadRequest)
		fmt.Println("error writing property log messages to response: ", err)
		return
	}
}

// Find a created property log message by ID
// @Summary      Find property log message by ID
// @Description  Find a property log message by ID
// @Tags         Property Log
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Property Log ID"
// @Success      200 {object} db.PropertyLog
// @Failure      400 {string} string "Can't find property log message with ID:"
// @Failure      400 {string} string "Invalid ID"
// @Router       /property-logs/{id} [get]
// @Security BearerToken
func (c propertyLogController) Find(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	// Query database for property log message using ID
	found, err := c.service.FindById(idParameter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find property log message with ID: %v\n", idParameter), http.StatusBadRequest)
		return
	}
	// Write found property log message to response
	err = helpers.WriteAsJSON(w, found)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find property log message with ID: %v\n", idParameter), http.StatusBadRequest)
		return
	}
}

// Create a new property log message
// @Summary      Create property log message
// @Description  Creates a new property log message
// @Tags         Property Log
// @Accept       json
// @Produce      json
// @Param        feature body models.CreatePropertyLog true "New Property Log Json"
// @Success      201 {string} string "Property log message creation successful!"
// @Failure      400 {string} string "Property log message creation failed."
// @Router       /property-logs [post]
func (c propertyLogController) Create(w http.ResponseWriter, r *http.Request) {
	// Init
	var recvLog models.RecvPropertyLog
	// Decode request body as JSON and store
	err := json.NewDecoder(r.Body).Decode(&recvLog)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&recvLog)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		helpers.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	//
	// Validate the token
	tokenData, err := auth.ValidateAndParseToken(w, r)
	fmt.Println("tokendata received: ", tokenData)
	// If error detected
	if err != nil {
		http.Error(w, "Error parsing authentication token", http.StatusForbidden)
		return
	}
	// Convert user id from token to int and store
	userIdFromToken, err := strconv.Atoi(tokenData.UserID)
	if err != nil {
		http.Error(w, "Issue with user id from token", http.StatusBadRequest)
		return
	}

	// Convert DTO to service required input model
	var propLog = models.CreatePropertyLog{
		LogMessage: recvLog.LogMessage,
		User:       db.User{ID: uint(userIdFromToken)},
		// All access through this handler must automatically apply a field value for the property log type
		Type: "input",
		// Having an issue here with the property field not being recognized as a property model
		Property: recvLog.Property,
	}

	fmt.Printf("propLog: %+v\n", propLog)

	// Create property log message
	_, createErr := c.service.Create(&propLog)
	if createErr != nil {
		fmt.Printf("Issue with prop log message creation: %v\n", createErr)
		http.Error(w, "Property log message creation failed.", http.StatusBadRequest)
		return
	}

	// Set status to created
	w.WriteHeader(http.StatusCreated)
	// Send user success message in body
	w.Write([]byte("Property log message creation successful!"))
}

// Update a property log message (using URL parameter id)
// @Summary      Update property log message
// @Description  Updates an existing property log message
// @Tags         Property Log
// @Accept       json
// @Produce      json
// @Param        feature body models.UpdatePropertyLog true "Update Property Log Json"
// @Param        id   path      int  true  "Log Message ID"
// @Success      200 {object} db.PropertyLog
// @Failure      400 {string} string "Failed property log message update"
// @Failure      400 {string} string "Invalid ID"
// @Failure      403 {string} string "Authentication Token not detected"
// @Router       /property-logs/{id} [put]
// @Security BearerToken
func (c propertyLogController) Update(w http.ResponseWriter, r *http.Request) {
	// Grab ID URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Init
	var log models.UpdatePropertyLog
	// Decode request body as JSON and store
	err = json.NewDecoder(r.Body).Decode(&log)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&log)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		helpers.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	// Update property log message
	updatedLogMessage, createErr := c.service.Update(idParameter, &log)
	if createErr != nil {
		http.Error(w, fmt.Sprintf("Failed property feature update: %s", createErr), http.StatusBadRequest)
		return
	}
	// Write property log message to output
	err = helpers.WriteAsJSON(w, updatedLogMessage)
	if err != nil {
		fmt.Printf("Error encountered when writing to JSON. Err: %s", err)
	}
}

// Delete property log message (using URL parameter id)
// @Summary      Delete Property log message
// @Description  Deletes an existing property log message
// @Tags         Property log message
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Log message ID"
// @Success      200 {string} string "Deletion successful!"
// @Failure      400 {string} string "Failed property log message deletion"
// @Failure      400 {string} string "Invalid ID"

// @Router       /property-logs/{id} [delete]
// @Security BearerToken
func (c propertyLogController) Delete(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Attampt to delete property log message using id
	err = c.service.Delete(idParameter)

	// If error detected
	if err != nil {
		http.Error(w, "Failed property log message deletion", http.StatusBadRequest)
		return
	}
	// Else write success
	w.Write([]byte("Deletion successful!"))
}

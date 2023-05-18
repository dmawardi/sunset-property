package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/db"
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
	log     service.PropertyLogService
}

func NewPropertyController(service service.PropertyService, log service.PropertyLogService) PropertyController {
	return &propertyController{service, log}
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

	// Generate a property log message frop property update in preparation for successful update
	genPropLogMessage := buildPropLogUpdate(prop)
	// Grab user id from token
	userID, err := auth.GetUserIDFromToken(w, r)
	if err != nil {
		http.Error(w, "Authentication Token not detected", http.StatusForbidden)
		return
	}

	// Update property
	updatedProperty, createErr := c.service.Update(idParameter, &prop)
	if createErr != nil {
		http.Error(w, fmt.Sprintf("Failed property update: %s", createErr), http.StatusBadRequest)
		return
	}
	// Proceed to update the property log with the update
	c.log.Create(&models.CreatePropertyLog{
		// From URL parameter
		Property: db.Property{
			ID: uint(idParameter),
		},
		// From JWT token
		User: db.User{ID: uint(userID)},
		// Generated message
		LogMessage: genPropLogMessage,
		Type:       "gen",
	})

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

// Build a log string for property updates
func buildPropLogUpdate(updateStruct interface{}) string {
	// Log update
	var updateString string = "UPDATE: "
	// Iterate through key value pairs within the struct
	// Get the type of the struct
	t := reflect.TypeOf(updateStruct)

	// Iterate through the fields of the struct
	for i := 0; i < t.NumField(); i++ {
		// Get the field
		field := t.Field(i)
		// Get the value of the field
		value := reflect.ValueOf(updateStruct).Field(i).Interface()
		// Get the type of value
		valueType := reflect.TypeOf(value)
		// fmt.Printf("\n%s (%s): %v", field.Name, valueType, value)

		// If value type is string
		if valueType.String() == "string" {
			// and not empty
			if value != "" {
				fmt.Printf("\nString value found in Field name: %v", field.Name)
				updateString += fmt.Sprintf("%s (%v), ", field.Name, value.(string)[0:5]+"...")
			}
			// else if value type is numeric
		} else if valueType.String() == "int" || valueType.String() == "int64" || valueType.String() == "float64" || valueType.String() == "float32" {
			// and not empty
			if value != "0" && value != "0.0" {
				updateString += fmt.Sprintf("%s, ", field.Name)
			}
			// else if value type is struct
		} else if strings.Contains(valueType.String(), "[]") {
			updateString += fmt.Sprintf("[]%s, ", field.Name)
		}

	}
	// Get length of string
	stringLength := len(updateString)
	// Remove last two characters of string (comma and space)
	removeFromEnd := stringLength - 2
	croppedLogMessage := updateString[:removeFromEnd]

	return croppedLogMessage
}

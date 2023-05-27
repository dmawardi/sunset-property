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

type VendorController interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	Find(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type vendorController struct {
	service service.VendorService
}

func NewVendorController(service service.VendorService) VendorController {
	return &vendorController{service}
}

// API/VENDORS
// Find a list of vendors
// @Summary      Find a list of vendors
// @Description  Accepts limit, offset, and order params and returns list of vendors
// @Tags         Vendors
// @Accept       json
// @Produce      json
// @Param        limit   path      int  true  "limit"
// @Param        offset   path      int  true  "offset"
// @Param        order   path      int  true  "order by"
// @Success      200 {object} []db.Vendor
// @Failure      400 {string} string "Can't find vendors"
// @Failure      400 {string} string "Must include limit parameter with a max value of 50"
// @Router       /vendors [get]
// @Security BearerToken
func (c vendorController) FindAll(w http.ResponseWriter, r *http.Request) {
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

	// Query database for all vendors using query params
	found, err := c.service.FindAll(limit, offset, orderBy)
	if err != nil {
		http.Error(w, "Can't find vendors", http.StatusBadRequest)
		return
	}
	// Write found vendors to response
	err = helpers.WriteAsJSON(w, found)
	if err != nil {
		http.Error(w, "Can't find vendors", http.StatusBadRequest)
		fmt.Println("error writing vendors to response: ", err)
		return
	}
}

// Find a created vendor by ID
// @Summary      Find vendor by ID
// @Description  Find a vendor by ID
// @Tags         Vendors
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Vendor ID"
// @Success      200 {object} db.Vendor
// @Failure      400 {string} string "Can't find vendor with ID:"
// @Failure      400 {string} string "Invalid ID"
// @Router       /vendors/{id} [get]
// @Security BearerToken
func (c vendorController) Find(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	// Query database for vendor using ID
	found, err := c.service.FindById(idParameter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find vendor with ID: %v\n%v", idParameter, err), http.StatusBadRequest)
		return
	}
	// Write found item to response
	err = helpers.WriteAsJSON(w, found)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find vendor with ID: %v\n", idParameter), http.StatusBadRequest)
		return
	}
}

// Create a new vendor
// @Summary      Create vendor
// @Description  Creates a new vendor
// @Tags         Vendors
// @Accept       json
// @Produce      json
// @Param        vendor body models.CreateVendor true "New Vendor Json"
// @Success      201 {string} string "Vendor creation successful!"
// @Failure      400 {string} string "Vendor creation failed."
// @Router       /vendors [post]
func (c vendorController) Create(w http.ResponseWriter, r *http.Request) {
	// Init
	var vendor models.CreateVendor
	// Decode request body from JSON and store
	err := json.NewDecoder(r.Body).Decode(&vendor)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&vendor)
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
	_, createErr := c.service.Create(&vendor)
	if createErr != nil {
		http.Error(w, "Vendor creation failed:."+createErr.Error(), http.StatusBadRequest)
		return
	}

	// Set status to created
	w.WriteHeader(http.StatusCreated)
	// Send user success message in body
	w.Write([]byte("Vendor creation successful!"))
}

// Update a vendor (using URL parameter id)
// @Summary      Update vendor
// @Description  Updates an existing vendor
// @Tags         Vendors
// @Accept       json
// @Produce      json
// @Param        vendor body models.UpdateWorkType true "Update Work Type Json"
// @Param        id   path      int  true  "Vendor ID"
// @Success      200 {object} db.MaintenanceRequest
// @Failure      400 {string} string "Failed vendor update"
// @Failure      400 {string} string "Invalid ID"
// @Failure      403 {string} string "Authentication Token not detected"
// @Router       /vendorss/{id} [put]
// @Security BearerToken
func (c vendorController) Update(w http.ResponseWriter, r *http.Request) {
	// Grab ID URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Init
	var vendor models.UpdateVendor
	// Decode request body from JSON and store
	err = json.NewDecoder(r.Body).Decode(&vendor)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&vendor)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		helpers.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	// Update vendor in db
	updatedVendor, createErr := c.service.Update(idParameter, &vendor)
	if createErr != nil {
		http.Error(w, fmt.Sprintf("Failed vendor update: %s", createErr), http.StatusBadRequest)
		return
	}
	// Write work type to output
	err = helpers.WriteAsJSON(w, updatedVendor)
	if err != nil {
		fmt.Printf("Error encountered when writing to JSON. Err: %s", err)
	}
}

// Delete vendor (using URL parameter id)
// @Summary      Delete vendor
// @Description  Deletes an existing vendor
// @Tags         Vendors
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Vendor ID"
// @Success      200 {string} string "Deletion successful!"
// @Failure      400 {string} string "Failed vendor deletion"
// @Failure      400 {string} string "Invalid ID"
// @Router       /vendors/{id} [delete]
// @Security BearerToken
func (c vendorController) Delete(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Attampt to delete vendor using id
	err = c.service.Delete(idParameter)

	// If error detected
	if err != nil {
		http.Error(w, "Failed vendor deletion", http.StatusBadRequest)
		return
	}
	// Else write success
	w.Write([]byte("Deletion successful!"))
}

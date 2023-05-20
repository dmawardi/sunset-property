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

type TransactionController interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	Find(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type transactionController struct {
	service service.TransactionService
}

func NewTransactionController(service service.TransactionService) TransactionController {
	return &transactionController{service}
}

// API/TRANSACTIONS
// Find a list of transactions
// @Summary      Find a list of transactions
// @Description  Accepts limit, offset, and order params and returns list of transactions
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        limit   path      int  true  "limit"
// @Param        offset   path      int  true  "offset"
// @Param        order   path      int  true  "order by"
// @Success      200 {object} []db.Transaction
// @Failure      400 {string} string "Can't find transactions"
// @Failure      400 {string} string "Must include limit parameter with a max value of 50"
// @Router       /transactions [get]
// @Security BearerToken
func (c transactionController) FindAll(w http.ResponseWriter, r *http.Request) {
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

	// Query database for all transactions using query params
	found, err := c.service.FindAll(limit, offset, orderBy)
	if err != nil {
		http.Error(w, "Can't find transactions", http.StatusBadRequest)
		return
	}
	// Write found transactions to response
	err = helpers.WriteAsJSON(w, found)
	if err != nil {
		http.Error(w, "Can't find transactions", http.StatusBadRequest)
		fmt.Println("error writing transactions to response: ", err)
		return
	}
}

// Find a created transaction by ID
// @Summary      Find transaction by ID
// @Description  Find a transaction by ID
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Transaction ID"
// @Success      200 {object} db.Transaction
// @Failure      400 {string} string "Can't find transaction with ID:"
// @Failure      400 {string} string "Invalid ID"
// @Router       /transactions/{id} [get]
// @Security BearerToken
func (c transactionController) Find(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	// Query database for transaction using ID
	found, err := c.service.FindById(idParameter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find transaction with ID: %v\n%v", idParameter, err), http.StatusBadRequest)
		return
	}
	// Write found transaction to response
	err = helpers.WriteAsJSON(w, found)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find transaction with ID: %v\n", idParameter), http.StatusBadRequest)
		return
	}
}

// Create a new transaction
// @Summary      Create transaction
// @Description  Creates a new transaction
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        transaction body models.CreateTransaction true "New Transaction Json"
// @Success      201 {string} string "Transaction creation successful!"
// @Failure      400 {string} string "Transaction creation failed."
// @Router       /transactions [post]
func (c transactionController) Create(w http.ResponseWriter, r *http.Request) {
	// Init
	var transaction models.CreateTransaction
	// Decode request body from JSON and store
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&transaction)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		helpers.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	// Create transaction in db
	_, createErr := c.service.Create(&transaction)
	if createErr != nil {
		fmt.Printf("Issue with transaction creation: %v\n", createErr)
		http.Error(w, "Transaction creation failed.", http.StatusBadRequest)
		return
	}

	// Set status to created
	w.WriteHeader(http.StatusCreated)
	// Send user success message in body
	w.Write([]byte("Transaction creation successful!"))
}

// Update a transaction (using URL parameter id)
// @Summary      Update transaction
// @Description  Updates an existing transaction
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        transaction body models.UpdateTransaction true "Update Transaction Json"
// @Param        id   path      int  true  "Transaction ID"
// @Success      200 {object} db.Transaction
// @Failure      400 {string} string "Failed transaction update"
// @Failure      400 {string} string "Invalid ID"
// @Failure      403 {string} string "Authentication Token not detected"
// @Router       /transactions/{id} [put]
// @Security BearerToken
func (c transactionController) Update(w http.ResponseWriter, r *http.Request) {
	// Grab ID URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Init
	var transaction models.UpdateTransaction
	// Decode request body from JSON and store
	err = json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&transaction)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		helpers.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	// Update transaction
	updatedTransaction, createErr := c.service.Update(idParameter, &transaction)
	if createErr != nil {
		http.Error(w, fmt.Sprintf("Failed transaction update: %s", createErr), http.StatusBadRequest)
		return
	}
	// Write transaction to output
	err = helpers.WriteAsJSON(w, updatedTransaction)
	if err != nil {
		fmt.Printf("Error encountered when writing to JSON. Err: %s", err)
	}
}

// Delete transaction (using URL parameter id)
// @Summary      Delete transaction
// @Description  Deletes an existing transaction
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Transaction ID"
// @Success      200 {string} string "Deletion successful!"
// @Failure      400 {string} string "Failed transaction deletion"
// @Failure      400 {string} string "Invalid ID"
// @Router       /transactions/{id} [delete]
// @Security BearerToken
func (c transactionController) Delete(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Failed transaction deletion", http.StatusBadRequest)
		return
	}
	// Else write success
	w.Write([]byte("Deletion successful!"))
}

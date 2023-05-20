package repository

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	FindAll(int, int, string) (*[]db.Transaction, error)
	FindById(int) (*db.Transaction, error)
	Create(*db.Transaction) (*db.Transaction, error)
	Update(int, *db.Transaction) (*db.Transaction, error)
	Delete(int) error
}

type transactionRepository struct {
	DB *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

// Creates a transaction in the database
func (r *transactionRepository) Create(transaction *db.Transaction) (*db.Transaction, error) {
	// Create new log message in database
	result := r.DB.Create(&transaction)
	if result.Error != nil {
		return nil, fmt.Errorf("failed creating transaction: %w", result.Error)
	}

	return transaction, nil
}

// Find a list of transactions in the database
func (r *transactionRepository) FindAll(limit int, offset int, order string) (*[]db.Transaction, error) {
	// Query all transactions based on the received parameters
	transactions, err := QueryAllTransactionsBasedOnParams(limit, offset, order, r.DB)
	if err != nil {
		fmt.Printf("Error querying db for list of transactions: %s", err)
		return nil, err
	}

	return &transactions, nil
}

// Find a transaction in database by ID
func (r *transactionRepository) FindById(id int) (*db.Transaction, error) {
	// Create an empty ref object of type transaction
	transaction := db.Transaction{}
	// Grab transaction from db if exists
	result := r.DB.Preload("Property").Preload("Contacts").First(&transaction, id)

	// If error detected
	if result.Error != nil {
		return nil, result.Error
	}
	// else
	return &transaction, nil
}

// Delete transaction in database
func (r *transactionRepository) Delete(id int) error {
	// Create an empty ref object of type transaction
	transaction := db.Transaction{}
	// Delete transaction from db if exists
	result := r.DB.Delete(&transaction, id)

	// If error detected
	if result.Error != nil {
		fmt.Println("error in deleting transaction: ", result.Error)
		return result.Error
	}
	// else
	return nil
}

// Updates transaction in database
func (r *transactionRepository) Update(id int, transaction *db.Transaction) (*db.Transaction, error) {
	// Init
	var err error
	// Find transaction by id to ensure it exists
	foundTransaction, err := r.FindById(id)
	if err != nil {
		fmt.Println("Transaction to update not found: ", err)
		return nil, err
	}

	// Update found transaction with details from transaction
	updateResult := r.DB.Model(&foundTransaction).Updates(transaction)
	if updateResult.Error != nil {
		fmt.Println("Transaction update failed: ", updateResult.Error)
		return nil, updateResult.Error
	}

	// Build associations
	assResult := r.DB.Model(&foundTransaction).Association("Contacts").Append(transaction.Contacts)
	// Check if association update failed
	if assResult != nil {
		fmt.Println("Property association update failed: ", assResult)
		return nil, assResult
	}

	// Retrieve updated transaction by id
	updatedTransaction, err := r.FindById(id)
	if err != nil {
		fmt.Println("Updated Transaction not found: ", err)
		return nil, err
	}
	return updatedTransaction, nil
}

// Takes limit, offset, and order parameters, builds a query and executes returning a list of transactions
func QueryAllTransactionsBasedOnParams(limit int, offset int, order string, dbClient *gorm.DB) ([]db.Transaction, error) {
	// Build model to query database
	transaction := []db.Transaction{}
	// Build base query for property log messages table
	query := dbClient.Model(&transaction)

	// Add parameters into query as needed
	if limit != 0 {
		query.Limit(limit)
	}
	if offset != 0 {
		query.Offset(offset)
	}
	// order format should be "column_name ASC/DESC" eg. "created_at ASC"
	if order != "" {
		query.Order(order)
	} else {
		query.Order("created_at DESC")
	}
	// Query database
	result := query.Find(&transaction)
	if result.Error != nil {
		return nil, result.Error
	}
	// Return if no errors with result
	return transaction, nil
}

package service

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/repository"
)

type TransactionService interface {
	FindAll(int, int, string) (*[]db.Transaction, error)
	FindById(int) (*db.Transaction, error)
	Create(*models.CreateTransaction) (*db.Transaction, error)
	Update(int, *models.UpdateTransaction) (*db.Transaction, error)
	Delete(int) error
}

type transactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{repo}
}

// Creates a transaction
func (s *transactionService) Create(transaction *models.CreateTransaction) (*db.Transaction, error) {
	// Create a new transaction from DTO
	transToCreate := db.Transaction{
		Type:             transaction.Type,
		Agency:           transaction.Agency,
		AgencyName:       transaction.AgencyName,
		IsLease:          transaction.IsLease,
		TenancyType:      transaction.TenancyType,
		TransactionNotes: transaction.TransactionNotes,
		TransactionValue: transaction.TransactionValue,
		Fee:              transaction.Fee,
		Property:         transaction.Property,
		TaskID:           transaction.Task.ID,
	}

	// Create transaction in database
	createdTransaction, err := s.repo.Create(&transToCreate)
	if err != nil {
		return nil, fmt.Errorf("failed creating transaction: %w", err)
	}

	return createdTransaction, nil
}

// Find a list of transactions
func (s *transactionService) FindAll(limit int, offset int, order string) (*[]db.Transaction, error) {
	transactions, err := s.repo.FindAll(limit, offset, order)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// Find transaction in database by ID
func (s *transactionService) FindById(id int) (*db.Transaction, error) {
	// Find transaction by id
	transaction, err := s.repo.FindById(id)
	// If error detected
	if err != nil {
		fmt.Println("error in finding transaction: ", err)
		return nil, err
	}
	// else
	return transaction, nil
}

// Delete transaction in database
func (s *transactionService) Delete(id int) error {
	err := s.repo.Delete(id)
	// If error detected
	if err != nil {
		fmt.Println("error in deleting transaction: ", err)
		return err
	}
	// else
	return nil
}

// Updates transaction in database
func (s *transactionService) Update(id int, transaction *models.UpdateTransaction) (*db.Transaction, error) {
	// Create a new transaction from DTO
	transToUpdate := db.Transaction{
		Type:                  transaction.Type,
		Agency:                transaction.Agency,
		AgencyName:            transaction.AgencyName,
		IsLease:               transaction.IsLease,
		TenancyType:           transaction.TenancyType,
		TransactionNotes:      transaction.TransactionNotes,
		TransactionValue:      transaction.TransactionValue,
		Fee:                   transaction.Fee,
		TransactionCompletion: transaction.TransactionCompletion,
		Contacts:              transaction.Contacts,
	}

	// Update using repo
	updatedTransaction, err := s.repo.Update(id, &transToUpdate)
	if err != nil {
		return nil, err
	}

	return updatedTransaction, nil
}

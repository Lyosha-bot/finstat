package service

import (
	"finstat/internal/models"
	"finstat/internal/repository"
	"time"

	"github.com/shopspring/decimal"
)

type TransactionService struct {
	repo repository.Transaction
}

func NewTransactionService(repo repository.Transaction) *TransactionService {
	return &TransactionService{
		repo: repo,
	}
}

func (s *TransactionService) InsertTransaction(userID uint, value decimal.Decimal, categoryID uint, description string, date time.Time) (uint, error) {
	return s.repo.InsertTransaction(userID, value, categoryID, description, date)
}

func (s *TransactionService) UpdateTransaction(userID uint, transactionID uint, newValue decimal.Decimal, newCategoryID uint, newDescription string, newDate time.Time) (bool, error) {
	return s.repo.UpdateTransaction(userID, transactionID, newValue, newCategoryID, newDescription, newDate)
}

func (s *TransactionService) DeleteTransaction(userID uint, transactionID uint) (bool, error) {
	return s.repo.DeleteTransaction(userID, transactionID)
}

func (s *TransactionService) Transactions(userID, limit, page uint, from, to *time.Time, transactionType int, categories []uint) ([]models.Transaction, error) {
	return s.repo.Transactions(userID, limit, page, from, to, transactionType, categories)
}

package service

import (
	"finstat/internal/repository"
	"time"

	"github.com/shopspring/decimal"
)

type Transaction = repository.Transaction

type TransactionRepo interface {
	AddTransaction(userID uint, value decimal.Decimal, categoryID uint, description string, date time.Time) (uint, error)
	UpdateTransaction(userID uint, transactionID uint, newValue decimal.Decimal, newCategoryID uint, newDescription string, newDate time.Time) (bool, error)
	DeleteTransaction(userID uint, transactionID uint) (bool, error)
	Transactions(userID, limit, page uint) ([]Transaction, error)
	TransactionsInPeriod(userID uint, limit, page uint, from, to time.Time) ([]Transaction, error)
	TransactionsFromDate(userID uint, limit, page uint, date time.Time, order bool) ([]Transaction, error)
}

type TransactionService struct {
	repo TransactionRepo
}

func NewTransactionService(repo TransactionRepo) *TransactionService {
	return &TransactionService{
		repo: repo,
	}
}

func (s *TransactionService) AddTransaction(userID uint, value decimal.Decimal, categoryID uint, description string, date time.Time) (uint, error) {
	return s.repo.AddTransaction(userID, value, categoryID, description, date)
}

func (s *TransactionService) UpdateTransaction(userID uint, transactionID uint, newValue decimal.Decimal, newCategoryID uint, newDescription string, newDate time.Time) (bool, error) {
	return s.repo.UpdateTransaction(userID, transactionID, newValue, newCategoryID, newDescription, newDate)
}

func (s *TransactionService) DeleteTransaction(userID uint, transactionID uint) (bool, error) {
	return s.repo.DeleteTransaction(userID, transactionID)
}

func (s *TransactionService) Transactions(userID, limit, page uint) ([]repository.Transaction, error) {
	return s.repo.Transactions(userID, limit, page)
}

func (s *TransactionService) TransactionsInPeriod(userID, limit, page uint, from, to time.Time) ([]repository.Transaction, error) {
	return s.repo.TransactionsInPeriod(userID, limit, page, from, to)
}

func (s *TransactionService) TransactionsFromDate(userID uint, limit, page uint, date time.Time, order bool) ([]repository.Transaction, error) {
	return s.repo.TransactionsFromDate(userID, limit, page, date, order)
}

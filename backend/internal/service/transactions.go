package service

import (
	ewrap "finstat/internal/lib"
	"finstat/internal/repository"
	"time"

	"github.com/shopspring/decimal"
)

type Transaction = repository.Transaction

type TransactionRepo interface {
	AddTransaction(userID uint, amount decimal.Decimal, categoryID uint, description string, date time.Time) (uint, error)
	UpdateTransaction(userID uint, transactionID uint, newAmount decimal.Decimal, newCategoryID uint, newDescription string, newDate time.Time) error
	DeleteTransaction(userID uint, transactionID uint) error
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

func (s *TransactionService) AddTransaction(userID uint, amount decimal.Decimal, categoryID uint, description string, date time.Time) (uint, error) {
	id, err := s.repo.AddTransaction(userID, amount, categoryID, description, date)
	if err != nil {
		return 0, ewrap.Wrap("Couldn't insert transaction", err)
	}

	return id, err
}

func (s *TransactionService) UpdateTransaction(userID uint, transactionID uint, newAmount decimal.Decimal, newCategoryID uint, newDescription string, newDate time.Time) error {
	if err := s.repo.UpdateTransaction(userID, transactionID, newAmount, newCategoryID, newDescription, newDate); err != nil {
		return ewrap.Wrap("Couldn't update transaction", err)
	}

	return nil
}

func (s *TransactionService) DeleteTransaction(userID uint, transactionID uint) error {
	if err := s.repo.DeleteTransaction(userID, transactionID); err != nil {
		return ewrap.Wrap("Couldn't update transaction", err)
	}

	return nil
}

func (s *TransactionService) Transactions(userID, limit, page uint) ([]repository.Transaction, error) {
	transactions, err := s.repo.Transactions(userID, limit, page)
	if err != nil {
		return nil, ewrap.Wrap("Couldn't get latest transactions", err)
	}

	return transactions, err
}

func (s *TransactionService) TransactionsInPeriod(userID, limit, page uint, from, to time.Time) ([]repository.Transaction, error) {
	transactions, err := s.repo.TransactionsInPeriod(userID, limit, page, from, to)
	if err != nil {
		return nil, ewrap.Wrap("Couldn't get transactions in period", err)
	}

	return transactions, err
}

func (s *TransactionService) TransactionsFromDate(userID uint, limit, page uint, date time.Time, order bool) ([]repository.Transaction, error) {
	transactions, err := s.repo.TransactionsFromDate(userID, limit, page, date, order)
	if err != nil {
		return nil, ewrap.Wrap("Couldn't get transactions in period", err)
	}

	return transactions, err
}

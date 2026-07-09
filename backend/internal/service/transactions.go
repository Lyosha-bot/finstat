package service

import (
	ewrap "finstat/internal/lib"
	"finstat/internal/repository"
	"time"

	"github.com/shopspring/decimal"
)

type Transaction = repository.Transaction

type TransactionRepo interface {
	AddTransaction(userID uint, amount decimal.Decimal, description string, date time.Time) (uint, error)
	Transactions(userID, limit, page uint) ([]Transaction, error)
	TransactionsInPeriod(userID uint, limit, page uint, from, to time.Time) ([]Transaction, error)
	TransactionsFromDate(userID uint, limit, page uint, date time.Time, order bool) ([]Transaction, error)
}

type TransactionsService struct {
	repo TransactionRepo
}

func NewTransactionService(repo TransactionRepo) *TransactionsService {
	return &TransactionsService{
		repo: repo,
	}
}

func (s *TransactionsService) AddTransaction(userID uint, amount decimal.Decimal, description string, date time.Time) (uint, error) {
	id, err := s.repo.AddTransaction(userID, amount, description, date)
	if err != nil {
		return 0, ewrap.Wrap("Couldn't insert transaction", err)
	}

	return id, err
}

func (s *TransactionsService) Transactions(userID, limit, page uint) ([]repository.Transaction, error) {
	transactions, err := s.repo.Transactions(userID, limit, page)
	if err != nil {
		return nil, ewrap.Wrap("Couldn't get latest transactions", err)
	}

	return transactions, err
}

func (s *TransactionsService) TransactionsInPeriod(userID, limit, page uint, from, to time.Time) ([]repository.Transaction, error) {
	transactions, err := s.repo.TransactionsInPeriod(userID, limit, page, from, to)
	if err != nil {
		return nil, ewrap.Wrap("Couldn't get transactions in period", err)
	}

	return transactions, err
}

func (s *TransactionsService) TransactionsFromDate(userID uint, limit, page uint, date time.Time, order bool) ([]repository.Transaction, error) {
	transactions, err := s.repo.TransactionsFromDate(userID, limit, page, date, order)
	if err != nil {
		return nil, ewrap.Wrap("Couldn't get transactions in period", err)
	}

	return transactions, err
}
